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

	service := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:        cr.Name + "-metrics",
			Namespace:   cr.Namespace,
			Labels:      labels,
			Annotations: cr.Spec.Metrics.Annotations,
		},
	}

	service.Spec.Type = corev1.ServiceTypeClusterIP

	service.Spec.Ports = []corev1.ServicePort{
		{
			Name:       "exporter",
			Port:       cr.Spec.Metrics.Port,
			Protocol:   corev1.ProtocolTCP,
			TargetPort: intstr.FromString("metrics"),
		},
	}

	service.Spec.Selector = baseLabels(cr)

	return service
}

func generateServiceStats(cr *appv1alpha1.NginxIngress) corev1.Service {
	labels := map[string]string{
		"app.improvado.io/component": "service",
	}

	service := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:        cr.Name + "-stats",
			Namespace:   cr.Namespace,
			Labels:      labels,
			Annotations: cr.Spec.Metrics.Annotations,
		},
	}

	service.Spec.Type = corev1.ServiceTypeClusterIP

	service.Spec.Ports = []corev1.ServicePort{
		{
			Name:       "stats",
			Port:       cr.Spec.Stats.Port,
			Protocol:   corev1.ProtocolTCP,
			TargetPort: intstr.FromString("stats"),
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
			TargetPort: intstr.FromString("http"),
			Protocol:   corev1.ProtocolTCP,
		},
		{
			Name:       "https",
			Port:       443,
			TargetPort: intstr.FromString("https"),
			Protocol:   corev1.ProtocolTCP,
		},
	}

	service.Spec.Selector = baseLabels(cr)

	return service
}

func generateDeployment(cr *appv1alpha1.NginxIngress) v1.Deployment {

	runAsUser := int64(33)
	if cr.Spec.NginxController.RunAsUser != nil {
		runAsUser = *cr.Spec.NginxController.RunAsUser
	}

	env := []corev1.EnvVar{}
	if cr.Spec.NginxController.Env != nil {
		env = *cr.Spec.NginxController.Env
	}

	// check affinity rules
	affinity := &corev1.Affinity{}
	if cr.Spec.NginxController.Affinity != nil {
		affinity = cr.Spec.NginxController.Affinity
	}

	annotations := map[string]string{}
	if cr.Spec.NginxController.Annotations != nil {
		annotations = *cr.Spec.NginxController.Annotations
	}

	resourcesLimits := corev1.ResourceList{}
	if cr.Spec.NginxController.PodLimits != nil {
		resourcesLimits = *cr.Spec.NginxController.PodLimits
	}

	resourcesRequests := corev1.ResourceList{}
	if cr.Spec.NginxController.PodRequests != nil {
		resourcesRequests = *cr.Spec.NginxController.PodRequests
	}

	// compile arguments from CR
	args := []string{"/nginx-ingress-controller"}

	if cr.Spec.DefaultBackend != nil {
		args = append(args, "--default-backend-service="+cr.Name+"-default-backend")
	} else if cr.Spec.NginxController.DefaultBackendService != "" {
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

	if cr.Spec.NginxController.PublishService && cr.Spec.NginxController.PublishServicePath != "" {
		args = append(args, "--publish-service="+cr.Spec.NginxController.PublishServicePath)
	} else if cr.Spec.NginxController.PublishService {
		args = append(args, "--publish-service="+cr.Namespace+"/"+cr.Name+"-controller")
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
			ContainerPort: 10254,
			Protocol:      corev1.ProtocolTCP,
		}
		ports = append(ports, metricsPort)
	}

	if cr.Spec.Stats != nil {
		statsPort := corev1.ContainerPort{
			Name:          "stats",
			ContainerPort: 18080,
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
					Annotations: annotations,
				},

				Spec: corev1.PodSpec{
					Affinity:           affinity,
					DNSPolicy:          cr.Spec.NginxController.DNSPolicy,
					ServiceAccountName: cr.Spec.ServiceAccount,
					PriorityClassName:  cr.Spec.NginxController.PriorityClassName,
					SecurityContext: &corev1.PodSecurityContext{
						RunAsUser: &runAsUser,
					},
					Containers: []corev1.Container{
						{
							Name:  "nginx-ingress",
							Image: cr.Spec.NginxController.Image.Repository + ":" + cr.Spec.NginxController.Image.Tag,
							Args:  args,
							Ports: ports,
							Env:   env,
							Resources: corev1.ResourceRequirements{
								Limits:   resourcesLimits,
								Requests: resourcesRequests,
							},
						},
					},
				},
			},
		},
	}

	return deployment
}

func generateDefaultBackendDeployment(cr *appv1alpha1.NginxIngress) v1.Deployment {

	name := cr.Name + "-default-backend"
	if cr.Spec.DefaultBackend.Name != "" {
		name = cr.Spec.DefaultBackend.Name
	}

	runAsUser := int64(0)
	if cr.Spec.DefaultBackend.RunAsUser != nil {
		runAsUser = *cr.Spec.DefaultBackend.RunAsUser
	}

	env := []corev1.EnvVar{}
	if cr.Spec.DefaultBackend.Env != nil {
		env = *cr.Spec.DefaultBackend.Env
	}

	// check affinity rules
	affinity := &corev1.Affinity{}
	if cr.Spec.DefaultBackend.Affinity != nil {
		affinity = cr.Spec.DefaultBackend.Affinity
	}

	annotations := map[string]string{}
	if cr.Spec.DefaultBackend.Annotations != nil {
		annotations = *cr.Spec.DefaultBackend.Annotations
	}

	// add custom arguments from CR
	args := append([]string{}, cr.Spec.DefaultBackend.CustomArgs...)

	resourcesLimits := corev1.ResourceList{}
	if cr.Spec.DefaultBackend.PodLimits != nil {
		resourcesLimits = *cr.Spec.DefaultBackend.PodLimits
	}

	resourcesRequests := corev1.ResourceList{}
	if cr.Spec.DefaultBackend.PodRequests != nil {
		resourcesRequests = *cr.Spec.DefaultBackend.PodRequests
	}

	ports := []corev1.ContainerPort{
		{
			Name:          "http",
			ContainerPort: cr.Spec.DefaultBackend.Port,
			Protocol:      corev1.ProtocolTCP,
		},
	}

	deployment := v1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: cr.ObjectMeta.Namespace,
			Labels:    baseLabels(cr),
		},
		Spec: v1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: baseLabels(cr),
			},
			Replicas: &cr.Spec.DefaultBackend.Replicas,
			Strategy: v1.DeploymentStrategy{
				Type:          v1.RollingUpdateDeploymentStrategyType,
				RollingUpdate: nil,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: mergeMaps(baseLabels(cr),
						map[string]string{"app.improvado.io/component": "defaultbackend"},
					),
					Annotations: setAnnotations(cr, annotations),
				},
				Spec: corev1.PodSpec{
					Affinity:           affinity,
					ServiceAccountName: cr.Spec.ServiceAccount,
					SecurityContext: &corev1.PodSecurityContext{
						RunAsUser: &runAsUser,
					},
					Containers: []corev1.Container{
						{
							Name:  "nginx-ingress",
							Image: cr.Spec.DefaultBackend.Image.Repository + ":" + cr.Spec.DefaultBackend.Image.Tag,
							Args:  args,
							Ports: ports,
							Env:   env,
							Resources: corev1.ResourceRequirements{
								Limits:   resourcesLimits,
								Requests: resourcesRequests,
							},
						},
					},
				},
			},
		},
	}

	return deployment
}

func generateDefaultBackendService(cr *appv1alpha1.NginxIngress) corev1.Service {
	labels := map[string]string{
		"app.improvado.io/component": "service",
	}

	serviceAnnotations := map[string]string{}
	if cr.Spec.DefaultBackend.Annotations != nil {
		serviceAnnotations = *cr.Spec.DefaultBackend.ServiceAnnotations
	}

	name := cr.Name + "-default-backend"
	if cr.Spec.DefaultBackend.Name != "" {
		name = cr.Spec.DefaultBackend.Name
	}

	service := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   cr.Namespace,
			Labels:      labels,
			Annotations: serviceAnnotations,
		},
	}

	service.Spec.Type = corev1.ServiceTypeLoadBalancer

	service.Spec.Ports = []corev1.ServicePort{
		{
			Name:       "http",
			Port:       cr.Spec.DefaultBackend.Port,
			TargetPort: intstr.FromString("http"),
			Protocol:   corev1.ProtocolTCP,
		},
	}

	service.Spec.Selector = baseLabels(cr)

	return service
}
