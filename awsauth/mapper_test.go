/*

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

package awsauth

import (
	"testing"
	"time"

	"github.com/onsi/gomega"
	"k8s.io/client-go/kubernetes/fake"
)

var testARNs = map[string]string{
	"node-1": "arn:aws:iam::00000000000:role/node-1",
	"node-2": "arn:aws:iam::00000000000:role/node-2",
	"user-1": "arn:aws:iam::00000000000:user/user-1",
	"user-2": "arn:aws:iam::00000000000:user/user-2",
}

func TestMapper_Remove(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterTestingT(t)
	client := fake.NewSimpleClientset()
	mapper := NewMapper(client, true)
	createMockConfigMap(client)

	err := mapper.Remove(&Arguments{
		OperationType: RemoveOperation,
		DataType:      MapRoleData,
		Username:      "system:node:{{EC2PrivateDNSName}}",
	})
	g.Expect(err).NotTo(gomega.HaveOccurred())

	err = mapper.Remove(&Arguments{
		OperationType: RemoveOperation,
		DataType:      MapUserData,
		Username:      "admin",
		UserARN:       "doesn't matter",
	})
	g.Expect(err).NotTo(gomega.HaveOccurred())

	auth, _, err := ReadAuthMap(client)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(len(auth.MapRoles)).To(gomega.Equal(0))
	g.Expect(len(auth.MapUsers)).To(gomega.Equal(0))
}

func TestMapper_RemoveNotFound(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterTestingT(t)
	client := fake.NewSimpleClientset()
	mapper := NewMapper(client, true)
	createMockConfigMap(client)

	err := mapper.Remove(&Arguments{
		OperationType: RemoveOperation,
		DataType:      MapRoleData,
		Username:      "system:node:{{EC2PrivateDNSName}}-na",
	})
	g.Expect(err).To(gomega.HaveOccurred())

	err = mapper.Remove(&Arguments{
		OperationType: RemoveOperation,
		DataType:      MapUserData,
		Username:      "admin-na",
		UserARN:       "doesn't matter",
	})
	g.Expect(err).To(gomega.HaveOccurred())

	auth, _, err := ReadAuthMap(client)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(len(auth.MapRoles)).To(gomega.Equal(1))
	g.Expect(len(auth.MapUsers)).To(gomega.Equal(1))
}

func TestMapper_RemoveWithRetries(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterTestingT(t)
	client := fake.NewSimpleClientset()
	mapper := NewMapper(client, true)
	createMockConfigMap(client)

	err := mapper.Remove(&Arguments{
		OperationType: RemoveOperation,
		DataType:      MapRoleData,
		Username:      "system:node:{{EC2PrivateDNSName}}",
		WithRetries:   true,
		MinRetryTime:  time.Millisecond * 1,
		MaxRetryTime:  time.Millisecond * 2,
		MaxRetryCount: 3,
	})
	g.Expect(err).NotTo(gomega.HaveOccurred())

	err = mapper.Remove(&Arguments{
		OperationType: RemoveOperation,
		DataType:      MapUserData,
		Username:      "admin",
		WithRetries:   true,
		MinRetryTime:  time.Millisecond * 1,
		MaxRetryTime:  time.Millisecond * 2,
		MaxRetryCount: 3,
	})
	g.Expect(err).NotTo(gomega.HaveOccurred())

	auth, _, err := ReadAuthMap(client)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(len(auth.MapRoles)).To(gomega.Equal(0))
	g.Expect(len(auth.MapUsers)).To(gomega.Equal(0))
}

func TestMapper_UpsertInsert(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterTestingT(t)
	client := fake.NewSimpleClientset()
	mapper := NewMapper(client, true)
	createMockConfigMap(client)

	err := mapper.Upsert(&Arguments{
		OperationType: UpsertOperation,
		DataType:      MapRoleData,
		RoleARN:       testARNs["node-2"],
		Username:      "system:node:{{EC2PrivateDNSName}}",
		Groups:        []string{"system:bootstrappers", "system:nodes"},
	})
	g.Expect(err).NotTo(gomega.HaveOccurred())

	err = mapper.Upsert(&Arguments{
		OperationType: UpsertOperation,
		DataType:      MapUserData,
		UserARN:       testARNs["user-2"],
		Username:      "admin",
		Groups:        []string{"system:masters"},
	})
	g.Expect(err).NotTo(gomega.HaveOccurred())

	auth, _, err := ReadAuthMap(client)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(len(auth.MapRoles)).To(gomega.Equal(1))
	g.Expect(len(auth.MapUsers)).To(gomega.Equal(2))
}

func TestMapper_UpsertUpdate(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterTestingT(t)
	client := fake.NewSimpleClientset()
	mapper := NewMapper(client, true)
	createMockConfigMap(client)

	err := mapper.Upsert(&Arguments{
		OperationType: UpsertOperation,
		DataType:      MapRoleData,
		RoleARN:       testARNs["node-1"],
		Username:      "this:is:a:test",
		Groups:        []string{"system:some-role"},
	})
	g.Expect(err).NotTo(gomega.HaveOccurred())

	err = mapper.Upsert(&Arguments{
		OperationType: UpsertOperation,
		DataType:      MapUserData,
		UserARN:       testARNs["user-1"],
		Username:      "admin",
		Groups:        []string{"system:some-role"},
	})
	g.Expect(err).NotTo(gomega.HaveOccurred())

	auth, _, err := ReadAuthMap(client)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(len(auth.MapRoles)).To(gomega.Equal(2))
	g.Expect(len(auth.MapUsers)).To(gomega.Equal(1))
	g.Expect(auth.MapRoles[0].RoleARN).To(gomega.Equal(testARNs["node-1"]))
	g.Expect(auth.MapRoles[0].Username).To(gomega.Equal("system:node:{{EC2PrivateDNSName}}"))
	g.Expect(auth.MapRoles[0].Groups).To(gomega.Equal([]string{"system:bootstrappers", "system:nodes"}))
	g.Expect(auth.MapUsers[0].UserARN).To(gomega.Equal(testARNs["user-1"]))
	g.Expect(auth.MapUsers[0].Username).To(gomega.Equal("admin"))
	g.Expect(auth.MapUsers[0].Groups).To(gomega.Equal([]string{"system:some-role"}))
}

func TestMapper_UpsertNotNeeded(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterTestingT(t)
	client := fake.NewSimpleClientset()
	mapper := NewMapper(client, true)
	createMockConfigMap(client)

	err := mapper.Upsert(&Arguments{
		OperationType: UpsertOperation,
		DataType:      MapRoleData,
		RoleARN:       testARNs["node-1"],
		Username:      "system:node:{{EC2PrivateDNSName}}",
		Groups:        []string{"system:bootstrappers", "system:nodes"},
	})
	g.Expect(err).NotTo(gomega.HaveOccurred())

	err = mapper.Upsert(&Arguments{
		OperationType: UpsertOperation,
		DataType:      MapUserData,
		UserARN:       testARNs["user-1"],
		Username:      "admin",
		Groups:        []string{"system:masters"},
	})
	g.Expect(err).NotTo(gomega.HaveOccurred())

	auth, _, err := ReadAuthMap(client)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(len(auth.MapRoles)).To(gomega.Equal(1))
	g.Expect(len(auth.MapUsers)).To(gomega.Equal(1))
}

func TestMapper_UpsertWithCreate(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterTestingT(t)
	client := fake.NewSimpleClientset()
	mapper := NewMapper(client, true)

	err := mapper.Upsert(&Arguments{
		OperationType: UpsertOperation,
		DataType:      MapRoleData,
		RoleARN:       testARNs["node-1"],
		Username:      "this:is:a:test",
		Groups:        []string{"system:some-role"},
	})
	g.Expect(err).NotTo(gomega.HaveOccurred())

	err = mapper.Upsert(&Arguments{
		OperationType: UpsertOperation,
		DataType:      MapUserData,
		UserARN:       testARNs["user-1"],
		Username:      "this:is:a:test",
		Groups:        []string{"system:some-role"},
	})
	g.Expect(err).NotTo(gomega.HaveOccurred())

	auth, _, err := ReadAuthMap(client)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(len(auth.MapRoles)).To(gomega.Equal(1))
	g.Expect(len(auth.MapUsers)).To(gomega.Equal(1))
	g.Expect(auth.MapRoles[0].RoleARN).To(gomega.Equal(testARNs["node-1"]))
	g.Expect(auth.MapRoles[0].Username).To(gomega.Equal("this:is:a:test"))
	g.Expect(auth.MapRoles[0].Groups).To(gomega.Equal([]string{"system:some-role"}))
	g.Expect(auth.MapUsers[0].UserARN).To(gomega.Equal(testARNs["user-1"]))
	g.Expect(auth.MapUsers[0].Username).To(gomega.Equal("this:is:a:test"))
	g.Expect(auth.MapUsers[0].Groups).To(gomega.Equal([]string{"system:some-role"}))
}

func TestMapper_UpsertWithRetries(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterTestingT(t)
	client := fake.NewSimpleClientset()
	mapper := NewMapper(client, true)
	createMockConfigMap(client)

	err := mapper.Upsert(&Arguments{
		OperationType: UpsertOperation,
		DataType:      MapRoleData,
		RoleARN:       testARNs["node-1"],
		Username:      "system:node:{{EC2PrivateDNSName}}",
		Groups:        []string{"system:some-role"},
		WithRetries:   true,
		MaxRetryCount: 12,
		MaxRetryTime:  1 * time.Millisecond,
		MinRetryTime:  1 * time.Millisecond,
	})
	g.Expect(err).NotTo(gomega.HaveOccurred())

	err = mapper.Upsert(&Arguments{
		OperationType: UpsertOperation,
		DataType:      MapUserData,
		UserARN:       testARNs["user-1"],
		Username:      "this:is:a:test",
		Groups:        []string{"system:some-role"},
	})
	g.Expect(err).NotTo(gomega.HaveOccurred())

	auth, _, err := ReadAuthMap(client)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(len(auth.MapRoles)).To(gomega.Equal(1))
	g.Expect(len(auth.MapUsers)).To(gomega.Equal(1))
	g.Expect(auth.MapRoles[0].RoleARN).To(gomega.Equal(testARNs["node-1"]))
	g.Expect(auth.MapRoles[0].Username).To(gomega.Equal("system:node:{{EC2PrivateDNSName}}"))
	g.Expect(auth.MapRoles[0].Groups).To(gomega.Equal([]string{"system:some-role"}))
	g.Expect(auth.MapUsers[0].UserARN).To(gomega.Equal(testARNs["user-1"]))
	g.Expect(auth.MapUsers[0].Username).To(gomega.Equal("admin"))
	g.Expect(auth.MapUsers[0].Groups).To(gomega.Equal([]string{"system:some-role"}))
}
