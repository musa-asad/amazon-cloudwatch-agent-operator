// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"

	"k8s.io/client-go/kubernetes"
)

func main() {

	args := os.Args
	namespace := args[1]
	jsonFilePath := args[2]

	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("error getting user home dir: %v\n\n", err)
	}
	kubeConfigPath := filepath.Join(userHomeDir, ".kube", "config")
	fmt.Printf("Using kubeconfig: %s\n\n", kubeConfigPath)

	kubeConfig, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		fmt.Printf("Error getting kubernetes config: %v\n\n", err)
	}

	clientSet, err := kubernetes.NewForConfig(kubeConfig)

	if err != nil {
		fmt.Printf("error getting kubernetes config: %v\n\n", err)
	}

	success := verifyInstrumentationEnvVariables(clientSet, namespace, jsonFilePath)
	if !success {
		fmt.Println("Instrumentation Annotation Injection Test: FAIL")
		os.Exit(1)
	} else {
		fmt.Println("Instrumentation Annotation Injection Test: PASS")
	}
}

func verifyInstrumentationEnvVariables(clientset *kubernetes.Clientset, namespace, jsonPath string) bool {
	podList, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: "app=nginx",
		FieldSelector: "status.phase!=Terminating",
	})
	if err != nil {
		fmt.Println("Error retrieving pods:", err)
		return false
	}

	if len(podList.Items) == 0 {
		fmt.Println("No pods found with the specified label")
		return false
	}

	podName := podList.Items[0].Name
	fmt.Println("Pod name:", podName)

	envMap, err := getPodEnvVariables(clientset, podName, namespace)
	if err != nil {
		fmt.Println("Error fetching environment variables from the pod:", err)
		return false
	}
	fmt.Println("Pod environment variables:", envMap)

	fileData, err := ioutil.ReadFile(jsonPath)
	if err != nil {
		fmt.Println("Error reading JSON file:", err)
		return false
	}

	var jsonData map[string]string
	if err := json.Unmarshal(fileData, &jsonData); err != nil {
		fmt.Println("Error parsing JSON file:", err)
		return false
	}
	fmt.Println("JSON data:", jsonData)

	for key, value := range jsonData {
		if val, ok := envMap[key]; ok {
			if strings.ReplaceAll(val, " ", "") != strings.ReplaceAll(value, " ", "") {
				fmt.Printf("Mismatch: Key '%s' values do not match. Pod value: %s, JSON value: %s\n", key, val, value)
				return false
			} else {
				fmt.Printf("Match: Key '%s' values match. Pod value: %s, JSON value: %s\n", key, val, value)
			}
		} else {
			fmt.Printf("Key '%s' not found in pod environment variables\n", key)
			return false
		}
	}
	return true
}

func getPodEnvVariables(clientset *kubernetes.Clientset, podName, namespace string) (map[string]string, error) {
	pod, err := clientset.CoreV1().Pods("default").Get(context.TODO(), podName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	envMap := make(map[string]string)

	for _, container := range pod.Spec.Containers {
		for _, envVar := range container.Env {
			envMap[envVar.Name] = envVar.Value
		}
	}

	return envMap, nil
}
