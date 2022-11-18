package k8sclient

import (
	"context"
	"fmt"
	"log"
	"time"

	"k8s.io/client-go/kubernetes"

	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	eventv1 "k8s.io/api/events/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	configv1 "k8s.io/client-go/applyconfigurations/core/v1"
	configmetav1 "k8s.io/client-go/applyconfigurations/meta/v1"

	storagev1 "k8s.io/api/storage/v1"
)

type ClientImpl struct {
	clients  *kubernetes.Clientset
	apiSpecs ApiSpecs
}

func NewK8sClient(config *K8sClusterConfig) *ClientImpl {

	client, err := kubernetes.NewForConfig(config.config)
	if err != nil {
		log.Fatalf("load fail kubeclient %s", err.Error())
	}

	return &ClientImpl{clients: client}
}

func (s *ClientImpl) ApiSpecs() ApiSpecs {
	if s.apiSpecs != nil {
		return s.apiSpecs
	}
	r := make(ApiSpecs)

	preferredList, _ := s.clients.DiscoveryClient.ServerPreferredResources()
	for _, preferred := range preferredList {
		for _, res := range preferred.APIResources {
			r[res.Kind] = res
		}
	}

	s.apiSpecs = r
	return r
}

func (s *ClientImpl) ObjectToNamespace(object runtime.Object) *corev1.Namespace {
	return object.(*corev1.Namespace)
}

func (s *ClientImpl) ObjectToResourceQuota(object runtime.Object) *corev1.ResourceQuota {
	return object.(*corev1.ResourceQuota)
}

func (s *ClientImpl) ObjectToPod(object runtime.Object) *corev1.Pod {
	return object.(*corev1.Pod)
}

func (s *ClientImpl) ObjectToJob(object runtime.Object) *batchv1.Job {
	return object.(*batchv1.Job)
}

func (s *ClientImpl) ApplyNamespace(ctx context.Context, name string, labels K8sLabels) (*corev1.Namespace, error) {
	kind := TYPEMETA_KIND_NAMESPACE
	ver := TYPEMETA_APIVERSION_V1

	config := configv1.NamespaceApplyConfiguration{
		TypeMetaApplyConfiguration:   configmetav1.TypeMetaApplyConfiguration{Kind: &kind, APIVersion: &ver},
		ObjectMetaApplyConfiguration: &configmetav1.ObjectMetaApplyConfiguration{Name: &name, Labels: labels},
	}
	opt := metav1.ApplyOptions{FieldManager: FIELD_MANAGER}

	return s.clients.CoreV1().Namespaces().Apply(ctx, &config, opt)
}

func (s *ClientImpl) ListNamespace(ctx context.Context, selector string) (*corev1.NamespaceList, error) {
	opts := metav1.ListOptions{}
	if selector != "" {
		opts.LabelSelector = selector
	}
	return s.clients.CoreV1().Namespaces().List(ctx, opts)
}

func (s *ClientImpl) GetNamespace(ctx context.Context, name string) (*corev1.Namespace, error) {
	opt := metav1.GetOptions{}
	return s.clients.CoreV1().Namespaces().Get(ctx, name, opt)
}

func (s *ClientImpl) DeleteNamespace(ctx context.Context, name string) error {
	// GracePeriodSeconds
	// Preconditions
	// OrphanDependents
	// PropagationPolicy
	//metav1.NewDeleteOptions(grace)
	//metav1.NewPreconditionDeleteOptions(uid)
	//metav1.NewRVDeletionPrecondition(rv)
	//metav1.NewUIDPreconditions()
	opt := metav1.DeleteOptions{} // == metav1.NewDeleteOptions(0)
	return s.clients.CoreV1().Namespaces().Delete(ctx, name, opt)
}

func (s *ClientImpl) WatchNamespaces(ctx context.Context, selector string) (watch.Interface, error) {
	opt := metav1.ListOptions{}
	if selector != "" {
		opt.LabelSelector = selector
	}

	return s.clients.CoreV1().Namespaces().Watch(ctx, opt)
}

func (s *ClientImpl) ListPod(ctx context.Context, namespace string, selector string) (*corev1.PodList, error) {
	opts := metav1.ListOptions{}
	if selector != "" {
		opts.LabelSelector = selector
	}

	return s.clients.CoreV1().Pods(namespace).List(ctx, opts)
}

func (s *ClientImpl) GetPod(ctx context.Context, namespace string, name string) (*corev1.Pod, error) {
	opt := metav1.GetOptions{}
	return s.clients.CoreV1().Pods(namespace).Get(ctx, name, opt)
}

func (s *ClientImpl) DeletePod(ctx context.Context, namespace string, name string) error {
	opt := metav1.DeleteOptions{} // == metav1.NewDeleteOptions(0)
	return s.clients.CoreV1().Pods(namespace).Delete(ctx, name, opt)
}

func (s *ClientImpl) DeletePods(ctx context.Context, namespace string, selector string) error {
	deleteOpt := metav1.DeleteOptions{}
	selectOpt := metav1.ListOptions{}
	if selector != "" {
		selectOpt.LabelSelector = selector
	}
	return s.clients.CoreV1().Pods(namespace).DeleteCollection(ctx, deleteOpt, selectOpt)
}

func (s *ClientImpl) WatchPods(ctx context.Context, namespace string, selector string) (watch.Interface, error) {
	opt := metav1.ListOptions{}
	if selector != "" {
		opt.LabelSelector = selector
	}

	return s.clients.CoreV1().Pods(namespace).Watch(ctx, opt)
}

func (s *ClientImpl) ListConfigMap(ctx context.Context, namespace string, selector string) (*corev1.ConfigMapList, error) {
	opts := metav1.ListOptions{}
	if selector != "" {
		opts.LabelSelector = selector
	}

	return s.clients.CoreV1().ConfigMaps(namespace).List(ctx, opts)
}

func (s *ClientImpl) DeleteConfigMap(ctx context.Context, namespace string, name string) error {
	opt := metav1.DeleteOptions{} // == metav1.NewDeleteOptions(0)
	return s.clients.CoreV1().ConfigMaps(namespace).Delete(ctx, name, opt)
}

func (s *ClientImpl) DeleteConfigMaps(ctx context.Context, namespace string, selector string) error {
	deleteOpt := metav1.DeleteOptions{}
	selectOpt := metav1.ListOptions{}
	if selector != "" {
		selectOpt.LabelSelector = selector
	}
	return s.clients.CoreV1().ConfigMaps(namespace).DeleteCollection(ctx, deleteOpt, selectOpt)
}

func (s *ClientImpl) GetConfigMap(ctx context.Context, namespace string, name string) (*corev1.ConfigMap, error) {
	opt := metav1.GetOptions{}
	return s.clients.CoreV1().ConfigMaps(namespace).Get(ctx, name, opt)
}

func (s *ClientImpl) ListNode(ctx context.Context, selector string) (*corev1.NodeList, error) {
	opts := metav1.ListOptions{}
	if selector != "" {
		opts.LabelSelector = selector
	}

	return s.clients.CoreV1().Nodes().List(ctx, opts)
}

func (s *ClientImpl) GetNode(ctx context.Context, name string) (*corev1.Node, error) {
	opt := metav1.GetOptions{}
	return s.clients.CoreV1().Nodes().Get(ctx, name, opt)
}

func (s *ClientImpl) ListEndpoints(ctx context.Context, namespace string, selector string) (*corev1.EndpointsList, error) {
	opts := metav1.ListOptions{}
	if selector != "" {
		opts.LabelSelector = selector
	}

	return s.clients.CoreV1().Endpoints(namespace).List(ctx, opts)
}

func (s *ClientImpl) GetEndpoints(ctx context.Context, namspace string, name string) (*corev1.Endpoints, error) {
	opt := metav1.GetOptions{}
	return s.clients.CoreV1().Endpoints(namspace).Get(ctx, name, opt)
}

func (s *ClientImpl) ListEvent(ctx context.Context, namespace string, selector string) (*eventv1.EventList, error) {
	opts := metav1.ListOptions{}
	if selector != "" {
		opts.LabelSelector = selector
	}

	return s.clients.EventsV1().Events(namespace).List(ctx, opts)
}

func (s *ClientImpl) GetEvent(ctx context.Context, namspace string, name string) (*eventv1.Event, error) {
	opt := metav1.GetOptions{}

	return s.clients.EventsV1().Events(namspace).Get(ctx, name, opt)
}

func (s *ClientImpl) ListLimitRange(ctx context.Context, namespace string, selector string) (*corev1.LimitRangeList, error) {
	opts := metav1.ListOptions{}
	if selector != "" {
		opts.LabelSelector = selector
	}

	return s.clients.CoreV1().LimitRanges(namespace).List(ctx, opts)
}

func (s *ClientImpl) GetLimitRange(ctx context.Context, namspace string, name string) (*corev1.LimitRange, error) {
	opt := metav1.GetOptions{}
	return s.clients.CoreV1().LimitRanges(namspace).Get(ctx, name, opt)
}

func (s *ClientImpl) DeleteLimitRange(ctx context.Context, namespace string, name string) error {
	opt := metav1.DeleteOptions{} // == metav1.NewDeleteOptions(0)
	return s.clients.CoreV1().LimitRanges(namespace).Delete(ctx, name, opt)
}

func (s *ClientImpl) DeleteLimitRanges(ctx context.Context, namespace string, selector string) error {
	deleteOpt := metav1.DeleteOptions{}
	selectOpt := metav1.ListOptions{}
	if selector != "" {
		selectOpt.LabelSelector = selector
	}
	return s.clients.CoreV1().LimitRanges(namespace).DeleteCollection(ctx, deleteOpt, selectOpt)
}

func (s *ClientImpl) ListPersistentVolumeClaim(ctx context.Context, namespace string, selector string) (*corev1.PersistentVolumeClaimList, error) {
	opts := metav1.ListOptions{}
	if selector != "" {
		opts.LabelSelector = selector
	}

	return s.clients.CoreV1().PersistentVolumeClaims(namespace).List(ctx, opts)
}

func (s *ClientImpl) DeletePersistentVolumeClaim(ctx context.Context, namespace string, name string) error {
	opt := metav1.DeleteOptions{} // == metav1.NewDeleteOptions(0)
	return s.clients.CoreV1().PersistentVolumeClaims(namespace).Delete(ctx, name, opt)
}

func (s *ClientImpl) DeletePersistentVolumeClaims(ctx context.Context, namespace string, selector string) error {
	deleteOpt := metav1.DeleteOptions{}
	selectOpt := metav1.ListOptions{}
	if selector != "" {
		selectOpt.LabelSelector = selector
	}
	return s.clients.CoreV1().PersistentVolumeClaims(namespace).DeleteCollection(ctx, deleteOpt, selectOpt)
}

func (s *ClientImpl) GetPersistentVolumeClaim(ctx context.Context, namspace string, name string) (*corev1.PersistentVolumeClaim, error) {
	opt := metav1.GetOptions{}
	return s.clients.CoreV1().PersistentVolumeClaims(namspace).Get(ctx, name, opt)
}

func (s *ClientImpl) ListPersistentVolume(ctx context.Context, selector string) (*corev1.PersistentVolumeList, error) {
	opts := metav1.ListOptions{}
	if selector != "" {
		opts.LabelSelector = selector
	}

	return s.clients.CoreV1().PersistentVolumes().List(ctx, opts)
}

func (s *ClientImpl) GetPersistentVolume(ctx context.Context, name string) (*corev1.PersistentVolume, error) {
	opt := metav1.GetOptions{}
	return s.clients.CoreV1().PersistentVolumes().Get(ctx, name, opt)
}

func (s *ClientImpl) ListPodTemplate(ctx context.Context, namespace string, selector string) (*corev1.PodTemplateList, error) {
	opts := metav1.ListOptions{}
	if selector != "" {
		opts.LabelSelector = selector
	}

	return s.clients.CoreV1().PodTemplates(namespace).List(ctx, opts)
}

func (s *ClientImpl) DeletePodTemplate(ctx context.Context, namespace string, name string) error {
	opt := metav1.DeleteOptions{} // == metav1.NewDeleteOptions(0)
	return s.clients.CoreV1().PodTemplates(namespace).Delete(ctx, name, opt)
}

func (s *ClientImpl) DeletePodTemplates(ctx context.Context, namespace string, selector string) error {
	deleteOpt := metav1.DeleteOptions{}
	selectOpt := metav1.ListOptions{}
	if selector != "" {
		selectOpt.LabelSelector = selector
	}
	return s.clients.CoreV1().PodTemplates(namespace).DeleteCollection(ctx, deleteOpt, selectOpt)
}

func (s *ClientImpl) GetPodTemplate(ctx context.Context, namespace string, name string) (*corev1.PodTemplate, error) {
	opt := metav1.GetOptions{}
	return s.clients.CoreV1().PodTemplates(namespace).Get(ctx, name, opt)
}

func (s *ClientImpl) ListSecret(ctx context.Context, namespace string, selector string) (*corev1.SecretList, error) {
	opts := metav1.ListOptions{}
	if selector != "" {
		opts.LabelSelector = selector
	}

	return s.clients.CoreV1().Secrets(namespace).List(ctx, opts)
}

func (s *ClientImpl) GetSecret(ctx context.Context, namespace string, name string) (*corev1.Secret, error) {
	opt := metav1.GetOptions{}
	return s.clients.CoreV1().Secrets(namespace).Get(ctx, name, opt)
}

func (s *ClientImpl) DeleteSecret(ctx context.Context, namespace string, name string) error {
	opt := metav1.DeleteOptions{} // == metav1.NewDeleteOptions(0)
	return s.clients.CoreV1().Secrets(namespace).Delete(ctx, name, opt)
}

func (s *ClientImpl) DeleteSecrets(ctx context.Context, namespace string, selector string) error {
	deleteOpt := metav1.DeleteOptions{}
	selectOpt := metav1.ListOptions{}
	if selector != "" {
		selectOpt.LabelSelector = selector
	}
	return s.clients.CoreV1().Secrets(namespace).DeleteCollection(ctx, deleteOpt, selectOpt)
}

func (s *ClientImpl) ListReplicationController(ctx context.Context, namespace string, selector string) (*corev1.ReplicationControllerList, error) {
	opts := metav1.ListOptions{}
	if selector != "" {
		opts.LabelSelector = selector
	}

	return s.clients.CoreV1().ReplicationControllers(namespace).List(ctx, opts)
}

func (s *ClientImpl) GetReplicationController(ctx context.Context, namespace string, name string) (*corev1.ReplicationController, error) {
	opt := metav1.GetOptions{}
	return s.clients.CoreV1().ReplicationControllers(namespace).Get(ctx, name, opt)
}

func (s *ClientImpl) DeleteReplicationController(ctx context.Context, namespace string, name string) error {
	opt := metav1.DeleteOptions{} // == metav1.NewDeleteOptions(0)
	return s.clients.CoreV1().ReplicationControllers(namespace).Delete(ctx, name, opt)
}

func (s *ClientImpl) DeleteReplicationControllers(ctx context.Context, namespace string, selector string) error {
	deleteOpt := metav1.DeleteOptions{}
	selectOpt := metav1.ListOptions{}
	if selector != "" {
		selectOpt.LabelSelector = selector
	}
	return s.clients.CoreV1().ReplicationControllers(namespace).DeleteCollection(ctx, deleteOpt, selectOpt)
}

func (s *ClientImpl) ListServiceAccount(ctx context.Context, namespace string, selector string) (*corev1.ServiceAccountList, error) {
	opts := metav1.ListOptions{}
	if selector != "" {
		opts.LabelSelector = selector
	}

	return s.clients.CoreV1().ServiceAccounts(namespace).List(ctx, opts)
}

func (s *ClientImpl) GetServiceAccount(ctx context.Context, namespace string, name string) (*corev1.ServiceAccount, error) {
	opt := metav1.GetOptions{}
	return s.clients.CoreV1().ServiceAccounts(namespace).Get(ctx, name, opt)
}

func (s *ClientImpl) DeleteServiceAccount(ctx context.Context, namespace string, name string) error {
	opt := metav1.DeleteOptions{} // == metav1.NewDeleteOptions(0)
	return s.clients.CoreV1().ServiceAccounts(namespace).Delete(ctx, name, opt)
}

func (s *ClientImpl) DeleteServiceAccounts(ctx context.Context, namespace string, selector string) error {
	deleteOpt := metav1.DeleteOptions{}
	selectOpt := metav1.ListOptions{}
	if selector != "" {
		selectOpt.LabelSelector = selector
	}
	return s.clients.CoreV1().ServiceAccounts(namespace).DeleteCollection(ctx, deleteOpt, selectOpt)
}

func (s *ClientImpl) ListResourceQuota(ctx context.Context, namespace string, selector string) (*corev1.ResourceQuotaList, error) {
	opts := metav1.ListOptions{}
	if selector != "" {
		opts.LabelSelector = selector
	}

	return s.clients.CoreV1().ResourceQuotas(namespace).List(ctx, opts)
}

func (s *ClientImpl) GetResourceQuota(ctx context.Context, namespace string, name string) (*corev1.ResourceQuota, error) {
	opt := metav1.GetOptions{}
	return s.clients.CoreV1().ResourceQuotas(namespace).Get(ctx, name, opt)
}

func (s *ClientImpl) DeleteResourceQuota(ctx context.Context, namespace string, name string) error {
	opt := metav1.DeleteOptions{} // == metav1.NewDeleteOptions(0)
	return s.clients.CoreV1().ResourceQuotas(namespace).Delete(ctx, name, opt)
}

func (s *ClientImpl) DeleteResourceQuotas(ctx context.Context, namespace string, selector string) error {
	deleteOpt := metav1.DeleteOptions{}
	selectOpt := metav1.ListOptions{}
	if selector != "" {
		selectOpt.LabelSelector = selector
	}
	return s.clients.CoreV1().ResourceQuotas(namespace).DeleteCollection(ctx, deleteOpt, selectOpt)
}

func (s ClientImpl) WatchResourceQuotas(ctx context.Context, namespace string, selector string) (watch.Interface, error) {
	opts := metav1.ListOptions{}
	if selector != "" {
		opts.LabelSelector = selector
	}
	return s.clients.CoreV1().ResourceQuotas(namespace).Watch(ctx, opts)
}

func (s *ClientImpl) ListService(ctx context.Context, namespace string, selector string) (*corev1.ServiceList, error) {
	opts := metav1.ListOptions{}
	if selector != "" {
		opts.LabelSelector = selector
	}
	return s.clients.CoreV1().Services(namespace).List(ctx, opts)
}

func (s *ClientImpl) GetService(ctx context.Context, namespace string, name string) (*corev1.Service, error) {
	opt := metav1.GetOptions{}
	return s.clients.CoreV1().Services(namespace).Get(ctx, name, opt)
}

func (s *ClientImpl) DeleteService(ctx context.Context, namespace string, name string) error {
	opt := metav1.DeleteOptions{} // == metav1.NewDeleteOptions(0)
	return s.clients.CoreV1().Services(namespace).Delete(ctx, name, opt)
}

func (s *ClientImpl) ListDeployment(ctx context.Context, namespace string, selector string) (*appsv1.DeploymentList, error) {
	opts := metav1.ListOptions{}
	if selector != "" {
		opts.LabelSelector = selector
	}

	return s.clients.AppsV1().Deployments(namespace).List(ctx, opts)
}

func (s *ClientImpl) GetDeployment(ctx context.Context, namespace string, name string) (*appsv1.Deployment, error) {
	opt := metav1.GetOptions{}
	return s.clients.AppsV1().Deployments(namespace).Get(ctx, name, opt)
}

func (s *ClientImpl) DeleteDeployment(ctx context.Context, namespace string, name string) error {
	opt := metav1.DeleteOptions{} // == metav1.NewDeleteOptions(0)
	return s.clients.AppsV1().Deployments(namespace).Delete(ctx, name, opt)
}

func (s *ClientImpl) DeleteDeployments(ctx context.Context, namespace string, selector string) error {
	deleteOpt := metav1.DeleteOptions{}
	selectOpt := metav1.ListOptions{}
	if selector != "" {
		selectOpt.LabelSelector = selector
	}

	return s.clients.AppsV1().Deployments(namespace).DeleteCollection(ctx, deleteOpt, selectOpt)
}

func (s *ClientImpl) RestartDeployment(ctx context.Context, namespace string, name string) (*appsv1.Deployment, error) {

	opts := metav1.PatchOptions{FieldManager: FIELD_MANAGER}

	data := fmt.Sprintf(`{"spec":{"template":{"metadata":{"annotations":{"builder.aiblab.co.kr/restartedAt":"%s"}}}}}`, time.Now().String())
	return s.clients.AppsV1().Deployments(namespace).Patch(ctx, name, types.StrategicMergePatchType, []byte(data), opts)
}

func (s *ClientImpl) ListDaemonSet(ctx context.Context, namespace string, selector string) (*appsv1.DaemonSetList, error) {
	opts := metav1.ListOptions{}
	if selector != "" {
		opts.LabelSelector = selector
	}

	return s.clients.AppsV1().DaemonSets(namespace).List(ctx, opts)
}

func (s *ClientImpl) DeleteDaemonSet(ctx context.Context, namespace string, name string) error {
	opt := metav1.DeleteOptions{} // == metav1.NewDeleteOptions(0)
	return s.clients.AppsV1().DaemonSets(namespace).Delete(ctx, name, opt)
}

func (s *ClientImpl) DeleteDaemonSets(ctx context.Context, namespace string, selector string) error {
	deleteOpt := metav1.DeleteOptions{}
	selectOpt := metav1.ListOptions{}
	if selector != "" {
		selectOpt.LabelSelector = selector
	}
	return s.clients.AppsV1().DaemonSets(namespace).DeleteCollection(ctx, deleteOpt, selectOpt)
}

func (s *ClientImpl) GetDaemonSet(ctx context.Context, namespace string, name string) (*appsv1.DaemonSet, error) {
	opt := metav1.GetOptions{}
	return s.clients.AppsV1().DaemonSets(namespace).Get(ctx, name, opt)
}

func (s *ClientImpl) RestartDaemonSet(ctx context.Context, namespace string, name string) (*appsv1.DaemonSet, error) {

	opts := metav1.PatchOptions{FieldManager: FIELD_MANAGER}

	data := fmt.Sprintf(`{"spec":{"template":{"metadata":{"annotations":{"builder.aiblab.co.kr/restartedAt":"%s"}}}}}`, time.Now().String())
	return s.clients.AppsV1().DaemonSets(namespace).Patch(ctx, name, types.StrategicMergePatchType, []byte(data), opts)
}

func (s *ClientImpl) ListStatefulSet(ctx context.Context, namespace string, selector string) (*appsv1.StatefulSetList, error) {
	opts := metav1.ListOptions{}
	if selector != "" {
		opts.LabelSelector = selector
	}

	return s.clients.AppsV1().StatefulSets(namespace).List(ctx, opts)
}

func (s *ClientImpl) GetStatefulSet(ctx context.Context, namespace string, name string) (*appsv1.StatefulSet, error) {
	opt := metav1.GetOptions{}
	return s.clients.AppsV1().StatefulSets(namespace).Get(ctx, name, opt)
}

func (s *ClientImpl) DeleteStatefulSet(ctx context.Context, namespace string, name string) error {
	opt := metav1.DeleteOptions{} // == metav1.NewDeleteOptions(0)
	return s.clients.AppsV1().StatefulSets(namespace).Delete(ctx, name, opt)
}

func (s *ClientImpl) DeleteStatefulSets(ctx context.Context, namespace string, selector string) error {
	deleteOpt := metav1.DeleteOptions{}
	selectOpt := metav1.ListOptions{}
	if selector != "" {
		selectOpt.LabelSelector = selector
	}
	return s.clients.AppsV1().StatefulSets(namespace).DeleteCollection(ctx, deleteOpt, selectOpt)
}

func (s *ClientImpl) RestartStatefulSet(ctx context.Context, namespace string, name string) (*appsv1.StatefulSet, error) {

	opts := metav1.PatchOptions{FieldManager: FIELD_MANAGER}

	data := fmt.Sprintf(`{"spec":{"template":{"metadata":{"annotations":{"builder.aiblab.co.kr/restartedAt":"%s"}}}}}`, time.Now().String())
	return s.clients.AppsV1().StatefulSets(namespace).Patch(ctx, name, types.StrategicMergePatchType, []byte(data), opts)
}

func (s *ClientImpl) ListReplicaSet(ctx context.Context, namespace string, selector string) (*appsv1.ReplicaSetList, error) {
	opts := metav1.ListOptions{}
	if selector != "" {
		opts.LabelSelector = selector
	}

	return s.clients.AppsV1().ReplicaSets(namespace).List(ctx, opts)
}

func (s *ClientImpl) GetReplicaSet(ctx context.Context, namespace string, name string) (*appsv1.ReplicaSet, error) {
	opt := metav1.GetOptions{}
	return s.clients.AppsV1().ReplicaSets(namespace).Get(ctx, name, opt)
}

func (s *ClientImpl) DeleteReplicaSet(ctx context.Context, namespace string, name string) error {
	opt := metav1.DeleteOptions{} // == metav1.NewDeleteOptions(0)
	return s.clients.AppsV1().ReplicaSets(namespace).Delete(ctx, name, opt)
}

func (s *ClientImpl) DeleteReplicaSets(ctx context.Context, namespace string, selector string) error {
	deleteOpt := metav1.DeleteOptions{}
	selectOpt := metav1.ListOptions{}
	if selector != "" {
		selectOpt.LabelSelector = selector
	}
	return s.clients.AppsV1().ReplicaSets(namespace).DeleteCollection(ctx, deleteOpt, selectOpt)
}

func (s *ClientImpl) ListJob(ctx context.Context, namespace string, selector string) (*batchv1.JobList, error) {
	opts := metav1.ListOptions{}
	if selector != "" {
		opts.LabelSelector = selector
	}
	return s.clients.BatchV1().Jobs(namespace).List(ctx, opts)
}

func (s *ClientImpl) GetJob(ctx context.Context, namespace string, name string) (*batchv1.Job, error) {
	opt := metav1.GetOptions{}
	return s.clients.BatchV1().Jobs(namespace).Get(ctx, name, opt)
}

func (s *ClientImpl) DeleteJob(ctx context.Context, namespace string, name string) error {
	opt := metav1.DeleteOptions{} // == metav1.NewDeleteOptions(0)
	return s.clients.BatchV1().Jobs(namespace).Delete(ctx, name, opt)
}

func (s *ClientImpl) DeleteJobs(ctx context.Context, namespace string, selector string) error {
	deleteOpt := metav1.DeleteOptions{}
	selectOpt := metav1.ListOptions{}
	if selector != "" {
		selectOpt.LabelSelector = selector
	}
	return s.clients.BatchV1().Jobs(namespace).DeleteCollection(ctx, deleteOpt, selectOpt)
}

func (s *ClientImpl) WatchJobs(ctx context.Context, namespace string, selector string) (watch.Interface, error) {
	opt := metav1.ListOptions{}
	if selector != "" {
		opt.LabelSelector = selector
	}
	return s.clients.BatchV1().Jobs(namespace).Watch(ctx, opt)
}

func (s *ClientImpl) ListCronJob(ctx context.Context, namespace string, selector string) (*batchv1.CronJobList, error) {
	opts := metav1.ListOptions{}
	if selector != "" {
		opts.LabelSelector = selector
	}
	return s.clients.BatchV1().CronJobs(namespace).List(ctx, opts)
}

func (s *ClientImpl) GetCronJob(ctx context.Context, namespace string, name string) (*batchv1.CronJob, error) {
	opt := metav1.GetOptions{}
	return s.clients.BatchV1().CronJobs(namespace).Get(ctx, name, opt)
}

func (s *ClientImpl) DeleteCronJob(ctx context.Context, namespace string, name string) error {
	opt := metav1.DeleteOptions{} // == metav1.NewDeleteOptions(0)
	return s.clients.BatchV1().CronJobs(namespace).Delete(ctx, name, opt)
}

func (s *ClientImpl) DeleteCronJobs(ctx context.Context, namespace string, selector string) error {
	deleteOpt := metav1.DeleteOptions{}
	selectOpt := metav1.ListOptions{}
	if selector != "" {
		selectOpt.LabelSelector = selector
	}
	return s.clients.BatchV1().CronJobs(namespace).DeleteCollection(ctx, deleteOpt, selectOpt)
}

func (s *ClientImpl) GetStorageClass(ctx context.Context, namespace string, name string) (*storagev1.StorageClass, error) {
	opt := metav1.GetOptions{}

	return s.clients.StorageV1().StorageClasses().Get(ctx, name, opt)
}
