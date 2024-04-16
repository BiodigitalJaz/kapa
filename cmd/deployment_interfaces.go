package main

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type DeploymentServices interface {
	CreateDeployment(namespace string, deployment *appsv1.Deployment) error
	GetDeployment(namespace, name string) (*appsv1.Deployment, error)
	UpdateDeployment(namespace string, deployment *appsv1.Deployment) error
	DeleteDeployment(namespace, name string) error
}

// KubernetesDeploymentService implements the DeploymentService interface using Kubernetes client-go.
type KubernetesDeploymentService struct {
	Clientset *kubernetes.Clientset
}

func (kds *KubernetesDeploymentService) CreateDeployment(namespace string, deployment *appsv1.Deployment) error {
	_, err := kds.Clientset.AppsV1().Deployments(namespace).Create(context.Background(), deployment, metav1.CreateOptions{})
	return err
}

func (kds *KubernetesDeploymentService) GetDeployment(namespace, name string) (*appsv1.Deployment, error) {
	deployment, err := kds.Clientset.AppsV1().Deployments(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return deployment, nil
}

func (kds *KubernetesDeploymentService) UpdateDeployment(namespace string, deployment *appsv1.Deployment) error {
	_, err := kds.Clientset.AppsV1().Deployments(namespace).Update(context.Background(), deployment, metav1.UpdateOptions{})
	return err
}

func (kds *KubernetesDeploymentService) DeleteDeployment(namespace, name string) error {
	err := kds.Clientset.AppsV1().Deployments(namespace).Delete(context.Background(), name, metav1.DeleteOptions{})
	return err
}

// NewKubernetesDeploymentService creates a new instance of KubernetesDeploymentService.
func NewKubernetesDeploymentService(clientset *kubernetes.Clientset) *KubernetesDeploymentService {
	return &KubernetesDeploymentService{
		Clientset: clientset,
	}
}
