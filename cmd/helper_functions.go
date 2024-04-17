package main

import (
	"bytes"
	"io"
	"log"
	"path/filepath"

	v1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func setupK8Access() (*kubernetes.Clientset, error) {

	kubeconfigPath := filepath.Join(homedir.HomeDir(), ".kube", "config")

	config, err := rest.InClusterConfig()
	if err != nil {
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfigPath)
		if err != nil {
			log.Printf("error loading kubeconfig: %v\n", err)
		}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Printf("error: %v", err)
	}

	return clientset, nil
}

func convertYAMLToTypedDeployment(deploymentYAML []byte) (v1.Deployment, error) {
	// Parse YAML into unstructured object
	var deploymentObj unstructured.Unstructured
	decoder := yaml.NewYAMLOrJSONDecoder(io.NopCloser(bytes.NewReader(deploymentYAML)), 4096)
	if err := decoder.Decode(&deploymentObj); err != nil {
		return v1.Deployment{}, err
	}

	// Convert unstructured object to typed Deployment
	var typedDeployment v1.Deployment
	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(deploymentObj.Object, &typedDeployment); err != nil {
		return v1.Deployment{}, err
	}
	return typedDeployment, nil
}
