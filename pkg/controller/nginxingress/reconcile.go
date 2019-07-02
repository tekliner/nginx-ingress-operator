package nginxingress

import (
	"reflect"

	"k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

func reconcileDeployment(foundDeployment v1.Deployment, newDeployment v1.Deployment) (bool, v1.Deployment) {

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

	return reconcileDeployment, foundDeployment

}

func reconcileService(foundService corev1.Service, newService corev1.Service) (bool, corev1.Service) {

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

	return reconcileService, foundService
}

func reconcileConfigmap(foundConfigmap corev1.ConfigMap, newConfigmap corev1.ConfigMap) (bool, corev1.ConfigMap) {
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

	return reconcileConfigmap, foundConfigmap
}
