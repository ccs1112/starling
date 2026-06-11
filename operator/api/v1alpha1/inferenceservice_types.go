package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// InferenceServiceSpec is the desired state. It answers
// "which model image to serve and how many replicas?". It is a deliberately tiny
// replacement for KServe's InferenceService. The control plane mechanisms are identical
// whether the pod runs vLLM or `sleep`, so the model is a placeholder.
type InferenceServiceSpec struct {
	// Image is the container that serves the model.
	Image string `json:"image"`

	// Replicas is the desired number of serving pods.
	// The controller drives an owned Deployment toward this count.
	// The gang scheduler admits them all-or-nothing.
	Replicas int32 `json:"replicas"`
}

// InferenceServiceStatus is the observed state the controller writes back.
// It is kept separate from Spec so that the status subresource can be updated
// without racing a user's edit to Spec.
type InferenceServiceStatus struct {
	// ReadyReplicas mirrors the owned Deployment's ready count.
	ReadyReplicas int32 `json:"readyReplicas,omitempty"`
}

// InferenceService is one served model. The API server stores it. The controller
// reconciles a Deployment and Service to match.
type InferenceService struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   InferenceServiceSpec   `json:"spec,omitempty"`
	Status InferenceServiceStatus `json:"status,omitempty"`
}

// InferenceServiceList is what a List call returns. It is a page of items plus
// the list metadata (resourceVersion, continue token) the informer needs to Watch
// from the right point.
type InferenceServiceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []InferenceService `json:"items"`
}
