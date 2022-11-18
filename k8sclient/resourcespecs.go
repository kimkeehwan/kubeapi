package k8sclient

import (
	"fmt"

	"k8s.io/apimachinery/pkg/runtime/schema"
)

type ResourceSpecs map[string]interface{}

func (s ResourceSpecs) GetKind() string {
	return fmt.Sprint(s["kind"])
}

func (s ResourceSpecs) GetName() string {
	metadata := s["metadata"].(ResourceSpecs)
	return fmt.Sprint(metadata["name"])
}

func (s ResourceSpecs) GetApiVersion() string {
	return fmt.Sprint(s["apiVersion"])
}

func (s ResourceSpecs) GroupVersionResource(client *ClientImpl) schema.GroupVersionResource {

	resource := client.ApiSpecs().GetResource(s.GetKind())

	gv, _ := schema.ParseGroupVersion(s.GetApiVersion())
	r := schema.GroupVersionResource{Group: gv.Group, Version: gv.Version, Resource: resource.Name}

	return r

}
