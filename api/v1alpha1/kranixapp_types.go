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
	Strategy           string `json:"strategy,omitempty"`
	MaxUnavailable     int32  `json:"maxUnavailable,omitempty"`
	HealthCheckPath    string `json:"healthCheckPath,omitempty"`
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
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		out.Items = make([]KranixApp, len(in.Items))
		for i := range in.Items {
			in.Items[i].DeepCopyInto(&out.Items[i])
		}
	}
}

// AutoScalingConfig defines auto-scaling behavior.
type AutoScalingConfig struct {
	Enabled                  bool           `json:"enabled"`
	MinReplicas              int32          `json:"minReplicas"`
	MaxReplicas              int32          `json:"maxReplicas"`
	TargetCPUUtilization     int32          `json:"targetCPUUtilization,omitempty"`    // percentage
	TargetMemoryUtilization  int32          `json:"targetMemoryUtilization,omitempty"` // percentage
	CustomMetrics            []CustomMetric `json:"customMetrics,omitempty"`
	ScaleDownCooldownSeconds int32          `json:"scaleDownCooldownSeconds,omitempty"`
	ScaleUpCooldownSeconds   int32          `json:"scaleUpCooldownSeconds,omitempty"`
}

// CustomMetric defines a custom metric for auto-scaling.
type CustomMetric struct {
	Name       string       `json:"name"`
	Type       string       `json:"type"` // pods, object
	MetricName string       `json:"metricName"`
	Target     MetricTarget `json:"target"`
}

// MetricTarget defines the target value for a metric.
type MetricTarget struct {
	Type         string `json:"type"` // average, value
	AverageValue string `json:"averageValue,omitempty"`
	Value        string `json:"value,omitempty"`
}

// SchedulingConfig defines scheduling preferences.
type SchedulingConfig struct {
	CostAware        bool              `json:"costAware,omitempty"`
	PreferredRegions []string          `json:"preferredRegions,omitempty"`
	PreferredZones   []string          `json:"preferredZones,omitempty"`
	NodeSelectors    map[string]string `json:"nodeSelectors,omitempty"`
	Affinity         *AffinityConfig   `json:"affinity,omitempty"`
	Tolerations      []Toleration      `json:"tolerations,omitempty"`
	MaxCostPerHour   string            `json:"maxCostPerHour,omitempty"`
}

// AffinityConfig defines pod affinity/anti-affinity rules.
type AffinityConfig struct {
	NodeAffinity    *NodeAffinity `json:"nodeAffinity,omitempty"`
	PodAffinity     *PodAffinity  `json:"podAffinity,omitempty"`
	PodAntiAffinity *PodAffinity  `json:"podAntiAffinity,omitempty"`
}

// NodeAffinity defines node affinity rules.
type NodeAffinity struct {
	RequiredDuringScheduling  []NodeSelectorTerm        `json:"requiredDuringScheduling,omitempty"`
	PreferredDuringScheduling []PreferredSchedulingTerm `json:"preferredDuringScheduling,omitempty"`
}

// NodeSelectorTerm defines a node selector term.
type NodeSelectorTerm struct {
	MatchExpressions []NodeSelectorRequirement `json:"matchExpressions,omitempty"`
	MatchFields      []NodeSelectorRequirement `json:"matchFields,omitempty"`
}

// NodeSelectorRequirement defines a node selector requirement.
type NodeSelectorRequirement struct {
	Key      string   `json:"key"`
	Operator string   `json:"operator"` // In, NotIn, Exists, DoesNotExist, Gt, Lt
	Values   []string `json:"values,omitempty"`
}

// PreferredSchedulingTerm defines a preferred scheduling term.
type PreferredSchedulingTerm struct {
	Weight     int32            `json:"weight"`
	Preference NodeSelectorTerm `json:"preference"`
}

// PodAffinity defines pod affinity rules.
type PodAffinity struct {
	RequiredDuringScheduling  []PodAffinityTerm         `json:"requiredDuringScheduling,omitempty"`
	PreferredDuringScheduling []WeightedPodAffinityTerm `json:"preferredDuringScheduling,omitempty"`
}

// PodAffinityTerm defines a pod affinity term.
type PodAffinityTerm struct {
	LabelSelector map[string]string `json:"labelSelector,omitempty"`
	Namespaces    []string          `json:"namespaces,omitempty"`
	TopologyKey   string            `json:"topologyKey"`
}

// WeightedPodAffinityTerm defines a weighted pod affinity term.
type WeightedPodAffinityTerm struct {
	Weight          int32           `json:"weight"`
	PodAffinityTerm PodAffinityTerm `json:"podAffinityTerm"`
}

// Toleration defines a toleration for taints.
type Toleration struct {
	Key               string `json:"key,omitempty"`
	Operator          string `json:"operator,omitempty"` // Exists, Equal
	Value             string `json:"value,omitempty"`
	Effect            string `json:"effect,omitempty"` // NoSchedule, PreferNoSchedule, NoExecute
	TolerationSeconds *int64 `json:"tolerationSeconds,omitempty"`
}

// CanaryConfig defines canary deployment configuration.
type CanaryConfig struct {
	Replicas         int32    `json:"replicas"`
	Percentage       int32    `json:"percentage,omitempty"`
	AnalysisDuration string   `json:"analysisDuration,omitempty"`
	SuccessThreshold int32    `json:"successThreshold,omitempty"`
	Metrics          []string `json:"metrics,omitempty"`
	AutoPromote      bool     `json:"autoPromote,omitempty"`
}

// ABTestConfig defines A/B testing configuration.
type ABTestConfig struct {
	VariantA         string   `json:"variantA"`
	VariantB         string   `json:"variantB"`
	TrafficSplit     int32    `json:"trafficSplit"` // percentage for variant B
	AnalysisDuration string   `json:"analysisDuration,omitempty"`
	Metrics          []string `json:"metrics,omitempty"`
	AutoSelectWinner bool     `json:"autoSelectWinner,omitempty"`
}

func init() {
	SchemeBuilder.Register(&KranixApp{}, &KranixAppList{})
}
