package metricsserver

import (
	"fmt"

	"github.com/kubermatic/kubermatic/api/pkg/resources/vpnsidecar"

	"github.com/kubermatic/kubermatic/api/pkg/resources"
	"github.com/kubermatic/kubermatic/api/pkg/resources/apiserver"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

const (
	name = "metrics-server"

	tag = "v0.3.1"
)

// Deployment returns the metrics-server Deployment
func Deployment(data resources.DeploymentDataProvider, existing *appsv1.Deployment) (*appsv1.Deployment, error) {
	var dep *appsv1.Deployment
	if existing != nil {
		dep = existing
	} else {
		dep = &appsv1.Deployment{}
	}

	dep.Name = resources.MetricsServerDeploymentName
	dep.OwnerReferences = []metav1.OwnerReference{data.GetClusterRef()}

	dep.Labels = resources.BaseAppLabel(name, nil)

	dep.Spec.Replicas = resources.Int32(1)
	dep.Spec.Selector = &metav1.LabelSelector{
		MatchLabels: resources.BaseAppLabel(name, nil),
	}
	dep.Spec.Strategy.Type = appsv1.RollingUpdateStatefulSetStrategyType
	dep.Spec.Strategy.RollingUpdate = &appsv1.RollingUpdateDeployment{
		MaxSurge: &intstr.IntOrString{
			Type:   intstr.Int,
			IntVal: 1,
		},
		MaxUnavailable: &intstr.IntOrString{
			Type:   intstr.Int,
			IntVal: 0,
		},
	}

	volumes := getVolumes()
	podLabels, err := data.GetPodTemplateLabels(name, volumes, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create pod labels: %v", err)
	}

	dep.Spec.Template.ObjectMeta = metav1.ObjectMeta{
		Labels: podLabels,
	}

	dep.Spec.Template.Spec.Volumes = volumes

	apiserverIsRunningContainer, err := apiserver.IsRunningInitContainer(data)
	if err != nil {
		return nil, err
	}
	dep.Spec.Template.Spec.InitContainers = []corev1.Container{*apiserverIsRunningContainer}

	openvpnSidecar, err := vpnsidecar.OpenVPNSidecarContainer(data, "openvpn-client")
	if err != nil {
		return nil, fmt.Errorf("failed to get openvpn-client sidecar: %v", err)
	}

	dnatControllerSidecar, err := vpnsidecar.DnatControllerContainer(data, "dnat-controller")
	if err != nil {
		return nil, fmt.Errorf("failed to get dnat-controller sidecar: %v", err)
	}

	dep.Spec.Template.Spec.Containers = []corev1.Container{
		{
			Name:            name,
			Image:           data.ImageRegistry(resources.RegistryGCR) + "/google_containers/metrics-server-amd64:" + tag,
			ImagePullPolicy: corev1.PullIfNotPresent,
			Command:         []string{"/metrics-server"},
			Args: []string{
				"--kubeconfig", "/etc/kubernetes/kubeconfig/kubeconfig",
				"--authentication-kubeconfig", "/etc/kubernetes/kubeconfig/kubeconfig",
				"--authorization-kubeconfig", "/etc/kubernetes/kubeconfig/kubeconfig",
				"--kubelet-port", "10250",
				"--kubelet-insecure-tls",
				// We use the same as the API server as we use the same dnat-controller
				"--kubelet-preferred-address-types", "ExternalIP,InternalIP",
				"--v", "1",
				"--logtostderr",
			},
			TerminationMessagePath:   corev1.TerminationMessagePathDefault,
			TerminationMessagePolicy: corev1.TerminationMessageReadFile,
			VolumeMounts: []corev1.VolumeMount{
				{
					Name:      resources.MetricsServerKubeconfigSecretName,
					MountPath: "/etc/kubernetes/kubeconfig",
					ReadOnly:  true,
				},
			},
		},
		*openvpnSidecar,
		*dnatControllerSidecar,
	}

	return dep, nil
}

func getVolumes() []corev1.Volume {
	return []corev1.Volume{
		{
			Name: resources.MetricsServerKubeconfigSecretName,
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: resources.MetricsServerKubeconfigSecretName,
					// We have to make the secret readable for all for now because owner/group cannot be changed.
					// ( upstream proposal: https://github.com/kubernetes/kubernetes/pull/28733 )
					DefaultMode: resources.Int32(resources.DefaultAllReadOnlyMode),
				},
			},
		},
		{
			Name: resources.OpenVPNClientCertificatesSecretName,
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName:  resources.OpenVPNClientCertificatesSecretName,
					DefaultMode: resources.Int32(resources.DefaultOwnerReadOnlyMode),
				},
			},
		},
		{
			Name: resources.KubeletDnatControllerKubeconfigSecretName,
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName:  resources.KubeletDnatControllerKubeconfigSecretName,
					DefaultMode: resources.Int32(resources.DefaultOwnerReadOnlyMode),
				},
			},
		},
	}
}