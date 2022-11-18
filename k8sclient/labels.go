package k8sclient

type LabelName string

type K8sLabels map[string]string

func (s K8sLabels) Put(key string, val string) K8sLabels {
	s[key] = val

	return s
}

func (s K8sLabels) Get(key string) string {
	return s[key]
}
