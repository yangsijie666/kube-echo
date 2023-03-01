/*
Copyright 2023 yangsijie666.
*/

package controllers

import (
	"context"
	"fmt"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/uuid"
	"k8s.io/client-go/util/retry"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	gogogov1 "kube-echo/api/v1"
)

// EchoReconciler reconciles a Echo object
type EchoReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=gogogo.yangsijie666.github.com,resources=echoes,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=gogogo.yangsijie666.github.com,resources=echoes/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=gogogo.yangsijie666.github.com,resources=echoes/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Echo object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.10.0/pkg/reconcile
func (r *EchoReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	l := log.FromContext(ctx, "Reconcile", types.NamespacedName{
		Namespace: req.Namespace,
		Name:      req.Name,
	}, "traceId", uuid.NewUUID())
	ctx = log.IntoContext(ctx, l)

	instance := &gogogov1.Echo{}
	if err := r.Get(ctx, types.NamespacedName{
		Namespace: req.Namespace,
		Name:      req.Name,
	}, instance); err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		l.Error(err, "get instance error")
		return ctrl.Result{}, err
	}

	// handle deletion
	if instance.DeletionTimestamp != nil {
		l.V(1).Info("instance is deleting")
		return ctrl.Result{}, nil
	}

	// do echo
	newStatus := instance.Status.DeepCopy()
	newStatus.EchoResult = instance.Spec.SaySomeThing

	return r.updateStatus(ctx, instance, newStatus)
}

func (r *EchoReconciler) updateStatus(ctx context.Context, instance *gogogov1.Echo, newStatus *gogogov1.EchoStatus) (ctrl.Result, error) {
	updated := instance.Status.EchoResult != newStatus.EchoResult
	if !updated {
		log.FromContext(ctx).V(1).Info("status remain unchanged")
		return ctrl.Result{}, nil
	}

	retryCnt := 1
	if err := retry.RetryOnConflict(retry.DefaultBackoff, func() error {
		defer func() {
			retryCnt += 1
		}()

		log.FromContext(ctx).V(1).Info(fmt.Sprintf("The %d th time updating status for %v: %s/%s, ", retryCnt, instance.Kind, instance.Namespace, instance.Name) +
			fmt.Sprintf("echoResult %s->%s, ", instance.Status.EchoResult, newStatus.EchoResult))

		obj := &gogogov1.Echo{}
		if err := r.Client.Get(context.TODO(), client.ObjectKey{Namespace: instance.Namespace, Name: instance.Name}, obj); err != nil {
			return err
		}

		obj.Status = *newStatus
		obj.Status.ObservedGeneration = obj.Generation
		if err := r.Client.Status().Update(context.TODO(), obj); err != nil {
			return err
		}
		return nil
	}); err != nil {
		log.FromContext(ctx).Error(err, "update status error")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *EchoReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&gogogov1.Echo{}).
		Complete(r)
}
