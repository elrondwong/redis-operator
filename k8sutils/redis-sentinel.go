package k8sutils

import (
	"context"
	"encoding/json"
	"errors"

	commonapi "github.com/elrondwong/redis-operator/api"
	redisv1beta2 "github.com/elrondwong/redis-operator/api/v1beta2"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes"
)

// RedisSentinelSTS is a interface to call Redis Statefulset function
type RedisSentinelSTS struct {
	RedisStateFulType             string
	ExternalConfig                *string
	Affinity                      *corev1.Affinity `json:"affinity,omitempty"`
	TerminationGracePeriodSeconds *int64           `json:"terminationGracePeriodSeconds,omitempty" protobuf:"varint,4,opt,name=terminationGracePeriodSeconds"`
	ReadinessProbe                *commonapi.Probe
	LivenessProbe                 *commonapi.Probe
}

// RedisSentinelService is a interface to call Redis Service function
type RedisSentinelService struct {
	RedisServiceRole string
}

type RedisReplicationObject struct {
	RedisReplication *redisv1beta2.RedisReplication
}

type RedisReplicationInfo struct {
	MasterIP     string
	MasterSecret commonapi.ExistingPasswordSecret
}

// Redis Sentinel Create the Redis Sentinel Setup
func CreateRedisSentinel(ctx context.Context, client kubernetes.Interface, logger logr.Logger, cr *redisv1beta2.RedisSentinel) error {
	prop := RedisSentinelSTS{
		RedisStateFulType:             "sentinel",
		Affinity:                      cr.Spec.Affinity,
		ReadinessProbe:                &cr.Spec.ReadinessProbe.Probe,
		LivenessProbe:                 &cr.Spec.LivenessProbe.Probe,
		TerminationGracePeriodSeconds: cr.Spec.TerminationGracePeriodSeconds,
	}

	if cr.Spec.RedisSentinelConfig != nil && cr.Spec.RedisSentinelConfig.AdditionalSentinelConfig != nil {
		prop.ExternalConfig = cr.Spec.RedisSentinelConfig.AdditionalSentinelConfig
	}

	return prop.CreateRedisSentinelSetup(ctx, client, logger, cr)

}

// Create RedisSentinel Service
func CreateRedisSentinelService(cr *redisv1beta2.RedisSentinel) error {

	prop := RedisSentinelService{
		RedisServiceRole: "sentinel",
	}
	return prop.CreateRedisSentinelService(cr)
}

// Create Redis Sentinel Cluster Setup
func (service RedisSentinelSTS) CreateRedisSentinelSetup(ctx context.Context, client kubernetes.Interface, logger logr.Logger, cr *redisv1beta2.RedisSentinel) error {
	stateFulName := cr.ObjectMeta.Name + "-" + service.RedisStateFulType
	labels := getRedisLabels(stateFulName, sentinel, service.RedisStateFulType, cr.ObjectMeta.Labels)
	annotations := generateStatefulSetsAnots(cr.ObjectMeta, cr.Spec.KubernetesConfig.IgnoreAnnotations)
	objectMetaInfo := generateObjectMetaInformation(stateFulName, cr.Namespace, labels, annotations)
	err := CreateOrUpdateStateFul(
		cr.Namespace,
		objectMetaInfo,
		generateRedisSentinelParams(cr, service.getSentinelCount(cr), service.ExternalConfig, service.Affinity),
		redisSentinelAsOwner(cr),
		generateRedisSentinelInitContainerParams(cr),
		generateRedisSentinelContainerParams(ctx, client, logger, cr, service.ReadinessProbe, service.LivenessProbe),
		cr.Spec.Sidecars,
	)

	if err != nil {
		logger.Error(err, "Cannot create Sentinel statefulset for Redis")
		return err
	}
	return nil
}

// Create Redis Sentile Params for the statefulset
func generateRedisSentinelParams(cr *redisv1beta2.RedisSentinel, replicas int32, externalConfig *string, affinity *corev1.Affinity) statefulSetParameters {

	res := statefulSetParameters{
		Replicas:                      &replicas,
		ClusterMode:                   false,
		NodeConfVolume:                false,
		NodeSelector:                  cr.Spec.NodeSelector,
		PodSecurityContext:            cr.Spec.PodSecurityContext,
		PriorityClassName:             cr.Spec.PriorityClassName,
		Affinity:                      affinity,
		TerminationGracePeriodSeconds: cr.Spec.TerminationGracePeriodSeconds,
		Tolerations:                   cr.Spec.Tolerations,
		ServiceAccountName:            cr.Spec.ServiceAccountName,
		UpdateStrategy:                cr.Spec.KubernetesConfig.UpdateStrategy,
		IgnoreAnnotations:             cr.Spec.KubernetesConfig.IgnoreAnnotations,
	}

	if cr.Spec.KubernetesConfig.ImagePullSecrets != nil {
		res.ImagePullSecrets = cr.Spec.KubernetesConfig.ImagePullSecrets
	}
	if externalConfig != nil {
		res.ExternalConfig = externalConfig
	}
	if cr.Spec.RedisExporter != nil {
		res.EnableMetrics = cr.Spec.RedisExporter.Enabled
	}
	if _, found := cr.ObjectMeta.GetAnnotations()[AnnotationKeyRecreateStatefulset]; found {
		res.RecreateStatefulSet = true
	}
	return res
}

// generateRedisSentinelInitContainerParams generates Redis sentinel initcontainer information
func generateRedisSentinelInitContainerParams(cr *redisv1beta2.RedisSentinel) initContainerParameters {

	initcontainerProp := initContainerParameters{}

	if cr.Spec.InitContainer != nil {
		initContainer := cr.Spec.InitContainer

		initcontainerProp = initContainerParameters{
			Enabled:               initContainer.Enabled,
			Role:                  "sentinel",
			Image:                 initContainer.Image,
			ImagePullPolicy:       initContainer.ImagePullPolicy,
			Resources:             initContainer.Resources,
			AdditionalEnvVariable: initContainer.EnvVars,
			Command:               initContainer.Command,
			Arguments:             initContainer.Args,
		}

	}

	return initcontainerProp
}

// Create Redis Sentinel Statefulset Container Params
func generateRedisSentinelContainerParams(ctx context.Context, client kubernetes.Interface, logger logr.Logger, cr *redisv1beta2.RedisSentinel, readinessProbeDef *commonapi.Probe, livenessProbeDef *commonapi.Probe) containerParameters {
	trueProperty := true
	falseProperty := false
	containerProp := containerParameters{
		Role:                  "sentinel",
		Image:                 cr.Spec.KubernetesConfig.Image,
		ImagePullPolicy:       cr.Spec.KubernetesConfig.ImagePullPolicy,
		Resources:             cr.Spec.KubernetesConfig.Resources,
		SecurityContext:       cr.Spec.SecurityContext,
		AdditionalEnvVariable: getSentinelEnvVariable(ctx, client, logger, cr),
	}
	if cr.Spec.EnvVars != nil {
		containerProp.EnvVars = cr.Spec.EnvVars
	}
	if cr.Spec.KubernetesConfig.ExistingPasswordSecret != nil {
		containerProp.EnabledPassword = &trueProperty
		containerProp.SecretName = cr.Spec.KubernetesConfig.ExistingPasswordSecret.Name
		containerProp.SecretKey = cr.Spec.KubernetesConfig.ExistingPasswordSecret.Key
	} else {
		containerProp.EnabledPassword = &falseProperty
	}
	if cr.Spec.RedisExporter != nil {
		containerProp.RedisExporterImage = cr.Spec.RedisExporter.Image
		containerProp.RedisExporterImagePullPolicy = cr.Spec.RedisExporter.ImagePullPolicy

		if cr.Spec.RedisExporter.Resources != nil {
			containerProp.RedisExporterResources = cr.Spec.RedisExporter.Resources
		}

		if cr.Spec.RedisExporter.EnvVars != nil {
			containerProp.RedisExporterEnv = cr.Spec.RedisExporter.EnvVars
		}
	}
	if readinessProbeDef != nil {
		containerProp.ReadinessProbe = readinessProbeDef
	}
	if livenessProbeDef != nil {
		containerProp.LivenessProbe = livenessProbeDef
	}
	if cr.Spec.TLS != nil {
		containerProp.TLSConfig = cr.Spec.TLS
	}

	return containerProp

}

// Get the Count of the Sentinel
func (service RedisSentinelSTS) getSentinelCount(cr *redisv1beta2.RedisSentinel) int32 {
	return cr.Spec.GetSentinelCounts(service.RedisStateFulType)
}

// Create the Service for redis sentinel
func (service RedisSentinelService) CreateRedisSentinelService(cr *redisv1beta2.RedisSentinel) error {
	serviceName := cr.ObjectMeta.Name + "-" + service.RedisServiceRole
	logger := serviceLogger(cr.Namespace, serviceName)
	labels := getRedisLabels(serviceName, sentinel, service.RedisServiceRole, cr.ObjectMeta.Labels)
	annotations := generateServiceAnots(cr.ObjectMeta, nil)

	if cr.Spec.RedisExporter != nil && cr.Spec.RedisExporter.Enabled {
		enableMetrics = true
	} else {
		enableMetrics = false
	}

	additionalServiceAnnotations := map[string]string{}
	if cr.Spec.KubernetesConfig.Service != nil {
		additionalServiceAnnotations = cr.Spec.KubernetesConfig.Service.ServiceAnnotations
	}

	objectMetaInfo := generateObjectMetaInformation(serviceName, cr.Namespace, labels, annotations)
	headlessObjectMetaInfo := generateObjectMetaInformation(serviceName+"-headless", cr.Namespace, labels, annotations)
	additionalObjectMetaInfo := generateObjectMetaInformation(serviceName+"-additional", cr.Namespace, labels, generateServiceAnots(cr.ObjectMeta, additionalServiceAnnotations))

	err := CreateOrUpdateService(cr.Namespace, headlessObjectMetaInfo, redisSentinelAsOwner(cr), false, true, "ClusterIP")
	if err != nil {
		logger.Error(err, "Cannot create headless service for Redis", "Setup.Type", service.RedisServiceRole)
		return err
	}
	err = CreateOrUpdateService(cr.Namespace, objectMetaInfo, redisSentinelAsOwner(cr), enableMetrics, false, "ClusterIP")
	if err != nil {
		logger.Error(err, "Cannot create service for Redis", "Setup.Type", service.RedisServiceRole)
		return err
	}

	additionalServiceType := "ClusterIP"
	if cr.Spec.KubernetesConfig.Service != nil {
		additionalServiceType = cr.Spec.KubernetesConfig.Service.ServiceType
	}
	err = CreateOrUpdateService(cr.Namespace, additionalObjectMetaInfo, redisSentinelAsOwner(cr), false, false, additionalServiceType)
	if err != nil {
		logger.Error(err, "Cannot create additional service for Redis", "Setup.Type", service.RedisServiceRole)
		return err
	}
	return nil

}

func getSentinelEnvVariable(ctx context.Context, client kubernetes.Interface, logger logr.Logger, cr *redisv1beta2.RedisSentinel) *[]corev1.EnvVar {
	if cr.Spec.RedisSentinelConfig == nil {
		return &[]corev1.EnvVar{}
	}

	redisReplicationInfo, err := getRedisReplicationInfo(ctx, client, logger, cr)
	if err != nil {
		logger.Error(err, "Cannot get redis replication info", "Setup.Type", cr.Name)
		return nil
	}
	var ip, secretName, secretKey string
	ip = redisReplicationInfo.MasterIP
	if redisReplicationInfo.MasterSecret.Name != nil {
		secretName = *redisReplicationInfo.MasterSecret.Name
	}
	if redisReplicationInfo.MasterSecret.Key != nil {
		secretKey = *redisReplicationInfo.MasterSecret.Key
	}

	envVar := &[]corev1.EnvVar{
		{
			Name:  "MASTER_GROUP_NAME",
			Value: cr.Spec.RedisSentinelConfig.MasterGroupName,
		},
		{
			Name:  "IP",
			Value: ip,
		},
		{
			Name:  "PORT",
			Value: cr.Spec.RedisSentinelConfig.RedisPort,
		},
		{
			Name:  "QUORUM",
			Value: cr.Spec.RedisSentinelConfig.Quorum,
		},
		{
			Name:  "DOWN_AFTER_MILLISECONDS",
			Value: cr.Spec.RedisSentinelConfig.DownAfterMilliseconds,
		},
		{
			Name:  "PARALLEL_SYNCS",
			Value: cr.Spec.RedisSentinelConfig.ParallelSyncs,
		},
		{
			Name:  "FAILOVER_TIMEOUT",
			Value: cr.Spec.RedisSentinelConfig.FailoverTimeout,
		},
	}
	if redisReplicationInfo.MasterSecret.Name != nil && redisReplicationInfo.MasterSecret.Key != nil {
		*envVar = append(*envVar, corev1.EnvVar{
			Name: "MASTER_PASSWORD",
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{Name: secretName},
					Key:                  secretKey,
				},
			},
		})
	}

	return envVar

}

func getRedisReplicationInfo(ctx context.Context, client kubernetes.Interface, logger logr.Logger, cr *redisv1beta2.RedisSentinel) (*RedisReplicationInfo, error) {
	replicationInfo := &RedisReplicationInfo{}
	dClient, err := GenerateK8sDynamicClient(GenerateK8sConfig)
	if err != nil {
		logger.Error(err, "Failed to generate dynamic client")
		return nil, err
	}

	replicationName := cr.Spec.RedisSentinelConfig.RedisReplicationName
	replicationNamespace := cr.Namespace

	var replicationInstance redisv1beta2.RedisReplication
	var realMasterPod string

	// Get Request on Dynamic Client
	customObject, err := dClient.Resource(schema.GroupVersionResource{
		Group:    "redis.redis.opstreelabs.in",
		Version:  "v1beta2",
		Resource: "redisreplications",
	}).Namespace(replicationNamespace).Get(context.TODO(), replicationName, v1.GetOptions{})

	if err != nil {
		logger.Error(err, "Failed to Execute Get Request", "replication name", replicationName, "namespace", replicationNamespace)
		return nil, err
	} else {
		logger.V(1).Info("Successfully Execute the Get Request", "replication name", replicationName, "namespace", replicationNamespace)
	}

	// Marshal CustomObject to JSON
	replicationJSON, err := customObject.MarshalJSON()
	if err != nil {
		logger.Error(err, "Failed To Load JSON")
		return nil, err
	}

	// Unmarshal The JSON on Object
	if err := json.Unmarshal(replicationJSON, &replicationInstance); err != nil {
		logger.Error(err, "Failed To Unmarshal JSON over the Object")
		return nil, err
	}

	masterPods := GetRedisNodesByRole(ctx, client, logger, &replicationInstance, "master")

	if len(masterPods) == 0 {
		realMasterPod = ""
		err := errors.New("no master pods found")
		logger.Error(err, "")
	} else if len(masterPods) == 1 {
		realMasterPod = masterPods[0]
	} else {
		realMasterPod = checkAttachedSlave(ctx, client, logger, &replicationInstance, masterPods)
	}

	realMasterInfo := RedisDetails{
		PodName:   realMasterPod,
		Namespace: replicationNamespace,
	}

	replicationInfo.MasterIP = getRedisServerIP(client, logger, realMasterInfo)
	if replicationInstance.Spec.KubernetesConfig.ExistingPasswordSecret != nil {
		replicationInfo.MasterSecret = commonapi.ExistingPasswordSecret{
			Name: replicationInstance.Spec.KubernetesConfig.ExistingPasswordSecret.Name,
			Key:  replicationInstance.Spec.KubernetesConfig.ExistingPasswordSecret.Key,
		}
	}
	return replicationInfo, nil
}
