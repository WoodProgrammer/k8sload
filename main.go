package main

import (
	"fmt"
	"os"

	lib "github.com/WoodProgrammer/k8sload/lib"
	"github.com/rs/zerolog/log"
)

func main() {
	fmt.Println("k8load v0.0.1")
	output, err := lib.GenerateManifestFile("load.yaml", "lib/_base_template.tmpl")

	if err != nil {
		log.Err(err).Msg("Error while running lib.GenerateManifestFile()")
		os.Exit(1)
	}
	fmt.Println(output)
}
