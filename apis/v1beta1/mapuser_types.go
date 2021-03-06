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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// MapUserSpec defines the desired state of MapUser
type MapUserSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// The User ARN to associate with the MapUser
	UserARN string `json:"userarn"`

	// The Kubernetes groups to associate with the MapUser
	// +kubebuilder:validation:Optional
	Groups []string `json:"groups"`

	// A useful description of the MapUser
	// +kubebuilder:validation:Optional
	Description string `json:"description"`

	// The email address of a contact person for the MapUser
	// +kubebuilder:validation:Optional
	Email string `json:"email"`
}

// MapUserStatus defines the observed state of MapUser
type MapUserStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:resource:scope=Cluster
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="User ARN",type=string,JSONPath=`.spec.userarn`
//+kubebuilder:printcolumn:name="Groups",type=string,JSONPath=`.spec.groups`
//+kubebuilder:printcolumn:name="Email",type=string,JSONPath=`.spec.email`
//+kubebuilder:printcolumn:name="Description",type=string,JSONPath=`.spec.description`

// MapUser is the Schema for the users API
type MapUser struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MapUserSpec   `json:"spec,omitempty"`
	Status MapUserStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// MapUserList contains a list of MapUser
type MapUserList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MapUser `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MapUser{}, &MapUserList{})
}
