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
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	kcorev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/sambatv/aws-auth-operator/apis/v1beta1"
)

const (
	MapRoleNamespace   = "default"
	MapRoleKind        = "MapRole"
	MapRoleARN         = "arn:aws:iam::123456789012:role/test-role"
	MapRoleUsername    = "test-role"
	MapRoleDescription = "A test role"
	MapRoleEmail       = "test@samba.tv"

	timeout  = time.Second * 10
	interval = time.Millisecond * 250
)

var _ = Describe("MapRole controller", func() {
	awsAuthObjectKey := client.ObjectKey{Name: "aws-auth", Namespace: "kube-system"}

	Context("When managing MapRole objects", func() {
		It("Should ensure we start with an empty kube-system/aws-auth configmap", func() {
			By("Getting the kube-system/aws-auth and looking at its data.")
			ctx := context.Background()
			var awsAuth kcorev1.ConfigMap
			err := k8sClient.Get(ctx, awsAuthObjectKey, &awsAuth)
			Expect(err).ToNot(HaveOccurred())
			Expect(len(awsAuth.Data)).Should(Equal(0))
		})

		It("Should add the MapRole", func() {
			By("Creating a new MapRole with a MapRoleARN")
			ctx := context.Background()
			mapRole := &v1beta1.MapRole{
				TypeMeta: metav1.TypeMeta{
					APIVersion: APIVersion,
					Kind:       MapRoleKind,
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      MapRoleUsername,
					Namespace: MapRoleNamespace,
				},
				Spec: v1beta1.MapRoleSpec{
					Description: MapRoleDescription,
					Email:       MapRoleEmail,
					Groups:      []string{},
					RoleARN:     MapRoleARN,
				},
			}
			Expect(k8sClient.Create(ctx, mapRole)).Should(Succeed())

			roleLookupKey := types.NamespacedName{Name: MapRoleUsername, Namespace: MapRoleNamespace}
			createdRole := &v1beta1.MapRole{}

			// We'll need to retry getting this newly created MapRole, given that creation may not immediately happen.
			Eventually(func() bool {
				err := k8sClient.Get(ctx, roleLookupKey, createdRole)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())
			// Ensure Role data is correct.
			Expect(createdRole.Spec.Description).Should(Equal(MapRoleDescription))
			Expect(createdRole.Spec.Email).Should(Equal(MapRoleEmail))
			Expect(createdRole.Spec.RoleARN).Should(Equal(MapRoleARN))
		})

		//It("Should ensure the Role data has been added to kube-system/aws-auth configmap data.mapRoles", func() {
		//	By("Reloading it and looking at data.mapRoles")
		//	ctx := context.Background()
		//	var awsAuth kcorev1.ConfigMap
		//	err := k8sClient.Get(ctx, awsAuthObjectKey, &awsAuth)
		//	Expect(err).ToNot(HaveOccurred())
		//	fmt.Fprintf(GinkgoWriter, "data=%v", awsAuth.Data)
		//	mapRoles, ok := awsAuth.Data["mapRoles"]
		//	Expect(ok).Should(Equal(true))
		//
		//	text, err := json.Marshal(awsAuth)
		//	Expect(err).ToNot(HaveOccurred())
		//	fmt.Fprintf(GinkgoWriter, "\nawsAuth=%s", text)
		//	Expect(len(mapRoles)).Should(Equal(1))
		//})
	})
})
