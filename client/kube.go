package main

import (
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// Load client kubeconfig
func LoadKubeConfig(kubeconfig *string) (*rest.Config, error) {
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		fmt.Printf("Error %s, when getting Loading Local Kubeconfig\n", err.Error())
	}
	return config, err
}

// Load client kubeconfig
func LoadInClusterConfig() (*rest.Config, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		fmt.Printf("Error %s, when getting InClusterConfig\n", err.Error())
	}
	if cliConf.Verbose {
		fmt.Println(*config)
	}
	return config, err
}

// Init kube client
func InitKubeConfig(config *rest.Config) (*kubernetes.Clientset, error) {
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Printf("Error %s, when initializing kube client", err.Error())
	}
	return clientset, err
}

// Get pods from a given namespace
func GetPods(client kubernetes.Clientset, namespace string) (*corev1.PodList, error) {
	pods, err := client.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		fmt.Printf("Error %s, when fetching pods\n", err.Error())
		return nil, err
	}
	return pods, nil
}

// Get deployments from a given namespace
func GetDeployments(client kubernetes.Clientset, namespace string) (*appsv1.DeploymentList, error) {
	deployments, err := client.AppsV1().Deployments(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		fmt.Printf("Error %s, when fetching deployments\n", err.Error())
		return nil, err
	}
	return deployments, nil
}
