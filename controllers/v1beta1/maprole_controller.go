/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1beta1

import (
	"context"

	"github.com/go-logr/logr"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	pkgruntime "k8s.io/apimachinery/pkg/runtime"
	ctrlruntime "sigs.k8s.io/controller-runtime"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/sambatv/aws-auth-operator/apis/v1beta1"
	"github.com/sambatv/aws-auth-operator/awsauth"
	"github.com/sambatv/aws-auth-operator/kube"
)

// MapRoleReconciler reconciles a MapRole object
type MapRoleReconciler struct {
	ctrlclient.Client
	Log    logr.Logger
	Scheme *pkgruntime.Scheme
}

//+kubebuilder:rbac:groups=core,resources=configmaps,verbs=get;list;watch;update;patch
//+kubebuilder:rbac:groups=core,resources=namespaces,verbs=get;list;watch
//+kubebuilder:rbac:groups=aws-auth.samba.tv,resources=maproles,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=aws-auth.samba.tv,resources=maproles/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=aws-auth.samba.tv,resources=maproles/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.8.3/pkg/reconcile
func (r *MapRoleReconciler) Reconcile(ctx context.Context, req ctrlruntime.Request) (ctrlruntime.Result, error) {
	mapRoleName := req.NamespacedName.Name

	log := r.Log.WithValues("MapRole", mapRoleName)
	log.Info("reconciling MapRole...")

	kubeClient, err := kube.GetClient()
	if err != nil {
		log.Error(err, "failure getting kube client")
		return ctrlruntime.Result{}, err
	}

	// Get a new aws auth service object.
	awsauthSvc, err := awsauth.NewService(&awsauth.ServiceConfig{
		KubeClient: kubeClient,
		Log:        r.Log,
	})
	if err != nil {
		log.Error(err, "failure creating new aws auth service")
		return ctrlruntime.Result{}, err
	}

	// Load the MapRole object by name (its AWS IAM role ARN).
	var mapRole v1beta1.MapRole
	if err := r.Get(ctx, req.NamespacedName, &mapRole); err != nil {
		// If any error other than a "NotFound" API error, it's a problem.
		statusErr, ok := err.(*apierrors.StatusError)
		if !ok || (ok && statusErr.ErrStatus.Reason != "NotFound") {
			log.Error(err, "failure getting MapRole")
			return ctrlruntime.Result{}, err
		}

		if err := awsauthSvc.RemoveMapRole(mapRoleName); err != nil {
			log.Error(err, "error removing mapRole data in aws-auth configmap")
			return ctrlruntime.Result{}, nil
		}
		log.Info("removed mapRole data in aws-auth configmap")
		return ctrlruntime.Result{}, nil
	}

	// Ensure that any changes are synced to the kube-system:aws-auth ConfigMap.
	if err := awsauthSvc.UpsertMapRole(mapRole.Name, awsauth.MapRole{
		RoleARN: mapRole.Spec.RoleARN,
		Groups:  mapRole.Spec.Groups,
	}); err != nil {
		log.Error(err, "error upserting MapRole in aws-auth")
		return ctrlruntime.Result{}, err
	}
	log.Info("upserted MapRole")
	return ctrlruntime.Result{}, nil
}

// SetupWithManager sets up the controller with the Mapper.
func (r *MapRoleReconciler) SetupWithManager(mgr ctrlruntime.Manager) error {
	return ctrlruntime.NewControllerManagedBy(mgr).
		For(&v1beta1.MapRole{}).
		Complete(r)
}
