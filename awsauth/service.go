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

package awsauth

import (
	"errors"
	"time"

	"github.com/go-logr/logr"
	"k8s.io/client-go/kubernetes"
)

// ServiceConfig is the configuration for a Service object.
type ServiceConfig struct {
	KubeClient    kubernetes.Interface
	Log           logr.Logger
	MaxRetryCount int
	MaxRetryTime  time.Duration
	MinRetryTime  time.Duration
	WithRetries   bool
}

// Service provides aws-auth configmap management behavior.
type Service interface {
	// UpsertMapRole upserts a MapRole into the configmap keyed by username.
	UpsertMapRole(username string, mapRole MapRole) error

	// RemoveMapRole removes a MapRole from the configmap by keyed by username
	RemoveMapRole(username string) error

	// UpsertMapUser upserts a MapUser into the configmap keyed by username.
	UpsertMapUser(username string, mapUser MapUser) error

	// RemoveMapUser removes a MapUser from the configmap keyed by username
	RemoveMapUser(username string) error
}

// NewService returns an implementation of the Service interface.
func NewService(cfg *ServiceConfig) (Service, error) {
	if cfg.WithRetries {
		if cfg.MaxRetryCount < 1 {
			return nil, errors.New("retry max count config must be greater than zero")
		}
	}
	return impl{cfg: *cfg}, nil
}

type impl struct {
	cfg ServiceConfig
}

// UpsertMapRole upserts a MapRole into the configmap keyed by username.
func (svc impl) UpsertMapRole(username string, mapRole MapRole) error {
	mapper := NewMapper(svc.cfg.KubeClient, false)
	err := mapper.Upsert(&Arguments{
		DataType:      MapRoleData,
		RoleARN:       mapRole.RoleARN,
		Username:      username,
		Groups:        mapRole.Groups,
		WithRetries:   svc.cfg.WithRetries,
		MaxRetryCount: svc.cfg.MaxRetryCount,
		MaxRetryTime:  svc.cfg.MaxRetryTime,
		MinRetryTime:  svc.cfg.MinRetryTime,
	})
	if err != nil {
		svc.cfg.Log.Error(err, "failure to upsert mapRole", "username", username)
	}
	return err
}

// RemoveMapRole removes a MapRole from the configmap keyed by username.
func (svc impl) RemoveMapRole(username string) error {
	mapper := NewMapper(svc.cfg.KubeClient, false)
	err := mapper.Remove(&Arguments{
		DataType:      MapRoleData,
		Username:      username,
		WithRetries:   svc.cfg.WithRetries,
		MaxRetryCount: svc.cfg.MaxRetryCount,
		MaxRetryTime:  svc.cfg.MaxRetryTime,
		MinRetryTime:  svc.cfg.MinRetryTime,
	})
	if err != nil {
		svc.cfg.Log.Info("mapRole not found", "username", username)
	}
	return err
}

// UpsertMapUser upserts a MapUser into the configmap keyed by username.
func (svc impl) UpsertMapUser(username string, mapUser MapUser) error {
	mapper := NewMapper(svc.cfg.KubeClient, false)
	err := mapper.Upsert(&Arguments{
		DataType:      MapUserData,
		UserARN:       mapUser.UserARN,
		Username:      username,
		Groups:        mapUser.Groups,
		WithRetries:   svc.cfg.WithRetries,
		MaxRetryCount: svc.cfg.MaxRetryCount,
		MaxRetryTime:  svc.cfg.MaxRetryTime,
		MinRetryTime:  svc.cfg.MinRetryTime,
	})
	if err != nil {
		svc.cfg.Log.Error(err, "failure to upsert mapUser", "username", username)
	}
	return err
}

// RemoveMapUser removes a MapUser from the configmap keyed by username.
func (svc impl) RemoveMapUser(username string) error {
	mapper := NewMapper(svc.cfg.KubeClient, false)
	err := mapper.Remove(&Arguments{
		DataType:      MapUserData,
		Username:      username,
		WithRetries:   svc.cfg.WithRetries,
		MaxRetryCount: svc.cfg.MaxRetryCount,
		MaxRetryTime:  svc.cfg.MaxRetryTime,
		MinRetryTime:  svc.cfg.MinRetryTime,
	})
	if err != nil {
		svc.cfg.Log.Info("mapUser not found", "username", username)
	}
	return err
}
