// Code generated by protoc-gen-deepcopy. DO NOT EDIT.
package selectorv1

import (
	proto "github.com/golang/protobuf/proto"
)

// DeepCopyInto supports using Selector within kubernetes types, where deepcopy-gen is used.
func (in *Selector) DeepCopyInto(out *Selector) {
	p := proto.Clone(in).(*Selector)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Selector. Required by controller-gen.
func (in *Selector) DeepCopy() *Selector {
	if in == nil {
		return nil
	}
	out := new(Selector)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new Selector. Required by controller-gen.
func (in *Selector) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using ServiceSelector within kubernetes types, where deepcopy-gen is used.
func (in *ServiceSelector) DeepCopyInto(out *ServiceSelector) {
	p := proto.Clone(in).(*ServiceSelector)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ServiceSelector. Required by controller-gen.
func (in *ServiceSelector) DeepCopy() *ServiceSelector {
	if in == nil {
		return nil
	}
	out := new(ServiceSelector)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new ServiceSelector. Required by controller-gen.
func (in *ServiceSelector) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using FlowSelector within kubernetes types, where deepcopy-gen is used.
func (in *FlowSelector) DeepCopyInto(out *FlowSelector) {
	p := proto.Clone(in).(*FlowSelector)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FlowSelector. Required by controller-gen.
func (in *FlowSelector) DeepCopy() *FlowSelector {
	if in == nil {
		return nil
	}
	out := new(FlowSelector)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new FlowSelector. Required by controller-gen.
func (in *FlowSelector) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using ControlPoint within kubernetes types, where deepcopy-gen is used.
func (in *ControlPoint) DeepCopyInto(out *ControlPoint) {
	p := proto.Clone(in).(*ControlPoint)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ControlPoint. Required by controller-gen.
func (in *ControlPoint) DeepCopy() *ControlPoint {
	if in == nil {
		return nil
	}
	out := new(ControlPoint)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new ControlPoint. Required by controller-gen.
func (in *ControlPoint) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}
