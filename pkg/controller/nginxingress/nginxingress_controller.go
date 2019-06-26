package nginxingress

import (
	"context"
	"os"
	"reflect"

	raven "github.com/getsentry/raven-go"
	appv1alpha1 "github.com/tekliner/nginx-ingress-operator/pkg/apis/app/v1alpha1"
	"k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_nginxingress")

func init() {
	if os.Getenv("SENTRY_DSN") != "" {
		raven.SetDSN(os.Getenv("SENTRY_DSN"))
	}
}

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new NginxIngress Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileNginxIngress{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("nginxingress-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource NginxIngress
	err = c.Watch(&source.Kind{Type: &appv1alpha1.NginxIngress{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner NginxIngress
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &appv1alpha1.NginxIngress{},
	})
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileNginxIngress{}

// ReconcileNginxIngress reconciles a NginxIngress object
type ReconcileNginxIngress struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a NginxIngress object and makes changes based on the state read
// and what is in the NginxIngress.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileNginxIngress) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling NginxIngress")

	// Fetch the NginxIngress instance
	instance := &appv1alpha1.NginxIngress{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		raven.CaptureErrorAndWait(err, nil)
		return reconcile.Result{}, err
	}

	// reconcile deployment
	newDeployment := generateDeployment(instance)

	if err := controllerutil.SetControllerReference(instance, &newDeployment, r.scheme); err != nil {
		raven.CaptureErrorAndWait(err, nil)
		return reconcile.Result{}, err
	}

	foundDeployment := &v1.Deployment{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: newDeployment.Name, Namespace: newDeployment.Namespace}, foundDeployment)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new Deployment", "Namespace", newDeployment.Namespace, "Name", newDeployment.Name)
		err = r.client.Create(context.TODO(), &newDeployment)
		if err != nil {
			raven.CaptureErrorAndWait(err, nil)
			return reconcile.Result{}, err
		}
	} else if err != nil {
		raven.CaptureErrorAndWait(err, nil)
		return reconcile.Result{}, err
	}

	reconcileDeployment := false

	if !reflect.DeepEqual(foundDeployment.Spec, newDeployment.Spec) {
		foundDeployment.Spec.Replicas = newDeployment.Spec.Replicas
		foundDeployment.Spec.Template = newDeployment.Spec.Template
		reconcileDeployment = true
	}

	if !reflect.DeepEqual(foundDeployment.Annotations, newDeployment.Annotations) {
		foundDeployment.Annotations = newDeployment.Annotations
		reconcileDeployment = true
	}

	if !reflect.DeepEqual(foundDeployment.Labels, newDeployment.Labels) {
		foundDeployment.Labels = newDeployment.Labels
		reconcileDeployment = true
	}

	if reconcileDeployment {
		if err = r.client.Update(context.TODO(), foundDeployment); err != nil {
			reqLogger.Info("Reconcile deployment error", "Namespace", foundDeployment.Namespace, "Name", foundDeployment.Name)
			raven.CaptureErrorAndWait(err, nil)
			return reconcile.Result{}, err
		}
	}

	// reconcile service
	newService := generateService(instance)

	if err := controllerutil.SetControllerReference(instance, &newService, r.scheme); err != nil {
		raven.CaptureErrorAndWait(err, nil)
		return reconcile.Result{}, err
	}

	foundService := &corev1.Service{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: newService.Name, Namespace: newService.Namespace}, foundService)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new Service", "Namespace", newService.Namespace, "Name", newService.Name)
		err = r.client.Create(context.TODO(), &newService)
		if err != nil {
			raven.CaptureErrorAndWait(err, nil)
			return reconcile.Result{}, err
		}
	} else if err != nil {
		raven.CaptureErrorAndWait(err, nil)
		return reconcile.Result{}, err
	}

	reconcileService := false

	if !reflect.DeepEqual(foundService.Spec, newService.Spec) {
		foundService.Spec.Ports = newService.Spec.Ports
		reconcileService = true
	}

	if !reflect.DeepEqual(foundService.Annotations, newService.Annotations) {
		foundService.Annotations = newService.Annotations
		reconcileService = true
	}

	if !reflect.DeepEqual(foundService.Labels, newService.Labels) {
		foundService.Labels = newService.Labels
		reconcileService = true
	}

	if reconcileService {
		if err = r.client.Update(context.TODO(), foundService); err != nil {
			reqLogger.Info("Reconcile service error", "Namespace", foundService.Namespace, "Name", foundService.Name)
			raven.CaptureErrorAndWait(err, nil)
			return reconcile.Result{}, err
		}
	}

	// if configmap name set manually in CR, typically not used ever, but hell knows...
	configmapName := instance.Name + "-controller"
	if instance.Spec.NginxController.ConfigMap != "" {
		configmapName = instance.Spec.NginxController.ConfigMap
	}

	// reconcile configmap
	newConfigmap := generateConfigmap(instance, configmapName)

	if err := controllerutil.SetControllerReference(instance, &newConfigmap, r.scheme); err != nil {
		raven.CaptureErrorAndWait(err, nil)
		return reconcile.Result{}, err
	}

	foundConfigmap := &corev1.ConfigMap{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: newConfigmap.Name, Namespace: newConfigmap.Namespace}, foundConfigmap)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new Configmap", "Namespace", newConfigmap.Namespace, "Name", newConfigmap.Name)
		err = r.client.Create(context.TODO(), &newConfigmap)
		if err != nil {
			raven.CaptureErrorAndWait(err, nil)
			return reconcile.Result{}, err
		}
	} else if err != nil {
		raven.CaptureErrorAndWait(err, nil)
		return reconcile.Result{}, err
	}

	reconcileConfigmap := false

	if !reflect.DeepEqual(foundConfigmap.Data, newConfigmap.Data) {
		foundConfigmap.Data = newConfigmap.Data
		reconcileConfigmap = true
	}

	if !reflect.DeepEqual(foundConfigmap.Annotations, newConfigmap.Annotations) {
		foundConfigmap.Annotations = newConfigmap.Annotations
		reconcileConfigmap = true
	}

	if !reflect.DeepEqual(foundConfigmap.Labels, newConfigmap.Labels) {
		foundConfigmap.Labels = newConfigmap.Labels
		reconcileConfigmap = true
	}

	if reconcileConfigmap {
		if err = r.client.Update(context.TODO(), foundConfigmap); err != nil {
			reqLogger.Info("Reconcile configmap error", "Namespace", foundConfigmap.Namespace, "Name", foundConfigmap.Name)
			raven.CaptureErrorAndWait(err, nil)
			return reconcile.Result{}, err
		}
	}

	reqLogger.Info("Reconcile fully complete", "Namespace", foundDeployment.Namespace, "Name", foundDeployment.Name)
	return reconcile.Result{}, nil
}
