package v1alpha1

import (
	"k8s.io/api/core/v1"
	"k8s.io/api/policy/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

type ImageSpec struct {
	Tag        string         `json:"tag"`
	Repository string         `json:"repository"`
	PullPolicy *v1.PullPolicy `json:"pullPolicy,omitempty"`
}

type DefaultBackendSpec struct {
	Name               string                      `json:"name,omitempty"`               // default cr.name+"-default-backend"
	Namespace          string                      `json:"namespace,omitempty"`          // default current namespace
	Image              ImageSpec                   `json:"image,omitempty"`              // default k8s.gcr.io/defaultbackend-amd64:1.5
	CustomArgs         []string                    `json:"customArgs,omitempty"`         // pass arguments to container
	RunAsUser          *int64                      `json:"runAsUser,omitempty"`          // default 65534
	Env                *[]v1.EnvVar                `json:"env,omitempty"`                // pass ENVs to container
	Port               *int32                      `json:"port,omitempty"`               // default 8080
	Affinity           *v1.Affinity                `json:"affinity,omitempty"`           // default empty
	Replicas           *int32                      `json:"replicas,omitempty"`           // default 1
	Annotations        *map[string]string          `json:"annotations,omitempty"`        // default empty
	PodRequests        *v1.ResourceList            `json:"podRequests,omitempty"`        // default empty
	PodLimits          *v1.ResourceList            `json:"podLimits,omitempty"`          // default empty
	ServiceAnnotations *map[string]string          `json:"serviceAnnotations,omitempty"` // default empty
	Pdb                v1beta1.PodDisruptionBudget `json:"pdb"`
}

func (in *DefaultBackendSpec) GetDisruptionBudget(name string, namespace string, replicas int32, selector map[string]string) v1beta1.PodDisruptionBudget {
	podDisruptionBudget := v1beta1.PodDisruptionBudget{}
	labelSelector := metav1.LabelSelector{MatchLabels: selector}
	if replicas >= 2 {
		specPDB := v1beta1.PodDisruptionBudgetSpec{Selector: &labelSelector}

		if in.Pdb.Spec.MinAvailable != nil && in.Pdb.Spec.MaxUnavailable != nil {
			specPDB.MinAvailable = in.Pdb.Spec.MinAvailable
		} else if in.Pdb.Spec.MinAvailable != nil {
			specPDB.MinAvailable = in.Pdb.Spec.MinAvailable
		} else if in.Pdb.Spec.MaxUnavailable != nil {
			specPDB.MaxUnavailable = in.Pdb.Spec.MaxUnavailable
		} else {
			specPDB.MinAvailable = func() *intstr.IntOrString { v := intstr.FromInt(1); return &v }()
		}

		podDisruptionBudget = v1beta1.PodDisruptionBudget{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name + "-backend",
				Namespace: namespace,
			},
			Spec: specPDB,
		}
	} else {
		maxUnavailable := intstr.FromInt(1)
		specPDB := v1beta1.PodDisruptionBudgetSpec{
			Selector:       &labelSelector,
			MaxUnavailable: &maxUnavailable,
		}
		podDisruptionBudget = v1beta1.PodDisruptionBudget{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name + "-backend",
				Namespace: namespace,
			},
			Spec: specPDB,
		}
	}
	return podDisruptionBudget
}

func (in *NginxIngress) GetBackendLabels() map[string]string {
	return map[string]string{
		"app.improvado.io/application": "nginx-ingress-controller",
		"app.improvado.io/instance":    in.Name,
		"app.improvado.io/component":   "backend",
	}
}

func (in *NginxIngress) GetControllerLabels() map[string]string {
	return map[string]string{
		"app.improvado.io/application": "nginx-ingress-controller",
		"app.improvado.io/instance":    in.Name,
		"app.improvado.io/component":   "controller",
	}
}

func (in *NginxControllerSpec) GetDisruptionBudget(name string, namespace string, replicas int32, selector map[string]string) v1beta1.PodDisruptionBudget {
	podDisruptionBudget := v1beta1.PodDisruptionBudget{}
	labelSelector := metav1.LabelSelector{MatchLabels: selector}
	if replicas >= 2 {
		specPDB := v1beta1.PodDisruptionBudgetSpec{Selector: &labelSelector}

		if in.Pdb.Spec.MinAvailable != nil && in.Pdb.Spec.MaxUnavailable != nil {
			specPDB.MinAvailable = in.Pdb.Spec.MinAvailable
		} else if in.Pdb.Spec.MinAvailable != nil {
			specPDB.MinAvailable = in.Pdb.Spec.MinAvailable
		} else if in.Pdb.Spec.MaxUnavailable != nil {
			specPDB.MaxUnavailable = in.Pdb.Spec.MaxUnavailable
		} else {
			specPDB.MinAvailable = func() *intstr.IntOrString { v := intstr.FromInt(1); return &v }()
		}

		podDisruptionBudget = v1beta1.PodDisruptionBudget{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name + "-controller",
				Namespace: namespace,
			},
			Spec: specPDB,
		}
	} else {
		maxUnavailable := intstr.FromInt(1)
		specPDB := v1beta1.PodDisruptionBudgetSpec{
			Selector:       &labelSelector,
			MaxUnavailable: &maxUnavailable,
		}
		podDisruptionBudget = v1beta1.PodDisruptionBudget{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name + "-controller",
				Namespace: namespace,
			},
			Spec: specPDB,
		}
	}
	return podDisruptionBudget
}

type MetricsServiceSpecs struct { // define to turn on
	Port        int32             `json:"port"` // required
	Annotations map[string]string `json:"annotations,omitempty"`
}

type StatsSpec struct { // define to turn on
	Port int32 `json:"port"` // required
}

type IngressServiceSpec struct {
	Type                  v1.ServiceType                      `json:"serviceType,omitempty"`
	ExternalTrafficPolicy v1.ServiceExternalTrafficPolicyType `json:"externalTrafficPolicy,omitempty"`
}

type NginxControllerSpec struct {
	// default quay.io/kubernetes-ingress-controller/nginx-ingress-controller:0.24.1
	Image       ImageSpec          `json:"image"`
	Env         *[]v1.EnvVar       `json:"env,omitempty"`
	ElectionID  string             `json:"electionID"`
	CustomArgs  []string           `json:"customArgs,omitempty"`
	RunAsUser   *int64             `json:"runAsUser,omitempty"` // default 33
	Affinity    *v1.Affinity       `json:"affinity,omitempty"`
	PodRequests *v1.ResourceList   `json:"pod_requests,omitempty"`
	PodLimits   *v1.ResourceList   `json:"pod_limits,omitempty"`
	Annotations *map[string]string `json:"annotations,omitempty"`

	// nginx configuration
	// set CM name manually, if user want to use own CM, f.e. one config for many instances
	// if not defined CM will be autogenerated with data below
	ConfigMap string `json:"configmap,omitempty"`
	// data, ignored if CM name set manually
	Config map[string]string `json:"config,omitempty"`

	Labels             map[string]string `json:"labels,omitempty"`
	Headers            string            `json:"headers,omitempty"`
	HostNetwork        bool              `json:"hostNetwork,omitempty"`
	ConfigMapTCP       string            `json:"configmapTCP,omitempty"`
	ConfigMapUDP       string            `json:"configmapUDP,omitempty"`
	WatchNamespace     string            `json:"watchNamespace,omitempty"`
	PublishService     bool              `json:"publishService,omitempty"`
	PublishServicePath string            `json:"publishServicePath,omitempty"` // override generated value

	PriorityClassName     string                      `json:"priorityClassName,omitempty"`
	DefaultBackendService string                      `json:"defaultBackendService,omitempty"` // format: namespace/svcname
	DNSPolicy             v1.DNSPolicy                `json:"dnsPolicy,omitempty"`
	IngressClass          string                      `json:"ingressClass,omitempty"`
	Pdb                   v1beta1.PodDisruptionBudget `json:"pdb,omitempty"`
}

// NginxIngressSpec defines the desired state of NginxIngress
// +k8s:openapi-gen=true
type NginxIngressSpec struct {
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html

	Replicas int32 `json:"replicas"`

	Metrics                 *MetricsServiceSpecs `json:"metrics,omitempty"`
	Stats                   *StatsSpec           `json:"stats,omitempty"`
	NginxController         NginxControllerSpec  `json:"nginxController"`
	DefaultBackend          *DefaultBackendSpec  `json:"defaultBackend,omitempty"`
	NginxServiceSpec        v1.ServiceSpec       `json:"service"`
	NginxServiceAnnotations map[string]string    `json:"serviceAnnotations,omitempty"`

	ServiceAccount string `json:"serviceAccount,omitempty"`
}

// NginxIngressStatus defines the observed state of NginxIngress
// +k8s:openapi-gen=true
type NginxIngressStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// NginxIngress is the Schema for the nginxingresses API
// +k8s:openapi-gen=true
type NginxIngress struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NginxIngressSpec   `json:"spec,omitempty"`
	Status NginxIngressStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// NginxIngressList contains a list of NginxIngress
type NginxIngressList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NginxIngress `json:"items"`
}

func init() {
	SchemeBuilder.Register(&NginxIngress{}, &NginxIngressList{})
}
