package controllers

import (
	"context"
	"time"

	"github.com/elrondwong/redis-operator/api/status"
	redisv1beta2 "github.com/elrondwong/redis-operator/api/v1beta2"
	"github.com/elrondwong/redis-operator/k8sutils"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// RedisSentinelReconciler reconciles a RedisSentinel object
type RedisSentinelReconciler struct {
	client.Client
	K8sClient  kubernetes.Interface
	Dk8sClinet dynamic.Interface
	Log        logr.Logger
	Scheme     *runtime.Scheme
}

// Reconcile is part of the main kubernetes reconciliation loop which aims
func (r *RedisSentinelReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	reqLogger := r.Log.WithValues("Request.Namespace", req.Namespace, "Request.Name", req.Name)
	reqLogger.Info("Reconciling opstree redis controller")
	instance := &redisv1beta2.RedisSentinel{}

	err := r.Client.Get(context.TODO(), req.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	if _, found := instance.ObjectMeta.GetAnnotations()["redissentinel.opstreelabs.in/skip-reconcile"]; found {
		reqLogger.Info("Found annotations redissentinel.opstreelabs.in/skip-reconcile, so skipping reconcile")
		return ctrl.Result{RequeueAfter: time.Second * 10}, nil
	}

	// Get total Sentinel Replicas
	// sentinelReplicas := instance.Spec.GetSentinelCounts("sentinel")

	if err = k8sutils.HandleRedisSentinelFinalizer(instance, r.Client); err != nil {
		return ctrl.Result{RequeueAfter: time.Second * 60}, err
	}

	if err = k8sutils.AddRedisSentinelFinalizer(instance, r.Client); err != nil {
		return ctrl.Result{RequeueAfter: time.Second * 60}, err
	}

	if instance.Status.State != status.RedisSenitnelReady {
		err = k8sutils.UpdateRedisSentinelStatus(instance, status.RedisSentinelInitializing, status.InitializingSentinelReason)
		if err != nil {
			return ctrl.Result{RequeueAfter: time.Second * 10}, err
		}
	}

	// Create Redis Sentinel
	err = k8sutils.CreateRedisSentinel(ctx, r.K8sClient, r.Log, instance)
	if err != nil {
		return ctrl.Result{}, err
	}

	err = k8sutils.ReconcileSentinelPodDisruptionBudget(instance, instance.Spec.PodDisruptionBudget)
	if err != nil {
		return ctrl.Result{}, err
	}

	// Create the Service for Redis Sentinel
	err = k8sutils.CreateRedisSentinelService(instance)
	if err != nil {
		return ctrl.Result{}, err
	}

	if instance.Status.State != status.RedisSenitnelReady && k8sutils.CheckRedisSentinelReady(instance) {
		err = k8sutils.UpdateRedisSentinelStatus(instance, status.RedisSenitnelReady, status.ReadySentinelReason)
		if err != nil {
			return ctrl.Result{RequeueAfter: time.Second * 10}, err
		}
	}

	reqLogger.Info("Will reconcile redis operator in again 10 seconds")
	return ctrl.Result{RequeueAfter: time.Second * 10}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *RedisSentinelReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&redisv1beta2.RedisSentinel{}).
		Complete(r)
}
