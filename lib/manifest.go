package lib

import (
	"bytes"
	"os"
	"text/template"

	"github.com/rs/zerolog/log"
	yaml "gopkg.in/yaml.v3"
)

func GenerateManifestFile(manifestFile, baseTemplateFile string) (string, error) {
	var cfg Config
	data, err := os.ReadFile(manifestFile)
	if err != nil {
		log.Err(err).Msgf("There is an error while reading data in this file %s", manifestFile)
		return "", err
	}

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		log.Err(err).Msg("There is an error while Unmarshalling Yaml structure")
		return "", err
	}

	tmpl, err := template.ParseFiles(baseTemplateFile)
	if err != nil {
		log.Err(err).Msgf("There is an error while parsing template file, please check base template file with that name %s", baseTemplateFile)
		return "", err

	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, cfg); err != nil {
		log.Err(err).Msg("There is an error while execute template output")
		return "", err
	}
	return buf.String(), nil
}
