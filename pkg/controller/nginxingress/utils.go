package nginxingress

import (
	appv1alpha1 "github.com/tekliner/nginx-ingress-operator/pkg/apis/app/v1alpha1"
	v1 "k8s.io/api/core/v1"
)

func mergeMaps(itermaps ...map[string]string) map[string]string {
	result := make(map[string]string)
	for _, rv := range itermaps {
		for k, v := range rv {
			result[k] = v
		}
	}
	return result
}

func baseLabels(cr *appv1alpha1.NginxIngress) map[string]string {
	return map[string]string{
		"app.improvado.io/application": "nginx-ingress-controller",
		"app.improvado.io/instance":    cr.Name,
	}
}

func returnDefaultAnnotations(cr *appv1alpha1.NginxIngress) map[string]string {
	return map[string]string{}
}

func returnDefaultENV() []v1.EnvVar {
	return []v1.EnvVar{
		v1.EnvVar{
			Name: "POD_NAME",
			ValueFrom: &v1.EnvVarSource{
				FieldRef: &v1.ObjectFieldSelector{FieldPath: "metadata.name"},
			},
		},
		v1.EnvVar{
			Name: "POD_NAMESPACE",
			ValueFrom: &v1.EnvVarSource{
				FieldRef: &v1.ObjectFieldSelector{FieldPath: "metadata.namespace"},
			},
		},
	}
}

func setAnnotations(cr *appv1alpha1.NginxIngress, templateAnnotations map[string]string) map[string]string {
	annotations := returnDefaultAnnotations(cr)
	if len(templateAnnotations) > 0 {
		annotations = templateAnnotations
	}
	return annotations
}
