package k8sclient

import (
	"context"
	"log"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	//metricsapi "k8s.io/metrics/pkg/apis/metrics"
	metricsV1beta1api "k8s.io/metrics/pkg/apis/metrics/v1beta1"
	metricsclientset "k8s.io/metrics/pkg/client/clientset/versioned"
)

type MetricsImpl struct {
	clients *metricsclientset.Clientset
}

func NewK8sMetricClient(config *K8sClusterConfig) *MetricsImpl {
	client, err := metricsclientset.NewForConfig(config.config)
	if err != nil {
		log.Fatalf("load fail metrics client %s", err.Error())
	}
	return &MetricsImpl{clients: client}
}

func (s *MetricsImpl) ListNode(ctx context.Context, selector string) (*metricsV1beta1api.NodeMetricsList, error) {
	opts := metav1.ListOptions{}
	if selector != "" {
		opts.LabelSelector = selector
	}
	return s.clients.MetricsV1beta1().NodeMetricses().List(ctx, opts)
}

func (s *MetricsImpl) GetNode(ctx context.Context, name string) (*metricsV1beta1api.NodeMetrics, error) {
	opt := metav1.GetOptions{}

	return s.clients.MetricsV1beta1().NodeMetricses().Get(ctx, name, opt)
}

func (s *MetricsImpl) ListPod(ctx context.Context, namespace string, selector string) (*metricsV1beta1api.PodMetricsList, error) {
	opts := metav1.ListOptions{}
	if selector != "" {
		opts.LabelSelector = selector
	}
	return s.clients.MetricsV1beta1().PodMetricses(namespace).List(ctx, opts)
}

func (s *MetricsImpl) GetPod(ctx context.Context, namespace string, name string) (*metricsV1beta1api.PodMetrics, error) {
	opt := metav1.GetOptions{}

	return s.clients.MetricsV1beta1().PodMetricses(namespace).Get(ctx, name, opt)
}

//func (s *MetricImpl)
