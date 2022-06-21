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

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	kcorev1 "k8s.io/api/core/v1"
	kmetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ktypes "k8s.io/apimachinery/pkg/types"
	kclient "sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/sambatv/aws-auth-operator/apis/v1beta1"
)

const (
	APIVersion = "aws-auth.samba.tv/v1beta1"

	MapUserNamespace   = "default"
	MapUserKind        = "MapUser"
	MapUserARN         = "arn:aws:iam::123456789012:user/test-user"
	MapUserUsername    = "test-user"
	MapUserDescription = "A test user"
	MapUserEmail       = "test@samba.tv"
)

var _ = Describe("MapUser controller", func() {
	awsAuthObjectKey := kclient.ObjectKey{Name: "aws-auth", Namespace: "kube-system"}

	Context("When managing MapUser objects", func() {
		It("Should ensure we start with an empty kube-system/aws-auth configmap", func() {
			By("Getting the kube-system/aws-auth and looking at its data.")
			ctx := context.Background()
			var awsAuth kcorev1.ConfigMap
			err := k8sClient.Get(ctx, awsAuthObjectKey, &awsAuth)
			Expect(err).ToNot(HaveOccurred())
			Expect(len(awsAuth.Data)).Should(Equal(0))
		})

		It("Should add the MapUser", func() {
			By("Creating a new MapUser with a User ARN")
			ctx := context.Background()
			user := &v1beta1.MapUser{
				TypeMeta: kmetav1.TypeMeta{
					APIVersion: APIVersion,
					Kind:       MapUserKind,
				},
				ObjectMeta: kmetav1.ObjectMeta{
					Name:      MapUserUsername,
					Namespace: MapUserNamespace,
				},
				Spec: v1beta1.MapUserSpec{
					Description: MapUserDescription,
					Email:       MapUserEmail,
					Groups:      []string{},
					UserARN:     MapUserARN,
				},
			}
			Expect(k8sClient.Create(ctx, user)).Should(Succeed())

			userLookupKey := ktypes.NamespacedName{Name: MapUserUsername, Namespace: MapUserNamespace}
			createdUser := &v1beta1.MapUser{}

			// We'll need to retry getting this newly created MapUser, given that creation may not immediately happen.
			Eventually(func() bool {
				err := k8sClient.Get(ctx, userLookupKey, createdUser)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())
			// Ensure MapUser data is correct.
			Expect(createdUser.Spec.UserARN).Should(Equal(MapUserARN))
			Expect(createdUser.Spec.Description).Should(Equal(MapUserDescription))
			Expect(createdUser.Spec.Email).Should(Equal(MapUserEmail))
		})

		//It("Should ensure the MapUser data has been added to kube-system/aws-auth configmap data.mapUsers", func() {
		//	By("Reloading it and looking at data.mapUsers")
		//	ctx := context.Background()
		//	var awsAuth kcorev1.ConfigMap
		//	err := k8sClient.Get(ctx, awsAuthObjectKey, &awsAuth)
		//	Expect(err).ToNot(HaveOccurred())
		//	fmt.Fprintf(GinkgoWriter, "data=%v", awsAuth.Data)
		//	mapUsers, ok := awsAuth.Data["mapUsers"]
		//	Expect(ok).Should(Equal(true))
		//
		//	text, err := json.Marshal(awsAuth)
		//	Expect(err).ToNot(HaveOccurred())
		//	fmt.Fprintf(GinkgoWriter, "\nawsAuth=%s", text)
		//	Expect(len(mapUsers)).Should(Equal(1))
		//})
	})
})
