/*
Copyright 2022.

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

package controllers

import (
	"context"
	"path/filepath"

	agentv1alpha1 "github.com/fluxninja/aperture/operator/api/agent/v1alpha1"
	controllerv1alpha1 "github.com/fluxninja/aperture/operator/api/controller/v1alpha1"
	"k8s.io/client-go/dynamic"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	// MutatingWebhookURI defines the URI for the Mutating Webhook for Pods.
	MutatingWebhookURI = "/mutate-pod"
	// AgentMutatingWebhookURI defines the URI for the Mutating Webhook for Agents.
	AgentMutatingWebhookURI = "agent-defaulter"
	// ControllerMutatingWebhookURI defines the URI for the Mutating Webhook for Controllers.
	ControllerMutatingWebhookURI = "controller-defaulter"
	// SecretKey defines the Kubernetes secret data key.
	SecretKey = "apiKey"
	// AppName defines name of the application.
	AppName = "aperture"
	// OperatorName defines operator name.
	OperatorName = AppName + "-operator"
	// ControllerServiceName defines controller service name.
	ControllerServiceName = AppName + "-controller"
	// AgentServiceName defines agent service name.
	AgentServiceName = AppName + "-agent"
	// PodMutatingWebhookName defines agent service name.
	PodMutatingWebhookName = AppName + "-injector"
	// AgentMutatingWebhookName defines agent service name.
	AgentMutatingWebhookName = AppName + "-" + AgentMutatingWebhookURI
	// ControllerMutatingWebhookName defines Controller Mutating Webhook Name.
	ControllerMutatingWebhookName = AppName + "-" + ControllerMutatingWebhookURI
	// ValidatingWebhookName defines Validating Webhook name.
	ValidatingWebhookName = ControllerServiceName + "-webhook"
	// FinalizerName defines finalizer name.
	FinalizerName = "fluxninja.com/finalizer"
	// SidecarKey defines sidecar key.
	SidecarKey = "sidecar.fluxninja.com"
	// SidecarAnnotationKey defines sidecar annotation key.
	SidecarAnnotationKey = SidecarKey + "/injection"
	// SidecarLabelKey defines sidecar label key.
	SidecarLabelKey = AppName + "-injection"
	// AgentGroupKey defines agent group key.
	AgentGroupKey = SidecarKey + "/agent-group"
	// V1Version defines v1 version.
	V1Version = "v1"
	// V1Alpha1Version defines v1alpha1 version.
	V1Alpha1Version = "v1alpha1"
	// Enabled string.
	Enabled = "enabled"
	// ValidatingWebhookSvcName defines Validating Webhook service name.
	ValidatingWebhookSvcName = ValidatingWebhookName
	// WebhookClientCertName defines client cert name.
	WebhookClientCertName = "client.pem"
	// ControllerCertKeyName defines controller key file name.
	ControllerCertKeyName = "key.pem"
	// ControllerCertName defines controller cert name.
	ControllerCertName = "crt.pem"
	// ControllerCertPath defines controller cert path.
	ControllerCertPath = "/etc/aperture/aperture-controller/certs"
	// Server string.
	Server = "server"
	// TCP string.
	TCP = "TCP"
	// DistCache string.
	DistCache = "dist-cache"
	// MemberList string.
	MemberList = "memberlist"
	// ApertureFluxNinjaPlugin defines FluxNinja plugin name.
	ApertureFluxNinjaPlugin = "aperture-plugin-fluxninja"
	// DefaulterAnnotationKey defines annotation key for set defaults.
	DefaulterAnnotationKey = "fluxninja.com/set-defaults"
	// FailedStatus string.
	FailedStatus = "failed"
	// PolicyValidatingWebhookName defines Validating Webhook name for Policy.
	PolicyValidatingWebhookName = "policy-validator.fluxninja.com"
	// PolicyValidatingWebhookURI defines Validating Webhook URI for Policy.
	PolicyValidatingWebhookURI = "/validate/policy"
)

var (
	// PolicyFilePath defines default path for the policies on Controller.
	PolicyFilePath = filepath.Join("/", "etc", "aperture", "aperture-controller", "policies")
	// Test string.
	Test = "test"
	// TestTwo string.
	TestTwo = "test2"
	// TestArray array.
	TestArray = []string{Test}
	// TestArrayTwo array.
	TestArrayTwo = []string{TestTwo, Test}
	// TestMap map.
	TestMap = map[string]string{
		Test: Test,
	}
	// TestMapTwo map.
	TestMapTwo = map[string]string{
		Test:    Test,
		TestTwo: TestTwo,
	}
	// K8sClient defines Kubernetes client for tests.
	K8sClient client.Client
	// K8sDynamicClient defines Kubernetes Dynamic client for tests.
	K8sDynamicClient dynamic.Interface
	// K8sManager defines Kubernetes Manager for tests.
	K8sManager ctrl.Manager
	// Ctx context.
	Ctx context.Context
	// DefaultAgentInstance defines default Agent instance for tests.
	DefaultAgentInstance *agentv1alpha1.Agent
	// DefaultControllerInstance defines default Controller instance for tests.
	DefaultControllerInstance *controllerv1alpha1.Controller
	// CertDir defines cert directory for tests.
	CertDir = filepath.Join(".", "certs")
	// PoliciesDir defines policies directory for tests.
	PoliciesDir = filepath.Join(".", "policies")
)
