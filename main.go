package main

import (
	"fmt"
	"os"

	lib "github.com/WoodProgrammer/k8sload/lib"
	"github.com/rs/zerolog/log"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
)

func NewKubernetesClient() lib.KubernetesClient {
	kubeconfig := os.Getenv("KUBECONFIG")
	if kubeconfig == "" {
		kubeconfig = os.Getenv("HOME") + "/.kube/config"
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err)
	}

	dynClient, err := dynamic.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	return &lib.KubernetesHandler{
		KubernetesClient: dynClient,
		Config:           config,
	}
}

func main() {
	log.Info().Msg("k8load v0.0.1")

	//gk8sClient := NewKubernetesClient()
	output, err := lib.GenerateManifestFile("load.yaml", "lib/_base_template.tmpl")

	if err != nil {
		log.Err(err).Msg("Error while running lib.GenerateManifestFile()")
		os.Exit(1)
	}
	fmt.Println(output)
	/*err = k8sClient.ApplyManifest(output)
	if err != nil {
		log.Err(err).Msg("Error while running k8sClient.ApplyManifest()")
		os.Exit(1)
	}*/
}
