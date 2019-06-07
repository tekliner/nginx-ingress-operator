package v1alpha1

import (
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type defaultBackend struct {
	Name  string `json:"name"`
	Image string `json:"image"`
}

// NginxIngressSpec defines the desired state of NginxIngress
// +k8s:openapi-gen=true
type NginxIngressSpec struct {
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html

	NginxImage string            `json:"nginxImage,omitempty"` // default quay.io/kubernetes-ingress-controller/nginx-ingress-controller:0.24.1
	Replicas   int               `json:"replicas"`
	RunAsUser  int               `json:"runAsUser,omitempty"` // default 33
	Ports      map[string]string `json:"ports"`
	Env        []v1.EnvVar       `json:"env"`
	ElectionID string            `json:"electionID"`

	Annotations       map[string]string `json:"annotations,omitempty"`
	Labels            map[string]string `json:"labels,omitempty"`
	IngressClass      string            `json:"ingressClass"`
	PriorityClassName string            `json:"priorityClassName.omitempty"`

	DNSPolicy string `json:"dnsPolicy,omitempty"`

	DefaultBackend        defaultBackend `json:"defaultBackend,omitempty"`
	DefaultBackendService string         `json:"defaultBackendService,omitempty"`

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
