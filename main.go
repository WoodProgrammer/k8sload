package main

import (
	"os"
	"sync"
	"time"

	lib "github.com/WoodProgrammer/k8sload/lib"
	"github.com/common-nighthawk/go-figure"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	file string
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

func LoadK8S() {
	var manifestArr []string
	maxConcurrency := 2

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, maxConcurrency)

	k8sClient := NewKubernetesClient()
	ProducerManifest, err := lib.GenerateManifestFile(file, lib.ProducerDeploymentTemplate)
	manifestArr = append(manifestArr, ProducerManifest)
	ProducerSvcManifest, err := lib.GenerateManifestFile(file, lib.ProducerSvcTemplate)
	manifestArr = append(manifestArr, ProducerSvcManifest)

	ConsumerManifest, err := lib.GenerateManifestFile(file, lib.ConsumerDeploymentTemplate)
	manifestArr = append(manifestArr, ConsumerManifest)
	ConsumerSvcManifest, err := lib.GenerateManifestFile(file, lib.ConsumerSvcTemplate)
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
			log.Info().Msgf("Manifest applied successfully")
		}(r)
	}
	wg.Wait()
	log.Info().Msgf("Please check the details on prometheus exposed services on consumer and producer side as well.")

}
func main() {

	consoleWriter := zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}
	consoleWriter.FormatLevel = func(i interface{}) string {
		// Add colors to log level
		switch i {
		case "debug":
			return "\x1b[36mDEBUG\x1b[0m"
		case "info":
			return "\x1b[32mINFO\x1b[0m"
		case "warn":
			return "\x1b[33mWARN\x1b[0m"
		case "error":
			return "\x1b[31mERROR\x1b[0m"
		default:
			return "\x1b[37m" + i.(string) + "\x1b[0m"
		}
	}
	log.Logger = zerolog.New(consoleWriter).With().Timestamp().Logger()

	asciiTitle := figure.NewFigure("K8sLoad", "", true)
	asciiTitle.Print()

	var rootCmd = &cobra.Command{
		Use:   "k8sload",
		Short: "Kubernetes plugin to spin up test cases",
		Run: func(cmd *cobra.Command, args []string) {
			LoadK8S()
		},
	}
	rootCmd.Flags().StringVarP(&file, "file", "F", file, "The k8sload base yaml file")

	rootCmd.MarkFlagRequired("file")

	if err := rootCmd.Execute(); err != nil {
		log.Err(err).Msg("CLI execution failed")
		os.Exit(1)
	}

}
