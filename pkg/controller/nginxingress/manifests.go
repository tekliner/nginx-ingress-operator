package nginxingress

import (
	appv1alpha1 "github.com/tekliner/nginx-ingress-operator/pkg/apis/app/v1alpha1"
	"k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func generateConfigmap(cr *appv1alpha1.NginxIngress, configmapName string) corev1.ConfigMap {
	data := cr.Spec.NginxController.Config

	return corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      configmapName,
			Namespace: cr.Namespace,
		},
		Data: data,
	}
}

func generateServiceMetrics(cr *appv1alpha1.NginxIngress) corev1.Service {
	labels := map[string]string{
		"app.improvado.io/component": "service",
	}

	annotations := map[string]string{}

	service := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:        cr.Name,
			Namespace:   cr.Namespace,
			Labels:      labels,
			Annotations: mergeMaps(cr.Spec.Metrics.Annotations, annotations),
		},
	}

	service.Spec.Type = corev1.ServiceTypeClusterIP

	service.Spec.Ports = []corev1.ServicePort{
		{
			Name:       "exporter",
			Port:       cr.Spec.Metrics.Port,
			Protocol:   corev1.ProtocolTCP,
			TargetPort: intstr.IntOrString{StrVal: "metrics"},
		},
	}

	service.Spec.Selector = baseLabels(cr)

	return service
}

func generateServiceStats(cr *appv1alpha1.NginxIngress) corev1.Service {
	labels := map[string]string{
		"app.improvado.io/component": "service",
	}

	annotations := map[string]string{}

	service := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:        cr.Name,
			Namespace:   cr.Namespace,
			Labels:      labels,
			Annotations: mergeMaps(cr.Spec.Metrics.Annotations, annotations),
		},
	}

	service.Spec.Type = corev1.ServiceTypeClusterIP

	service.Spec.Ports = []corev1.ServicePort{
		{
			Name:       "stats",
			Port:       cr.Spec.Metrics.Port,
			Protocol:   corev1.ProtocolTCP,
			TargetPort: intstr.IntOrString{StrVal: "stats"},
		},
	}

	service.Spec.Selector = baseLabels(cr)

	return service
}

func generateService(cr *appv1alpha1.NginxIngress) corev1.Service {
	labels := map[string]string{
		"app.improvado.io/component": "service",
	}

	service := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:        cr.Name,
			Namespace:   cr.Namespace,
			Labels:      labels,
			Annotations: cr.Spec.NginxServiceAnnotations,
		},
	}

	service.Spec.Type = corev1.ServiceTypeLoadBalancer

	service.Spec.Ports = []corev1.ServicePort{
		{
			Name:       "http",
			Port:       80,
			TargetPort: intstr.IntOrString{StrVal: "http"},
			Protocol:   corev1.ProtocolTCP,
		},
		{
			Name:       "https",
			Port:       443,
			TargetPort: intstr.IntOrString{StrVal: "https"},
			Protocol:   corev1.ProtocolTCP,
		},
	}

	service.Spec.Selector = baseLabels(cr)

	return service
}

func generateDeployment(cr *appv1alpha1.NginxIngress) v1.Deployment {

	// compile arguments from CR
	args := []string{"/nginx-ingress-controller"}

	if cr.Spec.NginxController.DefaultBackendService != "" {
		args = append(args, "--default-backend-service="+cr.Spec.NginxController.DefaultBackendService)
	}

	if cr.Spec.NginxController.ElectionID != "" {
		args = append(args, "--election-id="+cr.Spec.NginxController.ElectionID)
	} else {
		args = append(args, "--election-id=ingress-leader-election-"+cr.Name)
	}

	if cr.Spec.NginxController.IngressClass != "" {
		args = append(args, "--ingress-class="+cr.Spec.NginxController.IngressClass)
	} else {
		args = append(args, "--ingress-class=ingress-class-"+cr.Name)
	}

	if cr.Spec.NginxController.ConfigMap != "" {
		args = append(args, "--configmap="+cr.Spec.NginxController.ConfigMap)
	} else {
		args = append(args, "--configmap="+cr.Name)
	}

	if cr.Spec.NginxController.ConfigMapNginx != "" {
		args = append(args, "--nginx-configmap="+cr.Spec.NginxController.ConfigMapNginx)
	}

	if cr.Spec.NginxController.ConfigMapTCP != "" {
		args = append(args, "--tcp-services-configmap="+cr.Spec.NginxController.ConfigMapTCP)
	}

	if cr.Spec.NginxController.ConfigMapUDP != "" {
		args = append(args, "--udp-services-configmap="+cr.Spec.NginxController.ConfigMapUDP)
	}

	if cr.Spec.NginxController.WatchNamespace != "" {
		args = append(args, "--watch-namespace="+cr.Spec.NginxController.WatchNamespace)
	}

	// add custom arguments from CR
	args = append(args, cr.Spec.NginxController.CustomArgs...)

	ports := []corev1.ContainerPort{
		{
			Name:          "http",
			ContainerPort: 80,
			Protocol:      corev1.ProtocolTCP,
		},
		{
			Name:          "https",
			ContainerPort: 443,
			Protocol:      corev1.ProtocolTCP,
		},
	}

	if cr.Spec.Metrics != nil {
		metricsPort := corev1.ContainerPort{
			Name:          "metrics",
			ContainerPort: cr.Spec.Metrics.Port,
			Protocol:      corev1.ProtocolTCP,
		}
		ports = append(ports, metricsPort)
	}

	if cr.Spec.Stats != nil {
		statsPort := corev1.ContainerPort{
			Name:          "stats",
			ContainerPort: cr.Spec.Stats.Port,
			Protocol:      corev1.ProtocolTCP,
		}
		ports = append(ports, statsPort)
	}

	deployment := v1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name,
			Namespace: cr.ObjectMeta.Namespace,
			Labels:    baseLabels(cr),
		},
		Spec: v1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: baseLabels(cr),
			},
			Replicas: &cr.Spec.Replicas,
			Strategy: v1.DeploymentStrategy{Type: v1.RollingUpdateDeploymentStrategyType, RollingUpdate: nil},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: mergeMaps(baseLabels(cr),
						map[string]string{"app.improvado.io/component": "application"},
					),
					Annotations: setAnnotations(cr, cr.Annotations),
				},

				Spec: corev1.PodSpec{
					DNSPolicy:          cr.Spec.NginxController.DNSPolicy,
					ServiceAccountName: cr.Spec.ServiceAccount,
					PriorityClassName:  cr.Spec.NginxController.PriorityClassName,
					SecurityContext: &corev1.PodSecurityContext{
						RunAsUser: &cr.Spec.NginxController.RunAsUser,
					},
					Containers: []corev1.Container{
						{
							Name:  "nginx-ingress",
							Image: cr.Spec.NginxController.Image.Repository + ":" + cr.Spec.NginxController.Image.Tag,
							Args:  args,
							Ports: ports,
							Env:   append(returnDefaultENV(), cr.Spec.NginxController.Env...),
						},
					},
				},
			},
		},
	}

	return deployment
}
