package main

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type PodServices interface {
	GetPods(namespace string) ([]string, error)
	GetPodLogs(podName, namespace string) ([]string, error)
}

type KubernetesPodService struct {
	Clientset *kubernetes.Clientset
}

func (kps *KubernetesPodService) GetPods(namespace string) ([]string, error) {
	pods, err := kps.Clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var nsPods []string
	for _, pod := range pods.Items {
		nsPods = append(nsPods, pod.Name)
	}

	return nsPods, nil
}

func (kps *KubernetesPodService) GetPodLogs(podName, namespace string) ([]string, error) {
	// Create a pod logs request
	req := kps.Clientset.CoreV1().Pods(namespace).GetLogs(podName, &corev1.PodLogOptions{})

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

func NewKubernetesPodService(clientset *kubernetes.Clientset) *KubernetesPodService {
	return &KubernetesPodService{
		Clientset: clientset,
	}
}
