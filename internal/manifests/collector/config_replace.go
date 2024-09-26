// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package collector

import (
	"encoding/json"
	"fmt"
	"time"

	ta "github.com/aws/amazon-cloudwatch-agent-operator/internal/manifests/targetallocator/adapters"
	"github.com/aws/amazon-cloudwatch-agent-operator/internal/naming"
	"github.com/aws/amazon-cloudwatch-agent-operator/pkg/featuregate"

	promconfig "github.com/prometheus/prometheus/config"
	_ "github.com/prometheus/prometheus/discovery/install" // Package install has the side-effect of registering all builtin.
	"gopkg.in/yaml.v2"

	"github.com/aws/amazon-cloudwatch-agent-operator/apis/v1alpha1"
	"github.com/aws/amazon-cloudwatch-agent-operator/internal/manifests/collector/adapters"
)

type targetAllocator struct {
	Endpoint string        `yaml:"endpoint"`
	Interval time.Duration `yaml:"interval"`
	// HTTPSDConfig is a preference that can be set for the collector's target allocator, but the operator doesn't
	// care about what the value is set to. We just need this for validation when unmarshalling the configmap.
	HTTPSDConfig interface{} `yaml:"http_sd_config,omitempty"`
}

type Config struct {
	PromConfig        *promconfig.Config `yaml:"config"`
	TargetAllocConfig *targetAllocator   `yaml:"target_allocator,omitempty"`
}

func ReplaceConfig(instance v1alpha1.AmazonCloudWatchAgent) (string, error) {
	config, err := adapters.ConfigFromJSONString(instance.Spec.Config)
	if err != nil {
		return "", err
	}

	out, err := json.Marshal(config)
	if err != nil {
		return "", err
	}

	return string(out), nil
}

// ReplacePrometheusConfig replaces the prometheus configuration that the customer provides with itself (if the
// target-allocator isn't enabled) or the target_allocator configuration (if the target-allocator is enabled)
// and populates it into the prometheus.yaml file, which is seen in its ConfigMap.
func ReplacePrometheusConfig(instance v1alpha1.AmazonCloudWatchAgent) (string, error) {
	promConfigYaml, err := instance.Spec.Prometheus.Yaml()
	if err != nil {
		return "", fmt.Errorf("%s could not convert json to yaml", err)
	}

	// Check if TargetAllocator is enabled, if not, return the original config
	if !instance.Spec.TargetAllocator.Enabled {
		prometheusConfig, err := adapters.ConfigFromString(promConfigYaml)
		if err != nil {
			return "", err
		}

		prometheusConfigYAML, err := yaml.Marshal(prometheusConfig)
		if err != nil {
			return "", err
		}

		return string(prometheusConfigYAML), nil
	}

	promCfgMap, getCfgPromErr := adapters.ConfigFromString(promConfigYaml)
	if getCfgPromErr != nil {
		return "", getCfgPromErr
	}

	validateCfgPromErr := ta.ValidatePromConfig(promCfgMap, instance.Spec.TargetAllocator.Enabled, featuregate.EnableTargetAllocatorRewrite.IsEnabled())
	if validateCfgPromErr != nil {
		return "", validateCfgPromErr
	}

	if featuregate.EnableTargetAllocatorRewrite.IsEnabled() {
		updPromCfgMap, getCfgPromErr := ta.AddTAConfigToPromConfig(promCfgMap, naming.TAService(instance.Name))
		if getCfgPromErr != nil {
			return "", getCfgPromErr
		}

		out, updCfgMarshalErr := yaml.Marshal(updPromCfgMap)
		if updCfgMarshalErr != nil {
			return "", updCfgMarshalErr
		}

		return string(out), nil
	}

	updPromCfgMap, err := ta.AddHTTPSDConfigToPromConfig(promCfgMap, naming.TAService(instance.Name))
	if err != nil {
		return "", err
	}

	out, err := yaml.Marshal(updPromCfgMap)
	if err != nil {
		return "", err
	}

	return string(out), nil
}
