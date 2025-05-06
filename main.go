package main

import (
	"os"
	"sync"

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
	var manifestArr []string
	log.Info().Msg("k8load v0.0.1")
	maxConcurrency := 2

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, maxConcurrency)

	k8sClient := NewKubernetesClient()
	ProducerManifest, err := lib.GenerateManifestFile("load.yaml", "lib/_base_template_producer.tmpl")
	manifestArr = append(manifestArr, ProducerManifest)
	ProducerSvcManifest, err := lib.GenerateManifestFile("load.yaml", "lib/_base_template_producer_svc.tmpl")
	manifestArr = append(manifestArr, ProducerSvcManifest)

	ConsumerManifest, err := lib.GenerateManifestFile("load.yaml", "lib/_base_template_consumer.tmpl")
	manifestArr = append(manifestArr, ConsumerManifest)
	ConsumerSvcManifest, err := lib.GenerateManifestFile("load.yaml", "lib/_base_template_consumer_svc.tmpl")
	manifestArr = append(manifestArr, ConsumerSvcManifest)

	for _, r := range manifestArr {
		wg.Add(1)
		semaphore <- struct{}{} // acquire a token

		go func(manifest string) {
			defer wg.Done()
			defer func() { <-semaphore }() // release token

			if err != nil {
				log.Err(err).Msg("Error while running lib.GenerateManifestFile() ")
				os.Exit(1)
			}

			err = k8sClient.ApplyManifest(r)
			if err != nil {
				log.Err(err).Msg("Error while running k8sClient.ApplyManifest()")
				os.Exit(1)
			}
		}(r)
	}
	wg.Wait()
}
