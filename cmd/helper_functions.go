package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"path/filepath"

	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
func (a *APIServer) getPods(namespace string) ([]string, error) {
	pods, err := a.k8Client.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list pods: %v", err)
	}

	var nsPods []string
	for _, pod := range pods.Items {
		nsPods = append(nsPods, pod.Name)
	}

	return nsPods, nil
}

func (a *APIServer) getNamespaces() ([]string, error) {
	namespace, err := a.k8Client.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list namespaces: %v", err)
	}

	var namespaceNames []string
	for _, ns := range namespace.Items {
		namespaceNames = append(namespaceNames, ns.Name)
	}

	return namespaceNames, nil
}

func (a *APIServer) getNamespaceEvents(namespace string) ([]string, error) {
	events, err := a.k8Client.CoreV1().Events(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get events: %v", err)
	}

	var nsEvents []string
	for _, event := range events.Items {
		nsEvents = append(nsEvents, event.Message)
	}

	return nsEvents, nil
}

func (a *APIServer) getPodLogs(podName, namespace string) ([]string, error) {
	// Create a pod logs request
	req := a.k8Client.CoreV1().Pods(namespace).GetLogs(podName, &corev1.PodLogOptions{})

	// Retrieve logs
	logs, err := req.Stream(context.Background())
	if err != nil {
		return nil, err
	}
	defer logs.Close()

	// Read logs from the stream
	buf := make([]byte, 4096)
	var logOutput []string
	for {
		bytesRead, err := logs.Read(buf)
		if err != nil {
			break
		}
		logOutput = append(logOutput, string(buf[:bytesRead]))
	}

	return logOutput, nil
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
