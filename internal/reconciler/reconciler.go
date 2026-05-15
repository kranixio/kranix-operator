package reconciler

import (
	"context"
	"fmt"
	"time"

	"github.com/kranix-io/kranix-operator/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

// Reconciler handles reconciliation logic for all Kranix resources
type Reconciler struct {
	CoreClient CoreClient
}

// CoreClient is the interface for calling kranix-core
type CoreClient interface {
	DeployWorkload(ctx context.Context, spec *WorkloadSpec) (*WorkloadStatus, error)
	GetWorkloadStatus(ctx context.Context, name, namespace string) (*WorkloadStatus, error)
	CreateNamespace(ctx context.Context, name string, labels map[string]string, quota *ResourceQuota) error
	ApplyPolicy(ctx context.Context, namespace string, policy *PolicySpec) error
}

// WorkloadSpec is the spec sent to kranix-core
type WorkloadSpec struct {
	Name      string
	Namespace string
	Image     string
	Replicas  int32
	Env       map[string]string
	CPU       string
	Memory    string
	Ports     []Port
}

// Port represents a container port
type Port struct {
	ContainerPort int32
	Protocol      string
}

// WorkloadStatus is the status from kranix-core
type WorkloadStatus struct {
	Phase         string
	ReadyReplicas int32
}

// ResourceQuota represents resource quota
type ResourceQuota struct {
	CPU    string
	Memory string
}

// PolicySpec represents a policy
type PolicySpec struct {
	EnforceResourceLimits bool
	DefaultCpuLimit       string
	DefaultMemoryLimit    string
	AllowPrivileged       bool
	NetworkPolicy         *NetworkPolicy
}

// NetworkPolicy represents network policy
type NetworkPolicy struct {
	IngressFrom []string
}

// NewReconciler creates a new reconciler
func NewReconciler(coreClient CoreClient) *Reconciler {
	return &Reconciler{
		CoreClient: coreClient,
	}
}

// ReconcileApp reconciles a KranixApp
func (r *Reconciler) ReconcileApp(ctx context.Context, app *v1alpha1.KranixApp) (ctrl.Result, error) {
	// Build workload spec from KranixApp spec
	spec := &WorkloadSpec{
		Name:      app.Name,
		Namespace: app.Spec.Namespace,
		Image:     app.Spec.Image,
		Replicas:  app.Spec.Replicas,
		Env:       app.Spec.Env,
		CPU:       app.Spec.Resources.CPU,
		Memory:    app.Spec.Resources.Memory,
	}

	for _, p := range app.Spec.Ports {
		spec.Ports = append(spec.Ports, Port{
			ContainerPort: p.ContainerPort,
			Protocol:      string(p.Protocol),
		})
	}

	// Call kranix-core to deploy
	status, err := r.CoreClient.DeployWorkload(ctx, spec)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to deploy workload: %w", err)
	}

	// Update status
	now := metav1.Now()
	app.Status.Phase = v1alpha1.AppPhase(status.Phase)
	app.Status.ReadyReplicas = status.ReadyReplicas
	app.Status.LastReconciled = &now

	// Set condition
	if status.Phase == string(v1alpha1.AppPhaseRunning) {
		setCondition(&app.Status, "Ready", "True")
	} else {
		setCondition(&app.Status, "Ready", "False")
	}

	// Auto-heal if enabled and workload is degraded
	if app.Spec.AutoHeal && app.Status.Phase == v1alpha1.AppPhaseDegraded {
		// Trigger remediation via core
		// TODO: implement remediation logic
	}

	return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
}

// ReconcileNamespace reconciles a KranixNamespace
func (r *Reconciler) ReconcileNamespace(ctx context.Context, ns *v1alpha1.KranixNamespace) (ctrl.Result, error) {
	var quota *ResourceQuota
	if ns.Spec.ResourceQuota != nil {
		quota = &ResourceQuota{
			CPU:    ns.Spec.ResourceQuota.CPU,
			Memory: ns.Spec.ResourceQuota.Memory,
		}
	}

	err := r.CoreClient.CreateNamespace(ctx, ns.Name, ns.Spec.Labels, quota)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to create namespace: %w", err)
	}

	now := metav1.Now()
	ns.Status.Phase = "Ready"
	ns.Status.LastReconciled = &now

	return ctrl.Result{RequeueAfter: 5 * time.Minute}, nil
}

// ReconcilePolicy reconciles a KranixPolicy
func (r *Reconciler) ReconcilePolicy(ctx context.Context, policy *v1alpha1.KranixPolicy) (ctrl.Result, error) {
	policySpec := &PolicySpec{
		EnforceResourceLimits: policy.Spec.EnforceResourceLimits,
		DefaultCpuLimit:       policy.Spec.DefaultCpuLimit,
		DefaultMemoryLimit:    policy.Spec.DefaultMemoryLimit,
		AllowPrivileged:       policy.Spec.AllowPrivileged,
	}

	if policy.Spec.NetworkPolicy != nil {
		policySpec.NetworkPolicy = &NetworkPolicy{
			IngressFrom: policy.Spec.NetworkPolicy.IngressFrom,
		}
	}

	err := r.CoreClient.ApplyPolicy(ctx, policy.Namespace, policySpec)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to apply policy: %w", err)
	}

	now := metav1.Now()
	policy.Status.Phase = "Applied"
	policy.Status.LastReconciled = &now

	return ctrl.Result{RequeueAfter: 5 * time.Minute}, nil
}

func setCondition(status *v1alpha1.KranixAppStatus, condType string, statusStr string) {
	condition := metav1.Condition{
		Type:               condType,
		Status:             metav1.ConditionStatus(statusStr),
		LastTransitionTime: metav1.Now(),
	}

	for i, c := range status.Conditions {
		if c.Type == condType {
			status.Conditions[i] = condition
			return
		}
	}

	status.Conditions = append(status.Conditions, condition)
}
