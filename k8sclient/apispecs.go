package k8sclient

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

type ApiSpecs map[string]metav1.APIResource

func (s ApiSpecs) GetResource(kind string) metav1.APIResource {
	return s[kind]
}
