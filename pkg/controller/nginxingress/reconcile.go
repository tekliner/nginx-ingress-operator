package nginxingress

import (
	"k8s.io/api/policy/v1beta1"
	"reflect"

	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

func reconcileDeployment(foundDeployment v1.Deployment, newDeployment v1.Deployment) (bool, v1.Deployment) {

	reconcileRequired := false

	if !reflect.DeepEqual(foundDeployment.Spec.Replicas, newDeployment.Spec.Replicas) {
		foundDeployment.Spec.Replicas = newDeployment.Spec.Replicas
		reconcileRequired = true
	}

	if !reflect.DeepEqual(foundDeployment.Labels, newDeployment.Labels) {
		foundDeployment.Labels = newDeployment.Labels
		reconcileRequired = true
	}

	if !reflect.DeepEqual(foundDeployment.Spec.Template, newDeployment.Spec.Template) {
		foundDeployment.Spec.Template = newDeployment.Spec.Template
		reconcileRequired = true
	}

	return reconcileRequired, foundDeployment

}
func reconcilePdb(foundPdb v1beta1.PodDisruptionBudget, newPdb v1beta1.PodDisruptionBudget) (bool, v1beta1.PodDisruptionBudget) {

	reconcileRequired := false
	if !reflect.DeepEqual(foundPdb.Spec.Selector, newPdb.Spec.Selector) {
		foundPdb.Spec.Selector = newPdb.Spec.Selector
		reconcileRequired = true
	}

	if !reflect.DeepEqual(foundPdb.Spec.MaxUnavailable, newPdb.Spec.MaxUnavailable) {
		foundPdb.Spec.MaxUnavailable = newPdb.Spec.MaxUnavailable
		reconcileRequired = true
	}

	if !reflect.DeepEqual(foundPdb.Spec.MinAvailable, newPdb.Spec.MinAvailable) {
		foundPdb.Spec.MinAvailable = newPdb.Spec.MinAvailable
		reconcileRequired = true
	}
	return reconcileRequired, foundPdb
}

func reconcileService(foundService corev1.Service, newService corev1.Service) (bool, corev1.Service) {

	reconcileRequired := false

	if !reflect.DeepEqual(foundService.Spec.Type, newService.Spec.Type) {
		foundService.Spec.Type = newService.Spec.Type
		reconcileRequired = true
	}

	if !reflect.DeepEqual(foundService.Spec.ExternalTrafficPolicy, newService.Spec.ExternalTrafficPolicy) {
		foundService.Spec.Type = newService.Spec.Type
		reconcileRequired = true
	}

	if !reflect.DeepEqual(foundService.Labels, newService.Labels) {
		foundService.Labels = newService.Labels
		reconcileRequired = true
	}

	if !reflect.DeepEqual(foundService.Annotations, newService.Annotations) {
		foundService.Annotations = newService.Annotations
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

	if !reflect.DeepEqual(foundConfigmap.Labels, newConfigmap.Labels) {
		foundConfigmap.Labels = newConfigmap.Labels
		reconcileRequired = true
	}

	return reconcileRequired, foundConfigmap
}
