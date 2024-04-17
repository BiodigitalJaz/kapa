package main

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type NamespaceServices interface {
	GetNamespaces() ([]string, error)
	GetNamespaceEvents(namespace string) ([]string, error)
}

type KubernetesNamespaceService struct {
	Clientset *kubernetes.Clientset
}

func NewKubernetesNamespaceService(clientset *kubernetes.Clientset) *KubernetesNamespaceService {
	return &KubernetesNamespaceService{
		Clientset: clientset,
	}
}

func (kns *KubernetesNamespaceService) GetNamespaces() ([]string, error) {
	namespace, err := kns.Clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var namespaceNames []string
	for _, ns := range namespace.Items {
		namespaceNames = append(namespaceNames, ns.Name)
	}

	return namespaceNames, nil
}

func (kns *KubernetesNamespaceService) GetNamespaceEvents(namespace string) ([]string, error) {
	events, err := kns.Clientset.CoreV1().Events(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var nsEvents []string
	for _, event := range events.Items {
		nsEvents = append(nsEvents, event.Message)
	}

	return nsEvents, nil
}
