package main

import (
	"fmt"
	"os"
	"time"

	"k8s.io/client-go/kubernetes"
)

func main() {
	cliConfig := getCLIArgs(&cliConf)
	appConfig := getAppConfig(cliConf.AppConfig)

	if cliConf.Verbose {
		// Print the configurations for debugging
		fmt.Printf("CLI Config: %+v\n", cliConfig)
		fmt.Printf("App Config: %+v\n", appConfig)
	}

	config, err := LoadKubeConfig(&cliConfig.KubeConfig)
	if err != nil {
		fmt.Println("Trying to parse Inclutser config...")
		config, err = LoadInClusterConfig()
		if err != nil {
			fmt.Println("I am about to fail. Bye, bye")
		}
	}

	client, err := InitKubeConfig(config)
	if err != nil {
		os.Exit(1)
	}
	testLoop(client, cliConfig.Namespace)
}

func testLoop(client *kubernetes.Clientset, ns string) {
	loop := 0

	for {
		loop += 1
		fmt.Printf("Loop %d\n", loop)

		pods, err := GetPods(*client, ns)
		if err != nil {
			fmt.Printf("Failed to get pods: %v\n", err)
		} else {
			for _, pod := range pods.Items {
				fmt.Printf("Pod %s (%s)\n", pod.Name, pod.Namespace)
			}
		}

		deployments, err := GetDeployments(*client, cliConf.Namespace)
		if err != nil {
			fmt.Printf("Failed to get deployments: %v\n", err)
		} else {
			for _, deploy := range deployments.Items {
				fmt.Printf("Deployment %s (%s)\n", deploy.Name, deploy.Namespace)
			}
		}
		time.Sleep(1 * time.Second)
	}
}
