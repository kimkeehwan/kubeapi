package k8sclient

import (
	"log"

	metricsv "k8s.io/metrics/pkg/client/clientset/versioned"
)

type MetricImpl struct {
	clients *metricsv.Clientset
}

func NewK8sMetricClient(config *K8sClusterConfig) *MetricImpl {
	client, err := metricsv.NewForConfig(config.config)
	if err != nil {
		log.Fatalf("load fail metrics client %s", err.Error())
	}
	return &MetricImpl{clients: client}
}
