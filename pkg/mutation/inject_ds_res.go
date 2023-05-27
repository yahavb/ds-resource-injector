package mutation

import (
        "context"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
        "k8s.io/client-go/kubernetes"
        "k8s.io/client-go/rest"
        metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
        "strconv"
        "k8s.io/api/core/v1"
        "k8s.io/apimachinery/pkg/api/resource"
)

// injectDsRes is a container for the mutation injecting environment vars
type injectDsRes struct {
	Logger logrus.FieldLogger
}

// injectDsRes implements the podMutator interface
var _ podMutator = (*injectDsRes)(nil)

// Name returns the struct name
func (se injectDsRes) Name() string {
	return "inject_env"
}

// Mutate returns a new mutated pod according to set env rules
func (se injectDsRes) Mutate(pod *corev1.Pod) (*corev1.Pod, error) {

        nodeName := pod.Spec.Affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution.NodeSelectorTerms[0].MatchFields[0].Values[0]
	se.Logger.Debugf("nodeName=%s",nodeName)
        config, err := rest.InClusterConfig()
        if err != nil {
          panic(err)
        }
        client, err := kubernetes.NewForConfig(config)
        if err != nil {
          panic(err)
        }
        node, err := client.CoreV1().Nodes().Get(context.Background(), nodeName, metav1.GetOptions{})
        if err != nil {
          panic(err)
        }
        instanceCpu := node.Labels["karpenter.k8s.aws/instance-cpu"]
        se.Logger.Debugf("Node instance-cpu: %s\n", instanceCpu)

        instanceMilliCpu, err := strconv.Atoi(instanceCpu)
        if err != nil {
          panic(err)
        }
        instanceMilliCpu = instanceMilliCpu * 1000
        se.Logger.Debugf("Node instanceMilliCpu: %d\n", instanceMilliCpu)
        podCpuResources := int(float64(instanceMilliCpu) * 0.1)
        se.Logger.Debugf("Node pod-cpu resources: %d\n", podCpuResources)

	se.Logger = se.Logger.WithField("mutation", se.Name())
	mpod := pod.DeepCopy()
        se.Logger.Debugf("pod containers Resources before: %s\n", mpod.Spec.Containers[0].Resources)
        mpod.Spec.Containers[0].Resources = v1.ResourceRequirements{
          Limits: v1.ResourceList{
            "cpu": resource.MustParse(strconv.Itoa(podCpuResources)+"m"),
          },
          Requests: v1.ResourceList{
            "cpu": resource.MustParse(strconv.Itoa(podCpuResources)+"m"),
          },
        } 
        se.Logger.Debugf("pod containers Resources after: %s\n", mpod.Spec.Containers[0].Resources)

	return mpod, nil
}

// injectDsResVar injects a var in both containers and init containers of a pod
func injectDsResVar(pod *corev1.Pod, envVar corev1.EnvVar) {
	for i, container := range pod.Spec.Containers {
		if !HasEnvVar(container, envVar) {
			pod.Spec.Containers[i].Env = append(container.Env, envVar)
		}
	}
	for i, container := range pod.Spec.InitContainers {
		if !HasEnvVar(container, envVar) {
			pod.Spec.InitContainers[i].Env = append(container.Env, envVar)
		}
	}
}

// HasEnvVar returns true if environment variable exists false otherwise
func HasEnvVar(container corev1.Container, checkEnvVar corev1.EnvVar) bool {
	for _, envVar := range container.Env {
		if envVar.Name == checkEnvVar.Name {
			return true
		}
	}
	return false
}
