/*
Copyright 2023 yangsijie666.
*/

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// EchoSpec defines the desired state of Echo
type EchoSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// SaySomeThing
	SaySomeThing string `json:"saySomeThing"`
}

// EchoStatus defines the observed state of Echo
type EchoStatus struct {
	// EchoResult is equal to spec.SaySomeThing
	EchoResult string `json:"echoResult"`

	// ObservedGeneration is the most recent generation observed for this StatefulSet. It corresponds to the
	// ScaleTask's generation, which is updated on mutation by the API Server.
	//+kubebuilder:default=0
	ObservedGeneration int64 `json:"observedGeneration"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Namespaced
//+kubebuilder:printcolumn:name="EchoResult",type=string,JSONPath=`.status.echoResult`

// Echo is the Schema for the echoes API
type Echo struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EchoSpec   `json:"spec,omitempty"`
	Status EchoStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// EchoList contains a list of Echo
type EchoList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Echo `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Echo{}, &EchoList{})
}
