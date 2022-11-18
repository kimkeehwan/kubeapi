package k8sclient

import (
	"context"
	"log"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type DynamicImpl struct {
	clients dynamic.Interface
}

func NewK8sDynamicClient(config *K8sClusterConfig) *DynamicImpl {

	client, err := dynamic.NewForConfig(config.config)
	if err != nil {
		log.Fatalf("load fail dynamic client %s", err.Error())
	}
	return &DynamicImpl{clients: client}
}

func (s *DynamicImpl) Apply(ctx context.Context, namespace string, gvr schema.GroupVersionResource, resource ResourceSpecs) (*unstructured.Unstructured, error) {
	data := &unstructured.Unstructured{Object: resource}

	return s.clients.Resource(gvr).Namespace(namespace).Apply(ctx, resource.GetName(), data, metav1.ApplyOptions{FieldManager: FIELD_MANAGER})
}
