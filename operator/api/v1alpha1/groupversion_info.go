// Package v1alpha1 contains the InferenceService API types.
// This is just enough type machinery for the API server to register
// the kind and for a typed client to encode and decode it.
package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// GroupVersion is the group/version this package registers. The API server
// exposes these types under v1alpha1.
var GroupVersion = schema.GroupVersion{Group: "inference.noetics.dev", Version: "v1alpha1"}

// SchemeBuilder collects the functions that add our types to a runtime.Scheme.
// client-go's codecs use the Scheme to map Go structs to and from wire JSON.
// Without this registration the client cannot encode an InferenceService.
var SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)

// AddToScheme registers our types into a Scheme. main() calls this once at
// startup so that the typed client knows InferenceService and its list.
var AddToScheme = SchemeBuilder.AddToScheme

func addKnownTypes(scheme *runtime.Scheme) error {
	// Register the object and its list under the group and version.
	scheme.AddKnownTypes(GroupVersion,
		&InferenceService{},
		&InferenceServiceList{},
	)
	// Register the shared metav1 types like ListOption for this
	// GroupVersion so that the client can issue List and Watch with
	// the standard options against our group.
	metav1.AddToGroupVersion(scheme, GroupVersion)
	return nil
}
