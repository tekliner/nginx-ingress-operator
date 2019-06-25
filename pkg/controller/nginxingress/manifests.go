package nginxingress

import (
	appv1alpha1 "github.com/tekliner/nginx-ingress-operator/pkg/apis/app/v1alpha1"
	"k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

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
		Spec: cr.Spec.NginxServiceSpec,
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
	}
	if cr.Spec.NginxController.IngressClass != "" {
		args = append(args, "--ingress-class="+cr.Spec.NginxController.IngressClass)
	}
	if cr.Spec.NginxController.ConfigMap != "" {
		args = append(args, "--configmap="+cr.Spec.NginxController.ConfigMap)
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
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 80,
								},
							},
							Env: append(returnDefaultENV(), cr.Spec.NginxController.Env...),
						},
					},
				},
			},
		},
	}

	return deployment
}
