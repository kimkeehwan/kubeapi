package k8sclient

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type K8sClusterConfig struct {
	config *rest.Config
}

// type ClientImpl struct {
// 	clients *kubernetes.Clientset
// }

// type DynamicImpl struct {
// 	clients dynamic.Interface
// }

// type MetricImpl struct {
// 	clients *metricsv.Clientset
// }

func getKubeConfig(kubeconfigPath string) (*rest.Config, error) {
	var kubeconfig *rest.Config

	if kubeconfigPath != "" {
		config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
		if err != nil {
			return nil, fmt.Errorf("unable to load kubeconfig from %s: %v", kubeconfigPath, err)
		}
		kubeconfig = config
	} else {
		config, err := rest.InClusterConfig()
		if err != nil {
			return nil, fmt.Errorf("unable to load in-cluster config: %v", err)
		}
		kubeconfig = config
	}
	return kubeconfig, nil
}

func NewClusterConfig() *K8sClusterConfig {

	var kubeconfigfile *string

	if home := homedir.HomeDir(); home != "" {
		kubeconfigfile = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
		if _, err := os.Stat(*kubeconfigfile); err != nil {
			*kubeconfigfile = ""
		}
	}
	kubeconfig, err := getKubeConfig(*kubeconfigfile)
	if err != nil {
		log.Fatalf("load fail kube config %s", err.Error())
	}

	return &K8sClusterConfig{config: kubeconfig}
}

// func NewK8sClient(config *K8sClusterConfig) *ClientImpl {

// 	client, err := kubernetes.NewForConfig(config.config)
// 	if err != nil {
// 		log.Fatalf("load fail kubeclient %s", err.Error())
// 	}

// 	return &ClientImpl{clients: client}
// }

// func NewK8sDynamicClient(config *K8sClusterConfig) *DynamicImpl {

// 	client, err := dynamic.NewForConfig(config.config)
// 	if err != nil {
// 		log.Fatalf("load fail dynamic client %s", err.Error())
// 	}
// 	return &DynamicImpl{clients: client}
// }

// func NewK8sMetricClient(config *K8sClusterConfig) *MetricImpl {
// 	client, err := metricsv.NewForConfig(config.config)
// 	if err != nil {
// 		log.Fatalf("load fail metrics client %s", err.Error())
// 	}
// 	return &MetricImpl{clients: client}
// }
