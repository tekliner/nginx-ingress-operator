package nginxingress

import (
	"reflect"

	"k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

func reconcileDeployment(foundDeployment v1.Deployment, newDeployment v1.Deployment) (bool, v1.Deployment) {

	reconcileRequired := false

	if !reflect.DeepEqual(foundDeployment.Spec, newDeployment.Spec) {
		foundDeployment.Spec.Replicas = newDeployment.Spec.Replicas
		foundDeployment.Spec.Template = newDeployment.Spec.Template
		reconcileRequired = true
	}

	if !reflect.DeepEqual(foundDeployment.Annotations, newDeployment.Annotations) {
		foundDeployment.Annotations = newDeployment.Annotations
		reconcileRequired = true
	}

	if !reflect.DeepEqual(foundDeployment.Labels, newDeployment.Labels) {
		foundDeployment.Labels = newDeployment.Labels
		reconcileRequired = true
	}

	return reconcileRequired, foundDeployment

}

func reconcileService(foundService corev1.Service, newService corev1.Service) (bool, corev1.Service) {

	reconcileRequired := false

	if !reflect.DeepEqual(foundService.Spec, newService.Spec) {
		foundService.Spec.Ports = newService.Spec.Ports
		reconcileRequired = true
	}

	if !reflect.DeepEqual(foundService.Annotations, newService.Annotations) {
		foundService.Annotations = newService.Annotations
		reconcileRequired = true
	}

	if !reflect.DeepEqual(foundService.Labels, newService.Labels) {
		foundService.Labels = newService.Labels
		reconcileRequired = true
	}

	return reconcileRequired, foundService
}

func reconcileConfigmap(foundConfigmap corev1.ConfigMap, newConfigmap corev1.ConfigMap) (bool, corev1.ConfigMap) {
	reconcileRequired := false

	if !reflect.DeepEqual(foundConfigmap.Data, newConfigmap.Data) {
		foundConfigmap.Data = newConfigmap.Data
		reconcileRequired = true
	}

	if !reflect.DeepEqual(foundConfigmap.Annotations, newConfigmap.Annotations) {
		foundConfigmap.Annotations = newConfigmap.Annotations
		reconcileRequired = true
	}

	if !reflect.DeepEqual(foundConfigmap.Labels, newConfigmap.Labels) {
		foundConfigmap.Labels = newConfigmap.Labels
		reconcileRequired = true
	}

	return reconcileRequired, foundConfigmap
}
