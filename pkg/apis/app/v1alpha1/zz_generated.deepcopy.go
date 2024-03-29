// +build !ignore_autogenerated

// Code generated by operator-sdk. DO NOT EDIT.

package v1alpha1

import (
	v1 "k8s.io/api/core/v1"
	resource "k8s.io/apimachinery/pkg/api/resource"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DefaultBackendSpec) DeepCopyInto(out *DefaultBackendSpec) {
	*out = *in
	in.Image.DeepCopyInto(&out.Image)
	if in.CustomArgs != nil {
		in, out := &in.CustomArgs, &out.CustomArgs
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.RunAsUser != nil {
		in, out := &in.RunAsUser, &out.RunAsUser
		*out = new(int64)
		**out = **in
	}
	if in.Env != nil {
		in, out := &in.Env, &out.Env
		*out = new([]v1.EnvVar)
		if **in != nil {
			in, out := *in, *out
			*out = make([]v1.EnvVar, len(*in))
			for i := range *in {
				(*in)[i].DeepCopyInto(&(*out)[i])
			}
		}
	}
	if in.Port != nil {
		in, out := &in.Port, &out.Port
		*out = new(int32)
		**out = **in
	}
	if in.Affinity != nil {
		in, out := &in.Affinity, &out.Affinity
		*out = new(v1.Affinity)
		(*in).DeepCopyInto(*out)
	}
	if in.Replicas != nil {
		in, out := &in.Replicas, &out.Replicas
		*out = new(int32)
		**out = **in
	}
	if in.Annotations != nil {
		in, out := &in.Annotations, &out.Annotations
		*out = new(map[string]string)
		if **in != nil {
			in, out := *in, *out
			*out = make(map[string]string, len(*in))
			for key, val := range *in {
				(*out)[key] = val
			}
		}
	}
	if in.PodRequests != nil {
		in, out := &in.PodRequests, &out.PodRequests
		*out = new(v1.ResourceList)
		if **in != nil {
			in, out := *in, *out
			*out = make(map[v1.ResourceName]resource.Quantity, len(*in))
			for key, val := range *in {
				(*out)[key] = val.DeepCopy()
			}
		}
	}
	if in.PodLimits != nil {
		in, out := &in.PodLimits, &out.PodLimits
		*out = new(v1.ResourceList)
		if **in != nil {
			in, out := *in, *out
			*out = make(map[v1.ResourceName]resource.Quantity, len(*in))
			for key, val := range *in {
				(*out)[key] = val.DeepCopy()
			}
		}
	}
	if in.ServiceAnnotations != nil {
		in, out := &in.ServiceAnnotations, &out.ServiceAnnotations
		*out = new(map[string]string)
		if **in != nil {
			in, out := *in, *out
			*out = make(map[string]string, len(*in))
			for key, val := range *in {
				(*out)[key] = val
			}
		}
	}
	in.Pdb.DeepCopyInto(&out.Pdb)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DefaultBackendSpec.
func (in *DefaultBackendSpec) DeepCopy() *DefaultBackendSpec {
	if in == nil {
		return nil
	}
	out := new(DefaultBackendSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ImageSpec) DeepCopyInto(out *ImageSpec) {
	*out = *in
	if in.PullPolicy != nil {
		in, out := &in.PullPolicy, &out.PullPolicy
		*out = new(v1.PullPolicy)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ImageSpec.
func (in *ImageSpec) DeepCopy() *ImageSpec {
	if in == nil {
		return nil
	}
	out := new(ImageSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IngressServiceSpec) DeepCopyInto(out *IngressServiceSpec) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IngressServiceSpec.
func (in *IngressServiceSpec) DeepCopy() *IngressServiceSpec {
	if in == nil {
		return nil
	}
	out := new(IngressServiceSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MetricsServiceSpecs) DeepCopyInto(out *MetricsServiceSpecs) {
	*out = *in
	if in.Annotations != nil {
		in, out := &in.Annotations, &out.Annotations
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MetricsServiceSpecs.
func (in *MetricsServiceSpecs) DeepCopy() *MetricsServiceSpecs {
	if in == nil {
		return nil
	}
	out := new(MetricsServiceSpecs)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NginxControllerSpec) DeepCopyInto(out *NginxControllerSpec) {
	*out = *in
	in.Image.DeepCopyInto(&out.Image)
	if in.Env != nil {
		in, out := &in.Env, &out.Env
		*out = new([]v1.EnvVar)
		if **in != nil {
			in, out := *in, *out
			*out = make([]v1.EnvVar, len(*in))
			for i := range *in {
				(*in)[i].DeepCopyInto(&(*out)[i])
			}
		}
	}
	if in.CustomArgs != nil {
		in, out := &in.CustomArgs, &out.CustomArgs
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.RunAsUser != nil {
		in, out := &in.RunAsUser, &out.RunAsUser
		*out = new(int64)
		**out = **in
	}
	if in.Affinity != nil {
		in, out := &in.Affinity, &out.Affinity
		*out = new(v1.Affinity)
		(*in).DeepCopyInto(*out)
	}
	if in.PodRequests != nil {
		in, out := &in.PodRequests, &out.PodRequests
		*out = new(v1.ResourceList)
		if **in != nil {
			in, out := *in, *out
			*out = make(map[v1.ResourceName]resource.Quantity, len(*in))
			for key, val := range *in {
				(*out)[key] = val.DeepCopy()
			}
		}
	}
	if in.PodLimits != nil {
		in, out := &in.PodLimits, &out.PodLimits
		*out = new(v1.ResourceList)
		if **in != nil {
			in, out := *in, *out
			*out = make(map[v1.ResourceName]resource.Quantity, len(*in))
			for key, val := range *in {
				(*out)[key] = val.DeepCopy()
			}
		}
	}
	if in.Annotations != nil {
		in, out := &in.Annotations, &out.Annotations
		*out = new(map[string]string)
		if **in != nil {
			in, out := *in, *out
			*out = make(map[string]string, len(*in))
			for key, val := range *in {
				(*out)[key] = val
			}
		}
	}
	if in.Config != nil {
		in, out := &in.Config, &out.Config
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Labels != nil {
		in, out := &in.Labels, &out.Labels
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	in.Pdb.DeepCopyInto(&out.Pdb)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NginxControllerSpec.
func (in *NginxControllerSpec) DeepCopy() *NginxControllerSpec {
	if in == nil {
		return nil
	}
	out := new(NginxControllerSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NginxIngress) DeepCopyInto(out *NginxIngress) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NginxIngress.
func (in *NginxIngress) DeepCopy() *NginxIngress {
	if in == nil {
		return nil
	}
	out := new(NginxIngress)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *NginxIngress) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NginxIngressList) DeepCopyInto(out *NginxIngressList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]NginxIngress, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NginxIngressList.
func (in *NginxIngressList) DeepCopy() *NginxIngressList {
	if in == nil {
		return nil
	}
	out := new(NginxIngressList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *NginxIngressList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NginxIngressSpec) DeepCopyInto(out *NginxIngressSpec) {
	*out = *in
	if in.Metrics != nil {
		in, out := &in.Metrics, &out.Metrics
		*out = new(MetricsServiceSpecs)
		(*in).DeepCopyInto(*out)
	}
	if in.Stats != nil {
		in, out := &in.Stats, &out.Stats
		*out = new(StatsSpec)
		**out = **in
	}
	in.NginxController.DeepCopyInto(&out.NginxController)
	if in.DefaultBackend != nil {
		in, out := &in.DefaultBackend, &out.DefaultBackend
		*out = new(DefaultBackendSpec)
		(*in).DeepCopyInto(*out)
	}
	in.NginxServiceSpec.DeepCopyInto(&out.NginxServiceSpec)
	if in.NginxServiceAnnotations != nil {
		in, out := &in.NginxServiceAnnotations, &out.NginxServiceAnnotations
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NginxIngressSpec.
func (in *NginxIngressSpec) DeepCopy() *NginxIngressSpec {
	if in == nil {
		return nil
	}
	out := new(NginxIngressSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NginxIngressStatus) DeepCopyInto(out *NginxIngressStatus) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NginxIngressStatus.
func (in *NginxIngressStatus) DeepCopy() *NginxIngressStatus {
	if in == nil {
		return nil
	}
	out := new(NginxIngressStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *StatsSpec) DeepCopyInto(out *StatsSpec) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StatsSpec.
func (in *StatsSpec) DeepCopy() *StatsSpec {
	if in == nil {
		return nil
	}
	out := new(StatsSpec)
	in.DeepCopyInto(out)
	return out
}
