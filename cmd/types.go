package main

import "k8s.io/client-go/kubernetes"

type APIServer struct {
	ListenAddr string
	k8Client   *kubernetes.Clientset
}
