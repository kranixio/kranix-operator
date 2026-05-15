package controllers

import (
	"context"

	"github.com/kranix-io/kranix-operator/api/v1alpha1"
	"github.com/kranix-io/kranix-operator/internal/reconciler"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// KranixPolicyReconciler reconciles a KranixPolicy object
type KranixPolicyReconciler struct {
	client.Client
	Scheme    *runtime.Scheme
	Reconciler *reconciler.Reconciler
}

// +kubebuilder:rbac:groups=kranix.io,resources=kranixpolicies,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=kranix.io,resources=kranixpolicies/status,verbs=get;update;patch

// Reconcile is the reconciliation loop for KranixPolicy
func (r *KranixPolicyReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// Fetch the KranixPolicy instance
	var kranixPolicy v1alpha1.KranixPolicy
	if err := r.Get(ctx, req.NamespacedName, &kranixPolicy); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Delegate to the reconciler
	result, err := r.Reconciler.ReconcilePolicy(ctx, &kranixPolicy)
	if err != nil {
		logger.Error(err, "Failed to reconcile KranixPolicy")
		return ctrl.Result{}, err
	}

	return result, nil
}

// SetupWithManager sets up the controller with the Manager
func (r *KranixPolicyReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.KranixPolicy{}).
		Complete(r)
}
