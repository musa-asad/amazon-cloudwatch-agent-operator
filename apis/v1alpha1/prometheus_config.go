// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	"bytes"

	"gopkg.in/yaml.v3"
)

// AnyConfig represent parts of the config.
type AnyConfig struct {
	Object map[string]interface{} `json:"-" yaml:",inline"`
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AnyConfig) DeepCopyInto(out *AnyConfig) {
	*out = *in
	if in.Object != nil {
		in, out := &in.Object, &out.Object
		*out = make(map[string]interface{}, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AnyConfig.
func (in *AnyConfig) DeepCopy() *AnyConfig {
	if in == nil {
		return nil
	}
	out := new(AnyConfig)
	in.DeepCopyInto(out)
	return out
}

// PrometheusConfig encapsulates prometheus config.
type PrometheusConfig struct {
	// +kubebuilder:pruning:PreserveUnknownFields
	Config *AnyConfig `json:"config,omitempty" yaml:"config,omitempty"`
	// +kubebuilder:pruning:PreserveUnknownFields
	TrimMetricSuffixes bool `json:"trim_metric_suffixes,omitempty" yaml:"trim_metric_suffixes,omitempty"`
	// +kubebuilder:pruning:PreserveUnknownFields
	UseStartTimeMetric bool `json:"use_start_time_metric,omitempty" yaml:"use_start_time_metric,omitempty"`
	// +kubebuilder:pruning:PreserveUnknownFields
	StartTimeMetricRegex string `json:"start_time_metric_regex,omitempty" yaml:"start_time_metric_regex,omitempty"`
	// +kubebuilder:pruning:PreserveUnknownFields
	ReportExtraScrapeMetrics bool `json:"report_extra_scrape_metrics,omitempty" yaml:"report_extra_scrape_metrics,omitempty"`
	// +kubebuilder:pruning:PreserveUnknownFields
	TargetAllocator *AnyConfig `json:"target_allocator,omitempty" yaml:"target_allocator,omitempty"`
}

// Yaml encodes the current object and returns it as a string.
func (pc *PrometheusConfig) Yaml() (string, error) {
	var buf bytes.Buffer
	yamlEncoder := yaml.NewEncoder(&buf)
	yamlEncoder.SetIndent(2)
	if err := yamlEncoder.Encode(&pc); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// IsEmpty checks if the prometheus config is empty.
func (pc *PrometheusConfig) IsEmpty() bool {
	return pc.Config == nil &&
		!pc.TrimMetricSuffixes &&
		!pc.UseStartTimeMetric &&
		pc.StartTimeMetricRegex == "" &&
		!pc.ReportExtraScrapeMetrics &&
		pc.TargetAllocator == nil
}
