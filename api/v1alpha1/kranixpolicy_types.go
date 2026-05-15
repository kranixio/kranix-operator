package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// KranixPolicySpec defines the desired state of KranixPolicy
type KranixPolicySpec struct {
	// EnforceResourceLimits enables resource limit enforcement
	EnforceResourceLimits bool `json:"enforceResourceLimits,omitempty"`

	// DefaultCpuLimit is the default CPU limit for workloads
	DefaultCpuLimit string `json:"defaultCpuLimit,omitempty"`

	// DefaultMemoryLimit is the default memory limit for workloads
	DefaultMemoryLimit string `json:"defaultMemoryLimit,omitempty"`

	// AllowPrivileged allows privileged containers
	AllowPrivileged bool `json:"allowPrivileged,omitempty"`

	// NetworkPolicy defines network policy rules
	NetworkPolicy *NetworkPolicy `json:"networkPolicy,omitempty"`
}

// NetworkPolicy defines network policy
type NetworkPolicy struct {
	IngressFrom []string `json:"ingressFrom,omitempty"`
}

// KranixPolicyStatus defines the observed state of KranixPolicy
type KranixPolicyStatus struct {
	// Phase is the current phase
	Phase string `json:"phase,omitempty"`

	// LastReconciled is the timestamp of the last reconciliation
	LastReconciled *metav1.Time `json:"lastReconciled,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName=krnpol

// KranixPolicy is the Schema for the kranixpolicies API
type KranixPolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KranixPolicySpec   `json:"spec,omitempty"`
	Status KranixPolicyStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// KranixPolicyList contains a list of KranixPolicy
type KranixPolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KranixPolicy `json:"items"`
}

// DeepCopyObject implements runtime.Object
func (in *KranixPolicy) DeepCopyObject() runtime.Object {
	if in == nil {
		return nil
	}
	out := &KranixPolicy{}
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto copies from in into out
func (in *KranixPolicy) DeepCopyInto(out *KranixPolicy) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	out.Status = in.Status
}

// DeepCopyObject implements runtime.Object
func (in *KranixPolicyList) DeepCopyObject() runtime.Object {
	if in == nil {
		return nil
	}
	out := &KranixPolicyList{}
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto copies from in into out
func (in *KranixPolicyList) DeepCopyInto(out *KranixPolicyList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		out.Items = make([]KranixPolicy, len(in.Items))
		for i := range in.Items {
			in.Items[i].DeepCopyInto(&out.Items[i])
		}
	}
}

func init() {
	SchemeBuilder.Register(&KranixPolicy{}, &KranixPolicyList{})
}
