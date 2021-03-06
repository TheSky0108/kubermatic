/*
Copyright 2020 The Kubermatic Kubernetes Platform contributors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package etcd

import (
	"errors"
	"fmt"
	"strconv"

	kubermaticv1 "github.com/kubermatic/kubermatic/pkg/crd/kubermatic/v1"
	"github.com/kubermatic/kubermatic/pkg/resources"
	"github.com/kubermatic/kubermatic/pkg/resources/reconciling"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	name = "etcd"
	// ImageTag defines the image tag to use for the etcd image
	etcdImageTagV33 = "v3.3.18"
	etcdImageTagV34 = "v3.4.3"
)

var (
	baseTags = map[string]string{
		etcdImageTagV33: "v33",
		etcdImageTagV34: "v34",
	}

	defaultResourceRequirements = map[string]*corev1.ResourceRequirements{
		name: {
			Requests: corev1.ResourceList{
				corev1.ResourceMemory: resource.MustParse("256Mi"),
				corev1.ResourceCPU:    resource.MustParse("50m"),
			},
			Limits: corev1.ResourceList{
				corev1.ResourceMemory: resource.MustParse("2Gi"),
				corev1.ResourceCPU:    resource.MustParse("2"),
			},
		},
	}
)

type etcdStatefulSetCreatorData interface {
	Cluster() *kubermaticv1.Cluster
	GetPodTemplateLabels(string, []corev1.Volume, map[string]string) (map[string]string, error)
	ImageRegistry(string) string
	EtcdDiskSize() resource.Quantity
	GetClusterRef() metav1.OwnerReference
	SupportsFailureDomainZoneAntiAffinity() bool
}

// StatefulSetCreator returns the function to reconcile the etcd StatefulSet
func StatefulSetCreator(data etcdStatefulSetCreatorData, enableDataCorruptionChecks bool) reconciling.NamedStatefulSetCreatorGetter {
	return func() (string, reconciling.StatefulSetCreator) {
		return resources.EtcdStatefulSetName, func(set *appsv1.StatefulSet) (*appsv1.StatefulSet, error) {
			set.Name = resources.EtcdStatefulSetName

			set.Spec.Replicas = resources.Int32(resources.EtcdClusterSize)
			set.Spec.UpdateStrategy.Type = appsv1.RollingUpdateStatefulSetStrategyType
			set.Spec.PodManagementPolicy = appsv1.ParallelPodManagement
			set.Spec.ServiceName = resources.EtcdServiceName
			set.Spec.Template.Spec.ImagePullSecrets = []corev1.LocalObjectReference{{Name: resources.ImagePullSecretName}}

			baseLabels := getBasePodLabels(data.Cluster())
			set.Spec.Selector = &metav1.LabelSelector{
				MatchLabels: baseLabels,
			}

			volumes := getVolumes()
			podLabels, err := data.GetPodTemplateLabels(resources.EtcdStatefulSetName, volumes, baseLabels)
			if err != nil {
				return nil, fmt.Errorf("failed to create pod labels: %v", err)
			}

			set.Spec.Template.ObjectMeta = metav1.ObjectMeta{
				Name:   name,
				Labels: podLabels,
			}
			image, err := getLauncherImage(data)
			if err != nil {
				return nil, err
			}
			set.Spec.Template.Spec.Containers = []corev1.Container{
				{
					Name: resources.EtcdStatefulSetName,

					Image:           image,
					ImagePullPolicy: corev1.PullIfNotPresent,
					Command:         []string{"/usr/local/bin/etcd-launcher"},
					Env: []corev1.EnvVar{
						{
							Name: "POD_NAME",
							ValueFrom: &corev1.EnvVarSource{
								FieldRef: &corev1.ObjectFieldSelector{
									APIVersion: "v1",
									FieldPath:  "metadata.name",
								},
							},
						},
						{
							Name: "POD_IP",
							ValueFrom: &corev1.EnvVarSource{
								FieldRef: &corev1.ObjectFieldSelector{
									APIVersion: "v1",
									FieldPath:  "status.podIP",
								},
							},
						},
						{
							Name: "NAMESPACE",
							ValueFrom: &corev1.EnvVarSource{
								FieldRef: &corev1.ObjectFieldSelector{
									APIVersion: "v1",
									FieldPath:  "metadata.namespace",
								},
							},
						},
						{
							Name:  "TOKEN",
							Value: data.Cluster().Name,
						},
						{
							Name:  "ECTD_CLUSTER_SIZE",
							Value: strconv.Itoa(resources.EtcdClusterSize),
						},
						{
							Name:  "ENABLE_CORRUPTION_CHECK",
							Value: strconv.FormatBool(enableDataCorruptionChecks),
						},
						{
							Name:  "ETCDCTL_API",
							Value: "3",
						},
						{
							Name:  "ETCDCTL_CACERT",
							Value: "/etc/etcd/pki/ca/ca.crt",
						},
						{
							Name:  "ETCDCTL_CERT",
							Value: "/etc/etcd/pki/client/apiserver-etcd-client.crt",
						},
						{
							Name:  "ETCDCTL_KEY",
							Value: "/etc/etcd/pki/client/apiserver-etcd-client.key",
						},
						{
							Name:  "ETCDCTL_ENDPOINTS",
							Value: "https://127.0.0.1:2379",
						},
					},
					Ports: []corev1.ContainerPort{
						{
							ContainerPort: 2379,
							Protocol:      corev1.ProtocolTCP,
							Name:          "client",
						},
						{
							ContainerPort: 2380,
							Protocol:      corev1.ProtocolTCP,
							Name:          "peer",
						},
					},
					ReadinessProbe: &corev1.Probe{
						TimeoutSeconds:      10,
						PeriodSeconds:       30,
						SuccessThreshold:    1,
						FailureThreshold:    3,
						InitialDelaySeconds: 15,
						Handler: corev1.Handler{
							Exec: &corev1.ExecAction{
								Command: []string{
									"/usr/local/bin/etcdctl",
									"--command-timeout", "10s",
									"endpoint", "health",
								},
							},
						},
					},
					VolumeMounts: []corev1.VolumeMount{
						{
							Name:      "data",
							MountPath: "/var/run/etcd",
						},
						{
							Name:      resources.EtcdTLSCertificateSecretName,
							MountPath: "/etc/etcd/pki/tls",
						},
						{
							Name:      resources.CASecretName,
							MountPath: "/etc/etcd/pki/ca",
						},
						{
							Name:      resources.ApiserverEtcdClientCertificateSecretName,
							MountPath: "/etc/etcd/pki/client",
							ReadOnly:  true,
						},
					},
				},
			}
			err = resources.SetResourceRequirements(set.Spec.Template.Spec.Containers, defaultResourceRequirements, resources.GetOverrides(data.Cluster().Spec.ComponentsOverride), set.Annotations)
			if err != nil {
				return nil, fmt.Errorf("failed to set resource requirements: %v", err)
			}

			set.Spec.Template.Spec.Affinity = resources.HostnameAntiAffinity(resources.EtcdStatefulSetName, data.Cluster().Name)
			if data.SupportsFailureDomainZoneAntiAffinity() {
				antiAffinities := set.Spec.Template.Spec.Affinity.PodAntiAffinity.PreferredDuringSchedulingIgnoredDuringExecution
				antiAffinities = append(antiAffinities, resources.FailureDomainZoneAntiAffinity(resources.EtcdStatefulSetName, data.Cluster().Name))
				set.Spec.Template.Spec.Affinity.PodAntiAffinity.PreferredDuringSchedulingIgnoredDuringExecution = antiAffinities
			}

			set.Spec.Template.Spec.Volumes = volumes

			// Make sure, we don't change size of existing pvc's
			// Phase needs to be taken from an existing
			diskSize := data.EtcdDiskSize()
			if len(set.Spec.VolumeClaimTemplates) == 0 {
				set.Spec.VolumeClaimTemplates = []corev1.PersistentVolumeClaim{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:            "data",
							OwnerReferences: []metav1.OwnerReference{data.GetClusterRef()},
						},
						Spec: corev1.PersistentVolumeClaimSpec{
							StorageClassName: resources.String("kubermatic-fast"),
							AccessModes:      []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{corev1.ResourceStorage: diskSize},
							},
						},
					},
				}
			}

			return set, nil
		}
	}
}

func getVolumes() []corev1.Volume {
	return []corev1.Volume{
		{
			Name: resources.EtcdTLSCertificateSecretName,
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: resources.EtcdTLSCertificateSecretName,
				},
			},
		},
		{
			Name: resources.CASecretName,
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: resources.CASecretName,
					Items: []corev1.KeyToPath{
						{
							Path: resources.CACertSecretKey,
							Key:  resources.CACertSecretKey,
						},
					},
				},
			},
		},
		{
			Name: resources.ApiserverEtcdClientCertificateSecretName,
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: resources.ApiserverEtcdClientCertificateSecretName,
				},
			},
		},
	}
}

func getBasePodLabels(cluster *kubermaticv1.Cluster) map[string]string {
	additionalLabels := map[string]string{
		"cluster": cluster.Name,
	}
	return resources.BaseAppLabels(resources.EtcdStatefulSetName, additionalLabels)
}

// ImageTag returns the correct etcd image tag for a given Cluster
// TODO: Other functions use this function, swtich them to getLauncherImage
func ImageTag(c *kubermaticv1.Cluster) string {
	if c.IsOpenshift() || c.Spec.Version.Minor() < 17 {
		return etcdImageTagV33
	}
	return etcdImageTagV34
}

func getLauncherImage(data etcdStatefulSetCreatorData) (string, error) {
	etcdTag := ImageTag(data.Cluster())
	baseTag, ok := baseTags[etcdTag]
	if !ok {
		return "", errors.New("unknown etcd tag")
	}
	return data.ImageRegistry(resources.RegistryQuay) + "/kubermatic/etcd-launcher-" + baseTag + ":" + resources.KUBERMATICCOMMIT, nil
}
