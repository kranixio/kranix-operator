package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// KranixNamespaceSpec defines the desired state of KranixNamespace
type KranixNamespaceSpec struct {
	// Labels are labels to apply to the namespace
	Labels map[string]string `json:"labels,omitempty"`

	// ResourceQuota defines resource limits for the namespace
	ResourceQuota *ResourceQuota `json:"resourceQuota,omitempty"`
}

// ResourceQuota defines resource quota
type ResourceQuota struct {
	CPU    string `json:"cpu,omitempty"`
	Memory string `json:"memory,omitempty"`
}

// KranixNamespaceStatus defines the observed state of KranixNamespace
type KranixNamespaceStatus struct {
	// Phase is the current phase
	Phase string `json:"phase,omitempty"`

	// LastReconciled is the timestamp of the last reconciliation
	LastReconciled *metav1.Time `json:"lastReconciled,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName=krnens

// KranixNamespace is the Schema for the kranixnamespaces API
type KranixNamespace struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KranixNamespaceSpec   `json:"spec,omitempty"`
	Status KranixNamespaceStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// KranixNamespaceList contains a list of KranixNamespace
type KranixNamespaceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KranixNamespace `json:"items"`
}

// DeepCopyObject implements runtime.Object
func (in *KranixNamespace) DeepCopyObject() runtime.Object {
	if in == nil {
		return nil
	}
	out := &KranixNamespace{}
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto copies from in into out
func (in *KranixNamespace) DeepCopyInto(out *KranixNamespace) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	out.Status = in.Status
}

// DeepCopyObject implements runtime.Object
func (in *KranixNamespaceList) DeepCopyObject() runtime.Object {
	if in == nil {
		return nil
	}
	out := &KranixNamespaceList{}
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto copies from in into out
func (in *KranixNamespaceList) DeepCopyInto(out *KranixNamespaceList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		out.Items = make([]KranixNamespace, len(in.Items))
		for i := range in.Items {
			in.Items[i].DeepCopyInto(&out.Items[i])
		}
	}
}

func init() {
	SchemeBuilder.Register(&KranixNamespace{}, &KranixNamespaceList{})
}
