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
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"reflect"
	"time"

	"k8s.io/client-go/kubernetes"
)

// NewMapper returns a new Mapper object.
func NewMapper(client kubernetes.Interface, discardLogOutput bool) *Mapper {
	var mapper = &Mapper{}
	mapper.KubernetesClient = client

	if !discardLogOutput {
		log.SetOutput(ioutil.Discard)
	}
	return mapper
}

// Mapper is responsible for managing the auth map.
type Mapper struct {
	KubernetesClient kubernetes.Interface
}

// Remove removes a mapRole or mapUser from the auth map.
func (m *Mapper) Remove(args *Arguments) error {
	args.Validate()
	if args.WithRetries {
		return WithRetry(m.removeAuth, args)
	}
	return m.removeAuth(args)
}

func (m *Mapper) removeAuth(args *Arguments) error {
	authData, configMap, err := ReadAuthMap(m.KubernetesClient)
	if err != nil {
		return err
	}

	var removed bool

	if args.DataType == MapRoleData {
		var newRolesAuthMap []*MapRole
		for _, mapRole := range authData.MapRoles {
			if args.Username != mapRole.Username {
				newRolesAuthMap = append(newRolesAuthMap, mapRole)
			} else {
				removed = true
			}
		}
		authData.SetMapRoles(newRolesAuthMap)
	}

	if args.DataType == MapUserData {
		var newUsersAuthMap []*MapUser
		for _, mapUser := range authData.MapUsers {
			if args.Username != mapUser.Username {
				newUsersAuthMap = append(newUsersAuthMap, mapUser)
			} else {
				removed = true
			}
		}
		authData.SetMapUsers(newUsersAuthMap)
	}

	if !removed {
		return errors.New(fmt.Sprintf("%s with username '%s' not found in auth map", args.DataType, args.Username))
	}
	return UpdateAuthMap(m.KubernetesClient, authData, configMap)
}

// Upsert updates or inserts a mapRole or mapUser item into the auth map.
func (m *Mapper) Upsert(args *Arguments) error {
	args.Validate()
	if args.WithRetries {
		return WithRetry(m.upsertAuth, args)
	}
	return m.upsertAuth(args)
}

func (m *Mapper) upsertAuth(args *Arguments) error {
	authData, configMap, err := ReadAuthMap(m.KubernetesClient)
	if err != nil {
		return err
	}

	if args.DataType == MapRoleData {
		mapRole := NewMapRole(args.RoleARN, args.Username, args.Groups)
		newMap, ok := upsertRole(authData.MapRoles, mapRole)
		if ok {
			log.Printf("%s with username '%s' key has been updated\n", args.DataType, args.Username)
		} else {
			log.Printf("no updates needed to %s with username '%s'\n", args.DataType, args.Username)
		}
		authData.SetMapRoles(newMap)
	}

	if args.DataType == MapUserData {
		mapUser := NewMapUser(args.UserARN, args.Username, args.Groups)
		newMap, ok := upsertUser(authData.MapUsers, mapUser)
		if ok {
			log.Printf("%s with username '%s' key has been updated\n", args.DataType, args.Username)
		} else {
			log.Printf("%s with username '%s' key has been updated\n", args.DataType, args.Username)
		}
		authData.SetMapUsers(newMap)
	}

	return UpdateAuthMap(m.KubernetesClient, authData, configMap)
}

func upsertRole(authMaps []*MapRole, resource *MapRole) ([]*MapRole, bool) {
	var found, updated bool
	for _, existing := range authMaps {
		// Update existing role in auth map.
		if existing.Username == resource.Username {
			found = true
			if !reflect.DeepEqual(existing.Groups, resource.Groups) {
				existing.SetGroups(resource.Groups)
				updated = true
			}
			if existing.RoleARN != resource.RoleARN {
				existing.SetRoleARN(resource.RoleARN)
				updated = true
			}
		}
	}

	// Insert new role in auth map.
	if !found {
		updated = true
		authMaps = append(authMaps, resource)
	}
	return authMaps, updated
}

func upsertUser(authMaps []*MapUser, resource *MapUser) ([]*MapUser, bool) {
	var found, updated bool
	for _, existing := range authMaps {
		// Update existing user in auth map.
		if existing.UserARN == resource.UserARN {
			found = true
			if !reflect.DeepEqual(existing.Groups, resource.Groups) {
				existing.SetGroups(resource.Groups)
				updated = true
			}
			if existing.UserARN != resource.UserARN {
				existing.SetUserARN(resource.UserARN)
				updated = true
			}
		}
	}

	// Insert new user in auth map.
	if !found {
		updated = true
		authMaps = append(authMaps, resource)
	}
	return authMaps, updated
}

// Arguments are the arguments for management of the auth map.
type Arguments struct {
	OperationType OperationType
	DataType      DataType
	RoleARN       string
	UserARN       string
	Username      string
	Groups        []string
	WithRetries   bool
	MinRetryTime  time.Duration
	MaxRetryTime  time.Duration
	MaxRetryCount int
}

// Validate validates if all Arguments fields are valid.
func (args *Arguments) Validate() {
	if args.WithRetries && args.MaxRetryCount < 1 {
		log.Println("error: retry max count is invalid, must be greater than zero")
	}
	if args.Username == "" {
		log.Println("error: username not provided")
	}
	if args.OperationType == "" {
		log.Println("error: operation type not provided")
	}
	if args.OperationType != UpsertOperation && args.OperationType != RemoveOperation {
		log.Printf("error: operation type '%s' not valid\n", args.OperationType)
	}
	if args.DataType == "" {
		log.Println("error: data type not provided")
	}
	if args.DataType != MapRoleData && args.DataType != MapUserData {
		log.Printf("error: data type '%s' not valid\n", args.DataType)
	}
	if args.OperationType == UpsertOperation && args.DataType == MapRoleData && args.RoleARN == "" {
		log.Println("error: role arn not provided")
	}
	if args.OperationType == UpsertOperation && args.DataType == MapUserData && args.UserARN == "" {
		log.Println("error: user arn not provided")
	}
}

// OperationType indicates the auth map management operation.
type OperationType string

const (
	UpsertOperation OperationType = "upsert"
	RemoveOperation OperationType = "remove"
)

// DataType indicates the auth map management scope.
type DataType string

const (
	MapRoleData DataType = "mapRole"
	MapUserData DataType = "mapUser"
)
