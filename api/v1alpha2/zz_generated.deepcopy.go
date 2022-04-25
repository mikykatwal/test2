//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright 2022 Mondoo, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha2

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Admission) DeepCopyInto(out *Admission) {
	*out = *in
	out.Image = in.Image
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Admission.
func (in *Admission) DeepCopy() *Admission {
	if in == nil {
		return nil
	}
	out := new(Admission)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CertificateProvisioning) DeepCopyInto(out *CertificateProvisioning) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CertificateProvisioning.
func (in *CertificateProvisioning) DeepCopy() *CertificateProvisioning {
	if in == nil {
		return nil
	}
	out := new(CertificateProvisioning)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Image) DeepCopyInto(out *Image) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Image.
func (in *Image) DeepCopy() *Image {
	if in == nil {
		return nil
	}
	out := new(Image)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KubernetesResources) DeepCopyInto(out *KubernetesResources) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KubernetesResources.
func (in *KubernetesResources) DeepCopy() *KubernetesResources {
	if in == nil {
		return nil
	}
	out := new(KubernetesResources)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MondooAuditConfig) DeepCopyInto(out *MondooAuditConfig) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MondooAuditConfig.
func (in *MondooAuditConfig) DeepCopy() *MondooAuditConfig {
	if in == nil {
		return nil
	}
	out := new(MondooAuditConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *MondooAuditConfig) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MondooAuditConfigCondition) DeepCopyInto(out *MondooAuditConfigCondition) {
	*out = *in
	in.LastUpdateTime.DeepCopyInto(&out.LastUpdateTime)
	in.LastTransitionTime.DeepCopyInto(&out.LastTransitionTime)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MondooAuditConfigCondition.
func (in *MondooAuditConfigCondition) DeepCopy() *MondooAuditConfigCondition {
	if in == nil {
		return nil
	}
	out := new(MondooAuditConfigCondition)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MondooAuditConfigData) DeepCopyInto(out *MondooAuditConfigData) {
	*out = *in
	out.CertificateProvisioning = in.CertificateProvisioning
	in.Scanner.DeepCopyInto(&out.Scanner)
	out.KubernetesResources = in.KubernetesResources
	out.Nodes = in.Nodes
	out.Admission = in.Admission
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MondooAuditConfigData.
func (in *MondooAuditConfigData) DeepCopy() *MondooAuditConfigData {
	if in == nil {
		return nil
	}
	out := new(MondooAuditConfigData)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MondooAuditConfigList) DeepCopyInto(out *MondooAuditConfigList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]MondooAuditConfig, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MondooAuditConfigList.
func (in *MondooAuditConfigList) DeepCopy() *MondooAuditConfigList {
	if in == nil {
		return nil
	}
	out := new(MondooAuditConfigList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *MondooAuditConfigList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MondooAuditConfigStatus) DeepCopyInto(out *MondooAuditConfigStatus) {
	*out = *in
	if in.Pods != nil {
		in, out := &in.Pods, &out.Pods
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]MondooAuditConfigCondition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MondooAuditConfigStatus.
func (in *MondooAuditConfigStatus) DeepCopy() *MondooAuditConfigStatus {
	if in == nil {
		return nil
	}
	out := new(MondooAuditConfigStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Nodes) DeepCopyInto(out *Nodes) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Nodes.
func (in *Nodes) DeepCopy() *Nodes {
	if in == nil {
		return nil
	}
	out := new(Nodes)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Scanner) DeepCopyInto(out *Scanner) {
	*out = *in
	out.Image = in.Image
	in.Resources.DeepCopyInto(&out.Resources)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Scanner.
func (in *Scanner) DeepCopy() *Scanner {
	if in == nil {
		return nil
	}
	out := new(Scanner)
	in.DeepCopyInto(out)
	return out
}
