// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	v1 "k8s.io/api/core/v1"
)

// PortsSpec defines the AmazonCloudWatchAgent's container/service ports additional specifications.
type PortsSpec struct {
	// Allows defining which port to bind to the host in the Container.
	// +optional
	HostPort int32 `json:"hostPort,omitempty"`

	// Maintain previous fields in new struct
	v1.ServicePort `json:",inline"`
}

type AmazonCloudWatchAgentCommonFields struct {
	// ManagementState defines if the CR should be managed by the operator or not.
	// Default is managed.
	//
	// +required
	// +kubebuilder:validation:Required
	// +kubebuilder:default:=managed
	ManagementState ManagementStateType `json:"managementState,omitempty"`
	// Resources to set on generated pods.
	// +optional
	Resources v1.ResourceRequirements `json:"resources,omitempty"`
	// NodeSelector to schedule generated pods.
	// This only works with the following AmazonCloudWatchAgent mode's: daemonset, statefulset, and deployment.
	// +optional
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`
	// Args is the set of arguments to pass to the main container's binary.
	// +optional
	Args map[string]string `json:"args,omitempty"`
	// Replicas is the number of pod instances for the underlying replicaset. Set this if you are not using autoscaling.
	// +optional
	Replicas *int32 `json:"replicas,omitempty"`
	// PodDisruptionBudget specifies the pod disruption budget configuration to use
	// for the generated workload. By default, a PDB with a MaxUnavailable of one is set.
	//
	// +optional
	PodDisruptionBudget *PodDisruptionBudgetSpec `json:"podDisruptionBudget,omitempty"`
	// SecurityContext configures the container security context for
	// the generated main container.
	//
	// In deployment, daemonset, or statefulset mode, this controls
	// the security context settings for the primary application
	// container.
	//
	// In sidecar mode, this controls the security context for the
	// injected sidecar container.
	//
	// +optional
	SecurityContext *v1.SecurityContext `json:"securityContext,omitempty"`
	// PodSecurityContext configures the pod security context for the
	// generated pod, when running as a deployment, daemonset,
	// or statefulset.
	//
	// In sidecar mode, the amazon-cloudwatch-agent-operator will ignore this setting.
	//
	// +optional
	PodSecurityContext *v1.PodSecurityContext `json:"podSecurityContext,omitempty"`
	// PodAnnotations is the set of annotations that will be attached to
	// the generated pods.
	// +optional
	PodAnnotations map[string]string `json:"podAnnotations,omitempty"`
	// ServiceAccount indicates the name of an existing service account to use with this instance. When set,
	// the operator will not automatically create a ServiceAccount.
	// +optional
	ServiceAccount string `json:"serviceAccount,omitempty"`
	// Image indicates the container image to use for the generated pods.
	// +optional
	Image string `json:"image,omitempty"`
	// ImagePullPolicy indicates the pull policy to be used for retrieving the container image.
	// +optional
	ImagePullPolicy v1.PullPolicy `json:"imagePullPolicy,omitempty"`
	// VolumeMounts represents the mount points to use in the underlying deployment(s).
	// +optional
	// +listType=atomic
	VolumeMounts []v1.VolumeMount `json:"volumeMounts,omitempty"`
	// Ports allows a set of ports to be exposed by the underlying v1.Service & v1.ContainerPort. By default, the operator
	// will attempt to infer the required ports by parsing the .Spec.Config property but this property can be
	// used to open additional ports that can't be inferred by the operator, like for custom receivers.
	// +optional
	// +listType=atomic
	Ports []PortsSpec `json:"ports,omitempty"`
	// Environment variables to set on the generated pods.
	// +optional
	Env []v1.EnvVar `json:"env,omitempty"`
	// List of sources to populate environment variables on the generated pods.
	// +optional
	EnvFrom []v1.EnvFromSource `json:"envFrom,omitempty"`
	// Toleration to schedule the generated pods.
	// This only works with the following AmazonCloudWatchAgent mode's: daemonset, statefulset, and deployment.
	// +optional
	Tolerations []v1.Toleration `json:"tolerations,omitempty"`
	// Volumes represents which volumes to use in the underlying deployment(s).
	// +optional
	// +listType=atomic
	Volumes []v1.Volume `json:"volumes,omitempty"`
	// If specified, indicates the pod's scheduling constraints
	// +optional
	Affinity *v1.Affinity `json:"affinity,omitempty"`
	// Actions that the management system should take in response to container lifecycle events. Cannot be updated.
	// +optional
	Lifecycle *v1.Lifecycle `json:"lifecycle,omitempty"`
	// Duration in seconds the pod needs to terminate gracefully upon probe failure.
	// +optional
	TerminationGracePeriodSeconds *int64 `json:"terminationGracePeriodSeconds,omitempty"`
	// TopologySpreadConstraints embedded kubernetes pod configuration option,
	// controls how pods are spread across your cluster among failure-domains
	// such as regions, zones, nodes, and other user-defined topology domains
	// https://kubernetes.io/docs/concepts/workloads/pods/pod-topology-spread-constraints/
	// This only works with the following AmazonCloudWatchAgent mode's: statefulset, and deployment.
	// +optional
	TopologySpreadConstraints []v1.TopologySpreadConstraint `json:"topologySpreadConstraints,omitempty"`
	// HostNetwork indicates if the pod should run in the host networking namespace.
	// +optional
	HostNetwork bool `json:"hostNetwork,omitempty"`
	// ShareProcessNamespace indicates if the pod's containers should share process namespace.
	// +optional
	ShareProcessNamespace bool `json:"shareProcessNamespace,omitempty"`
	// If specified, indicates the pod's priority.
	// If not specified, the pod priority will be default or zero if there is no
	// default.
	// +optional
	PriorityClassName string `json:"priorityClassName,omitempty"`
	// InitContainers allows injecting initContainers to the generated pod definition.
	// These init containers can be used to fetch secrets for injection into the
	// configuration from external sources, run added checks, etc. Any errors during the execution of
	// an initContainer will lead to a restart of the Pod. More info:
	// https://kubernetes.io/docs/concepts/workloads/pods/init-containers/
	// +optional
	InitContainers []v1.Container `json:"initContainers,omitempty"`
	// AdditionalContainers allows injecting additional containers into the generated pod definition.
	// These sidecar containers can be used for authentication proxies, log shipping sidecars, agents for shipping
	// metrics to their cloud, or in general sidecars that do not support automatic injection.
	// This only works with the following AmazonCloudWatchAgent mode's: daemonset, statefulset, and deployment.
	//
	// Container names managed by the operator:
	// * `otc-container`
	//
	// Overriding containers managed by the operator is outside the scope of what the maintainers will support and by
	// doing so, you wil accept the risk of it breaking things.
	//
	// +optional
	AdditionalContainers []v1.Container `json:"additionalContainers,omitempty"`
	// PodDNSConfig defines the DNS parameters of a pod in addition to those generated from DNSPolicy.
	PodDNSConfig v1.PodDNSConfig `json:"podDnsConfig,omitempty"`
	// IPFamily represents the IP Family (IPv4 or IPv6). This type is used
	// to express the family of an IP expressed by a type (e.g. service.spec.ipFamilies).
	// +optional
	IpFamilies []v1.IPFamily `json:"ipFamilies,omitempty"`
	// IPFamilyPolicy represents the dual-stack-ness requested or required by a Service
	// +kubebuilder:default:=SingleStack
	// +optional
	IpFamilyPolicy *v1.IPFamilyPolicy `json:"ipFamilyPolicy,omitempty"`
}

type StatefulSetCommonFields struct {
	// VolumeClaimTemplates will provide stable storage using PersistentVolumes.
	// This only works with the following AmazonCloudWatchAgent mode's: statefulset.
	// +optional
	// +listType=atomic
	VolumeClaimTemplates []v1.PersistentVolumeClaim `json:"volumeClaimTemplates,omitempty"`
}
