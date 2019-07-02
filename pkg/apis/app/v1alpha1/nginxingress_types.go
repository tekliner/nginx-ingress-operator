package v1alpha1

import (
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ImageSpec struct {
	Tag        string `json:"tag"`
	Repository string `json:"repository"`
	PullPolicy string `json:"pullPolicy,omitempty"`
}

type DefaultBackendSpec struct {
	Name  string    `json:"name"`
	Image ImageSpec `json:"image"`
}

type MetricsServiceSpecs struct {
	Port        int32             `json:"port,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
}

type StatsSpec struct {
	Port int32 `json:"port,omitempty"`
}

type IngressServiceSpec struct {
	Type                  v1.ServiceType                      `json:"serviceType,omitempty"`
	ExternalTrafficPolicy v1.ServiceExternalTrafficPolicyType `json:"externalTrafficPolicy,omitempty"`
	Annotations           map[string]string                   `json:"annotations,omitempty"`
}

type NginxControllerSpec struct {
	Image      ImageSpec   `json:"image,omitempty"` // default quay.io/kubernetes-ingress-controller/nginx-ingress-controller:0.24.1
	Env        []v1.EnvVar `json:"env"`
	ElectionID string      `json:"electionID"`
	CustomArgs []string    `json:"customArgs,omitempty"`
	RunAsUser  int64       `json:"runAsUser,omitempty"` // default 33

	// nginx configuration
	Config         map[string]string `json:"config,omitempty"`
	Labels         map[string]string `json:"labels,omitempty"`
	Headers        string            `json:"headers,omitempty"`
	HostNetwork    bool              `json:"hostNetwork,omitempty"`
	ConfigMap      string            `json:"configmap,omitempty"`
	ConfigMapNginx string            `json:"configmapNginx,omitempty"`
	ConfigMapTCP   string            `json:"configmapTCP,omitempty"`
	ConfigMapUDP   string            `json:"configmapUDP,omitempty"`
	WatchNamespace string            `json:"watchNamespace,omitempty"`

	PriorityClassName     string       `json:"priorityClassName.omitempty"`
	DefaultBackendService string       `json:"defaultBackendService,omitempty"`
	DNSPolicy             v1.DNSPolicy `json:"dnsPolicy,omitempty"`
	IngressClass          string       `json:"ingressClass,omitempty"`
}

// NginxIngressSpec defines the desired state of NginxIngress
// +k8s:openapi-gen=true
type NginxIngressSpec struct {
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html

	Replicas int32 `json:"replicas"`

	Metrics                 *MetricsServiceSpecs `json:"metrics,omitempty"`
	Stats                   *StatsSpec           `json:"stats,omitempty"`
	NginxController         NginxControllerSpec  `json:"nginxController"`
	DefaultBackend          DefaultBackendSpec   `json:"defaultBackend,omitempty"`
	NginxServiceSpec        v1.ServiceSpec       `json:"service"`
	NginxServiceAnnotations map[string]string    `json:"serviceAnnotations,omitempty"`

	ServiceAccount string `json:"serviceAccount"`
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
