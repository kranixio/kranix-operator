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

// KranixAppReconciler reconciles a KranixApp object
type KranixAppReconciler struct {
	client.Client
	Scheme    *runtime.Scheme
	Reconciler *reconciler.Reconciler
}

// +kubebuilder:rbac:groups=kranix.io,resources=kranixapps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=kranix.io,resources=kranixapps/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=kranix.io,resources=kranixapps/finalizers,verbs=update

// Reconcile is the reconciliation loop for KranixApp
func (r *KranixAppReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// Fetch the KranixApp instance
	var kranixApp v1alpha1.KranixApp
	if err := r.Get(ctx, req.NamespacedName, &kranixApp); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Delegate to the reconciler
	result, err := r.Reconciler.ReconcileApp(ctx, &kranixApp)
	if err != nil {
		logger.Error(err, "Failed to reconcile KranixApp")
		return ctrl.Result{}, err
	}

	return result, nil
}

// SetupWithManager sets up the controller with the Manager
func (r *KranixAppReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.KranixApp{}).
		Complete(r)
}
