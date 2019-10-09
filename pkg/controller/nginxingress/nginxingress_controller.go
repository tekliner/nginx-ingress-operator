package nginxingress

import (
	"context"
	"os"

	raven "github.com/getsentry/raven-go"
	appv1alpha1 "github.com/tekliner/nginx-ingress-operator/pkg/apis/app/v1alpha1"
	"k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	v1beta1policy "k8s.io/api/policy/v1beta1"
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
		raven.CaptureErrorAndWait(err, nil)
		return err
	}

	// Watch for changes to primary resource NginxIngress
	err = c.Watch(&source.Kind{Type: &appv1alpha1.NginxIngress{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		raven.CaptureErrorAndWait(err, nil)
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner NginxIngress
	err = c.Watch(&source.Kind{Type: &v1.Deployment{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &appv1alpha1.NginxIngress{},
	})
	if err != nil {
		raven.CaptureErrorAndWait(err, nil)
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

	// defaultBackend deployment and service
	// checking at first, if defaultBackend not defined controller will use NginxController.DefaultBackendService string
	// https://kubernetes.github.io/ingress-nginx/user-guide/cli-arguments/
	// Service used to serve HTTP requests not matching any known server name (catch-all).
	if instance.Spec.NginxController.DefaultBackendService == "" {
		newBackendDeployment := generateDefaultBackendDeployment(instance)

		if err := controllerutil.SetControllerReference(instance, &newBackendDeployment, r.scheme); err != nil {
			raven.CaptureErrorAndWait(err, nil)
			return reconcile.Result{}, err
		}

		// default backend deployment
		foundBackendDeployment := v1.Deployment{}

		err = r.client.Get(context.TODO(), types.NamespacedName{Name: newBackendDeployment.Name, Namespace: newBackendDeployment.Namespace}, &foundBackendDeployment)
		if err != nil && errors.IsNotFound(err) {
			reqLogger.Info("Creating a new Deployment", "Namespace", newBackendDeployment.Namespace, "Name", newBackendDeployment.Name)
			err = r.client.Create(context.TODO(), &newBackendDeployment)
			if err != nil {
				raven.CaptureErrorAndWait(err, nil)
				return reconcile.Result{}, err
			}

		} else if err != nil {
			raven.CaptureErrorAndWait(err, nil)
			return reconcile.Result{}, err
		} else {
			if reconcileRequired, reconDeployment := reconcileDeployment(foundBackendDeployment, newBackendDeployment); reconcileRequired {
				reqLogger.Info("Updating Deployment", "Namespace", reconDeployment.Namespace, "Name", reconDeployment.Name)
				if err = r.client.Update(context.TODO(), &reconDeployment); err != nil {
					reqLogger.Info("Reconcile deployment error", "Namespace", foundBackendDeployment.Namespace, "Name", foundBackendDeployment.Name)
					raven.CaptureErrorAndWait(err, nil)
					return reconcile.Result{}, err
				}
			}
		}

		// reconcile default backend service
		newDefaultService := generateDefaultBackendService(instance)

		if err := controllerutil.SetControllerReference(instance, &newDefaultService, r.scheme); err != nil {
			raven.CaptureErrorAndWait(err, nil)
			return reconcile.Result{}, err
		}

		foundDefaultService := corev1.Service{}
		err = r.client.Get(context.TODO(), types.NamespacedName{Name: newDefaultService.Name, Namespace: newDefaultService.Namespace}, &foundDefaultService)
		if err != nil && errors.IsNotFound(err) {
			reqLogger.Info("Creating a new Service", "Namespace", newDefaultService.Namespace, "Name", newDefaultService.Name)
			err = r.client.Create(context.TODO(), &newDefaultService)
			if err != nil {
				raven.CaptureErrorAndWait(err, nil)
				return reconcile.Result{}, err
			}
		} else if err != nil {
			raven.CaptureErrorAndWait(err, nil)
			return reconcile.Result{}, err
		} else {
			if reconcileRequired, reconService := reconcileService(foundDefaultService, newDefaultService); reconcileRequired {
				reqLogger.Info("Updating Service", "Namespace", reconService.Namespace, "Name", reconService.Name)
				if err = r.client.Update(context.TODO(), &reconService); err != nil {
					reqLogger.Info("Reconcile service error", "Namespace", foundDefaultService.Namespace, "Name", foundDefaultService.Name)
					raven.CaptureErrorAndWait(err, nil)
					return reconcile.Result{}, err
				}
			}
		}
	}

	// reconcile service
	newService := generateService(instance)

	if err := controllerutil.SetControllerReference(instance, &newService, r.scheme); err != nil {
		raven.CaptureErrorAndWait(err, nil)
		return reconcile.Result{}, err
	}

	foundService := corev1.Service{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: newService.Name, Namespace: newService.Namespace}, &foundService)
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
	} else {
		if reconcileRequired, reconService := reconcileService(foundService, newService); reconcileRequired {
			reqLogger.Info("Updating Service", "Namespace", reconService.Namespace, "Name", reconService.Name)
			if err = r.client.Update(context.TODO(), &reconService); err != nil {
				reqLogger.Info("Reconcile service error", "Namespace", foundService.Namespace, "Name", foundService.Name)
				raven.CaptureErrorAndWait(err, nil)
				return reconcile.Result{}, err
			}
		}
	}

	// reconcile metrics service
	if instance.Spec.Metrics != nil {
		newService = generateServiceMetrics(instance)

		if err := controllerutil.SetControllerReference(instance, &newService, r.scheme); err != nil {
			raven.CaptureErrorAndWait(err, nil)
			return reconcile.Result{}, err
		}

		foundService = corev1.Service{}
		err = r.client.Get(context.TODO(), types.NamespacedName{Name: newService.Name, Namespace: newService.Namespace}, &foundService)
		if err != nil && errors.IsNotFound(err) {
			reqLogger.Info("Creating a new metrics Service", "Namespace", newService.Namespace, "Name", newService.Name)
			err = r.client.Create(context.TODO(), &newService)
			if err != nil {
				raven.CaptureErrorAndWait(err, nil)
				return reconcile.Result{}, err
			}
		} else if err != nil {
			raven.CaptureErrorAndWait(err, nil)
			return reconcile.Result{}, err
		} else {
			if reconcileRequired, reconService := reconcileService(foundService, newService); reconcileRequired {
				reqLogger.Info("Updating Metrics Service", "Namespace", reconService.Namespace, "Name", reconService.Name)
				if err = r.client.Update(context.TODO(), &reconService); err != nil {
					reqLogger.Info("Reconcile metrics Service error", "Namespace", foundService.Namespace, "Name", foundService.Name)
					raven.CaptureErrorAndWait(err, nil)
					return reconcile.Result{}, err
				}
			}
		}
	}

	// reconcile stats service
	if instance.Spec.Stats != nil {
		newService = generateServiceStats(instance)

		if err := controllerutil.SetControllerReference(instance, &newService, r.scheme); err != nil {
			raven.CaptureErrorAndWait(err, nil)
			return reconcile.Result{}, err
		}

		foundService = corev1.Service{}
		err = r.client.Get(context.TODO(), types.NamespacedName{Name: newService.Name, Namespace: newService.Namespace}, &foundService)
		if err != nil && errors.IsNotFound(err) {
			reqLogger.Info("Creating a new stats Service", "Namespace", newService.Namespace, "Name", newService.Name)
			err = r.client.Create(context.TODO(), &newService)
			if err != nil {
				raven.CaptureErrorAndWait(err, nil)
				return reconcile.Result{}, err
			}
		} else if err != nil {
			raven.CaptureErrorAndWait(err, nil)
			return reconcile.Result{}, err
		} else {
			if reconcileRequired, reconService := reconcileService(foundService, newService); reconcileRequired {
				reqLogger.Info("Updating Stats Service", "Namespace", reconService.Namespace, "Name", reconService.Name)
				if err = r.client.Update(context.TODO(), &reconService); err != nil {
					reqLogger.Info("Reconcile stats Service error", "Namespace", foundService.Namespace, "Name", foundService.Name)
					raven.CaptureErrorAndWait(err, nil)
					return reconcile.Result{}, err
				}
			}
		}
	}

	// if CM name not set manually then create CM based on nginxController.config
	if instance.Spec.NginxController.ConfigMap == "" {
		// reconcile configmap
		newConfigmap := generateConfigmap(instance)

		if err := controllerutil.SetControllerReference(instance, &newConfigmap, r.scheme); err != nil {
			raven.CaptureErrorAndWait(err, nil)
			return reconcile.Result{}, err
		}

		foundConfigmap := corev1.ConfigMap{}
		err = r.client.Get(context.TODO(), types.NamespacedName{Name: newConfigmap.Name, Namespace: newConfigmap.Namespace}, &foundConfigmap)
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
		} else {
			if reconcileRequired, reconConfigmap := reconcileConfigmap(foundConfigmap, newConfigmap); reconcileRequired {
				reqLogger.Info("Updating Configmap", "Namespace", reconConfigmap.Namespace, "Name", reconConfigmap.Name)
				if err = r.client.Update(context.TODO(), &reconConfigmap); err != nil {
					reqLogger.Info("Reconcile configmap error", "Namespace", foundConfigmap.Namespace, "Name", foundConfigmap.Name)
					raven.CaptureErrorAndWait(err, nil)
					return reconcile.Result{}, err
				}
			}
		}
	}

	// reconcile deployment
	newDeployment := generateDeployment(instance)

	if err := controllerutil.SetControllerReference(instance, &newDeployment, r.scheme); err != nil {
		raven.CaptureErrorAndWait(err, nil)
		return reconcile.Result{}, err
	}

	// controller deployment
	foundDeployment := v1.Deployment{}

	err = r.client.Get(context.TODO(), types.NamespacedName{Name: newDeployment.Name, Namespace: newDeployment.Namespace}, &foundDeployment)
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
	} else {
		if reconcileRequired, reconDeployment := reconcileDeployment(foundDeployment, newDeployment); reconcileRequired {
			reqLogger.Info("Updating Deployment", "Namespace", reconDeployment.Namespace, "Name", reconDeployment.Name)
			if err = r.client.Update(context.TODO(), &reconDeployment); err != nil {
				reqLogger.Info("Reconcile deployment error", "Namespace", foundDeployment.Namespace, "Name", foundDeployment.Name)
				raven.CaptureErrorAndWait(err, nil)
				return reconcile.Result{}, err
			}
		}
	}

	// reconcile backend podDisruptionBudget

	newBackendPDB := generatePodDisruptionBudget(instance, "-default-backend")

	if err := controllerutil.SetControllerReference(instance, &newBackendPDB, r.scheme); err != nil {
		raven.CaptureErrorAndWait(err, nil)
		return reconcile.Result{}, err
	}

	foundBackendPDB := v1beta1policy.PodDisruptionBudget{}

	err = r.client.Get(context.TODO(), types.NamespacedName{Name: newBackendPDB.Name, Namespace: newBackendPDB.Namespace}, &foundBackendPDB)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new PodDisruptionBudget", "Namespace", newBackendPDB.Namespace, "Name", newBackendPDB.Name)
		err = r.client.Create(context.TODO(), &newBackendPDB)
		if err != nil {
			raven.CaptureErrorAndWait(err, nil)
			return reconcile.Result{}, err
		}
	} else if err != nil {
		raven.CaptureErrorAndWait(err, nil)
		return reconcile.Result{}, err
	} else {
		if reconcileRequired, reconPDB := reconcilePdb(foundBackendPDB, newBackendPDB); reconcileRequired {
			reqLogger.Info("Updating PodDisruptionBudget", "Namespace", reconPDB.Namespace, "Name", reconPDB.Name)
			if err = r.client.Update(context.TODO(), &reconPDB); err != nil {
				reqLogger.Info("Reconcile PodDisruptionBudget error", "Namespace", foundBackendPDB.Namespace, "Name", foundBackendPDB.Name)
				raven.CaptureErrorAndWait(err, nil)
				return reconcile.Result{}, err
			}
		}
	}

	// reconcile controller podDisruptionBudget

	//newControllerBdp := generatePodDisruptionBudget(instance, "-controller")
	//
	//if err := controllerutil.SetControllerReference(instance, &newControllerBdp, r.scheme); err != nil {
	//	raven.CaptureErrorAndWait(err, nil)
	//	return reconcile.Result{}, err
	//}
	//
	//foundControllerPDB := v1beta1policy.PodDisruptionBudget{}
	//
	//err = r.client.Get(context.TODO(), types.NamespacedName{Name: newControllerBdp.Name, Namespace: newControllerBdp.Namespace}, &foundControllerPDB)
	//if err != nil && errors.IsNotFound(err) {
	//	reqLogger.Info("Creating a new ControllerPDB", "Namespace", newControllerBdp.Namespace, "Name", newControllerBdp.Name)
	//	err = r.client.Create(context.TODO(), &newControllerBdp)
	//	if err != nil {
	//		raven.CaptureErrorAndWait(err, nil)
	//		return reconcile.Result{}, err
	//	}
	//} else if err != nil {
	//	raven.CaptureErrorAndWait(err, nil)
	//	return reconcile.Result{}, err
	//} else {
	//	if reconcileRequired, reconPDB := reconcilePdb(foundControllerPDB, newControllerBdp); reconcileRequired {
	//		reqLogger.Info("Updating ControllerPDB", "Namespace", reconPDB.Namespace, "Name", reconPDB.Name)
	//		if err = r.client.Update(context.TODO(), &reconPDB); err != nil {
	//			reqLogger.Info("Reconcile ControllerPDB error", "Namespace", foundControllerPDB.Namespace, "Name", foundControllerPDB.Name)
	//			raven.CaptureErrorAndWait(err, nil)
	//			return reconcile.Result{}, err
	//		}
	//	}
	//}

	reqLogger.Info("Reconcile fully complete", "Namespace", foundDeployment.Namespace, "Name", foundDeployment.Name)
	return reconcile.Result{}, nil
}
