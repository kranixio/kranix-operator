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

// KranixNamespaceReconciler reconciles a KranixNamespace object
type KranixNamespaceReconciler struct {
	client.Client
	Scheme    *runtime.Scheme
	Reconciler *reconciler.Reconciler
}

// +kubebuilder:rbac:groups=kranix.io,resources=kranixnamespaces,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=kranix.io,resources=kranixnamespaces/status,verbs=get;update;patch

// Reconcile is the reconciliation loop for KranixNamespace
func (r *KranixNamespaceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// Fetch the KranixNamespace instance
	var kranixNamespace v1alpha1.KranixNamespace
	if err := r.Get(ctx, req.NamespacedName, &kranixNamespace); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Delegate to the reconciler
	result, err := r.Reconciler.ReconcileNamespace(ctx, &kranixNamespace)
	if err != nil {
		logger.Error(err, "Failed to reconcile KranixNamespace")
		return ctrl.Result{}, err
	}

	return result, nil
}

// SetupWithManager sets up the controller with the Manager
func (r *KranixNamespaceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.KranixNamespace{}).
		Complete(r)
}
