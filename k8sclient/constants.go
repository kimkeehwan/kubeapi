package k8sclient

const (
	LABEL_ENVIRONMENT         = LabelName("environment")
	LABEL_COMPONENT           = LabelName("component")
	LABEL_APP_IDPP2           = LabelName("app.idpp2")
	LABEL_PART_OF             = LabelName("part-of")
	LABEL_IDPP2_BUILDER       = LabelName("aiblab.co.kr/builder")
	LABEL_IDPP2_BUILDER_TAGET = LabelName("aiblab.co.kr/builder-target")
	LABEL_IDPP2_USERNAME      = LabelName("aiblab.co.kr/username")
	LABEL_JUPYTERUSERNAME     = LabelName("hub.jupyter.org/username")
	LABEL_NVIDIA_GPU          = LabelName(NVIDIA_GPU_NAME)
)

const (
	FIELD_MANAGER                  = "projectmanager"
	AIBLAB_BUILDER                 = "idpp2-builder"
	AIBLAB_BUILDER_TAGET_IPYKERNEL = "ipykernel"
	AIBLAB_BUILDER_TAGET_IRKERNEL  = "irkernel"
	AIBLAB_BUILDER_TAGET_TFSERVING = "tfserving"
	NVIDIA_GPU_NAME                = "nvidia.com/gpu"
)

const (
	TYPEMETA_KIND_NAMESPACE  = "Namespace"
	TYPEMETA_KIND_POD        = "Pod"
	TYPEMETA_KIND_SERVICE    = "Service"
	TYPEMETA_KIND_DEPLOYMENT = "Deployment"
	TYPEMETA_KIND_JOB        = "Job"
)

const (
	TYPEMETA_APIVERSION_V1              = "v1"
	TYPEMETA_APIVERSION_BATCH_V1        = "batch/v1"
	TYPEMETA_APIVERSION_METRICS_V1BETA1 = "metrics.k8s.io/v1beta1"
)
