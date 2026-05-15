package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// KranixAppSpec defines the desired state of KranixApp
type KranixAppSpec struct {
	// Image is the container image to deploy
	Image string `json:"image"`

	// Replicas is the number of desired replicas
	Replicas int32 `json:"replicas"`

	// Namespace is the target namespace for the workload
	Namespace string `json:"namespace"`

	// Env is environment variables for the container
	Env map[string]string `json:"env,omitempty"`

	// Resources defines CPU and memory limits
	Resources ResourceRequirements `json:"resources,omitempty"`

	// Ports are container ports to expose
	Ports []corev1.ContainerPort `json:"ports,omitempty"`

	// Rollout defines deployment strategy
	Rollout RolloutConfig `json:"rollout,omitempty"`

	// AutoHeal enables automatic remediation
	AutoHeal bool `json:"autoHeal,omitempty"`
}

// KranixAppStatus defines the observed state of KranixApp
type KranixAppStatus struct {
	// Phase is the current phase of the workload
	Phase AppPhase `json:"phase"`

	// ReadyReplicas is the number of ready replicas
	ReadyReplicas int32 `json:"readyReplicas,omitempty"`

	// LastReconciled is the timestamp of the last reconciliation
	LastReconciled *metav1.Time `json:"lastReconciled,omitempty"`

	// Conditions describe the current conditions
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// ResourceRequirements defines resource limits
type ResourceRequirements struct {
	CPU    string `json:"cpu,omitempty"`
	Memory string `json:"memory,omitempty"`
}

// RolloutConfig defines rollout strategy
type RolloutConfig struct {
	Strategy          string `json:"strategy,omitempty"`
	MaxUnavailable    int32  `json:"maxUnavailable,omitempty"`
	HealthCheckPath   string `json:"healthCheckPath,omitempty"`
	HealthCheckTimeout string `json:"healthCheckTimeout,omitempty"`
}

// AppPhase represents the phase of a KranixApp
type AppPhase string

const (
	AppPhasePending   AppPhase = "Pending"
	AppPhaseDeploying AppPhase = "Deploying"
	AppPhaseRunning   AppPhase = "Running"
	AppPhaseDegraded  AppPhase = "Degraded"
	AppPhaseFailed    AppPhase = "Failed"
)

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName=krane

// KranixApp is the Schema for the kranixapps API
type KranixApp struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KranixAppSpec   `json:"spec,omitempty"`
	Status KranixAppStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// KranixAppList contains a list of KranixApp
type KranixAppList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KranixApp `json:"items"`
}

// DeepCopyObject implements runtime.Object
func (in *KranixApp) DeepCopyObject() runtime.Object {
	if in == nil {
		return nil
	}
	out := &KranixApp{}
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto copies from in into out
func (in *KranixApp) DeepCopyInto(out *KranixApp) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	out.Status = in.Status
}

// DeepCopyObject implements runtime.Object
func (in *KranixAppList) DeepCopyObject() runtime.Object {
	if in == nil {
		return nil
	}
	out := &KranixAppList{}
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto copies from in into out
func (in *KranixAppList) DeepCopyInto(out *KranixAppList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		out.Items = make([]KranixApp, len(in.Items))
		for i := range in.Items {
			in.Items[i].DeepCopyInto(&out.Items[i])
		}
	}
}

func init() {
	SchemeBuilder.Register(&KranixApp{}, &KranixAppList{})
}
