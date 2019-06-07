package nginxingress

import (
	"context"
	"os"

	raven "github.com/getsentry/raven-go"
	appv1alpha1 "github.com/tekliner/nginx-ingress-operator/pkg/apis/app/v1alpha1"
	"k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
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
		return reconcile.Result{}, err
	}

	// Define a new Deployment object
	newDeployment := newDeployment(instance)

	// Set NginxIngress instance as the owner and controller
	if err := controllerutil.SetControllerReference(instance, newDeployment, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	foundDeployment := &v1.Deployment{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: foundDeployment.Name, Namespace: foundDeployment.Namespace}, foundDeployment)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new deployment", "Namespace", newDeployment.Namespace, "Name", newDeployment.Name)
		err = r.client.Create(context.TODO(), newDeployment)
		if err != nil {
			return reconcile.Result{}, err
		}
		return reconcile.Result{}, nil
	} else if err != nil {
		return reconcile.Result{}, err
	}

	// Deployment already exists - don't requeue
	reqLogger.Info("Skip reconcile: Deployment already exists", "Namespace", foundDeployment.Namespace, "Name", foundDeployment.Name)
	return reconcile.Result{}, nil
}

func newDeployment(cr *appv1alpha1.NginxIngress) *v1.Deployment {

	args := []string{}
	// add arguments for default command

	return &v1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name,
			Namespace: cr.ObjectMeta.Namespace,
			Labels:    baseLabels(cr),
		},

		Spec: v1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: baseLabels(cr),
			},
			Strategy: v1.DeploymentStrategy{Type: v1.RollingUpdateDeploymentStrategyType, RollingUpdate: nil},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: mergeMaps(baseLabels(cr),
						map[string]string{"app.improvado.io/component": "controller"},
					),
					Annotations: setAnnotations(cr, cr.Annotations),
				},

				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "nginx",
							Image: cr.Spec.NginxImage,
							Args:  args,
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 8000,
								},
							},
							Env: cr.Spec.Env,
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "static",
									MountPath: "/usr/src/app/jobprocessor/static",
								},
								{
									Name:      "media",
									MountPath: "/usr/src/app/jobprocessor/media",
								},
							},
							Resources: cr.Spec.Application.Resources,
						},
						{
							Name:  "jobprocessor-nginx",
							Image: cr.Spec.NginxImage.Repository + ":" + cr.Spec.NginxImage.Tag,
							Ports: []corev1.ContainerPort{
								{
									Name:          "http",
									ContainerPort: 80,
								},
							},
							LivenessProbe: &corev1.Probe{
								FailureThreshold:    3,
								InitialDelaySeconds: 20,
								PeriodSeconds:       40,
								SuccessThreshold:    1,
								TimeoutSeconds:      1,
								Handler: corev1.Handler{
									HTTPGet: &corev1.HTTPGetAction{
										Path:   "/-/liveness",
										Port:   intstr.IntOrString{IntVal: 8000},
										Scheme: corev1.URISchemeHTTP,
									},
								},
							},
							ReadinessProbe: &corev1.Probe{
								FailureThreshold:    3,
								InitialDelaySeconds: 10,
								PeriodSeconds:       20,
								SuccessThreshold:    1,
								TimeoutSeconds:      1,
								Handler: corev1.Handler{
									HTTPGet: &corev1.HTTPGetAction{
										Path:   "/-/readiness",
										Port:   intstr.IntOrString{IntVal: 8000},
										Scheme: corev1.URISchemeHTTP,
									},
								},
							},
							Lifecycle: &corev1.Lifecycle{
								PreStop: &corev1.Handler{
									Exec: &corev1.ExecAction{
										Command: []string{"sleep", "15"},
									},
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "nginx-config",
									MountPath: "/etc/nginx/nginx.conf",
									SubPath:   "nginx.conf",
								},
								{
									Name:      "static",
									MountPath: "/usr/src/app/jobprocessor/static",
								},
								{
									Name:      "media",
									MountPath: "/usr/src/app/jobprocessor/media",
								},
							},
						},
						{
							Name:  "statsd",
							Image: cr.Spec.StatsdImage.Repository + ":" + cr.Spec.StatsdImage.Tag,
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 9102,
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "nginx-config",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: cr.Name + "-nginx-config",
									},
								},
							},
						},
						{
							Name: "static",
							VolumeSource: corev1.VolumeSource{
								EmptyDir: &corev1.EmptyDirVolumeSource{},
							},
						},
						{
							Name: "media",
							VolumeSource: corev1.VolumeSource{
								EmptyDir: &corev1.EmptyDirVolumeSource{},
							},
						},
					},
				},
			},
		},
	}
}
