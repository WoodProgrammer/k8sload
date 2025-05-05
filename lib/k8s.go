package main

import (
	"context"

	"github.com/rs/zerolog/log"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
)

type KubernetesClient interface {
	ApplyManifest(manifestDetails string) error
}

type KubernetesHandler struct {
	KubernetesClient *dynamic.DynamicClient
	Config           *rest.Config
}

func (k *KubernetesHandler) ApplyManifest(manifestDetails string) error {
	manifestBytes := []byte(manifestDetails)

	decUnstructured := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
	obj := &unstructured.Unstructured{}
	_, gvk, err := decUnstructured.Decode(manifestBytes, nil, obj)
	if err != nil {
		return err
	}

	mapping, err := restMapper(k.Config).RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return err
	}

	var dr dynamic.ResourceInterface
	if mapping.Scope.Name() == "namespace" {
		dr = k.KubernetesClient.Resource(mapping.Resource).Namespace(obj.GetNamespace())
	} else {
		dr = k.KubernetesClient.Resource(mapping.Resource)
	}
	_, err = dr.Create(context.TODO(), obj, metav1.CreateOptions{})
	if err != nil {
		_, err = dr.Update(context.TODO(), obj, metav1.UpdateOptions{})
		return err
	}
	if err != nil {
		return err
	}

	log.Info().Msg("Manifest applied successfully.")
	return nil
}

func restMapper(config *rest.Config) meta.RESTMapper {
	dc, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		panic(err)
	}
	gr, err := restmapper.GetAPIGroupResources(dc)
	if err != nil {
		panic(err)
	}
	return restmapper.NewDiscoveryRESTMapper(gr)
}
