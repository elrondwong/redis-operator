package k8sutils

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/util/retry"
)

func podLogger(namespace string, name string) logr.Logger {
	reqLogger := log.WithValues("Request.Pod.Namespace", namespace, "Request.Pod.Name", name)
	return reqLogger
}

// GetPodsByLabels retrieves a list of Pods that match the specified labels.
func GetPodsByLabels(client kubernetes.Interface, labels map[string]string, namespace string) (*v1.PodList, error) {
	logger := podLogger(namespace, "podList")
	// Build the label selector for Pod list
	labelSelector := metav1.LabelSelector{MatchLabels: labels}
	podListOptions := metav1.ListOptions{LabelSelector: metav1.FormatLabelSelector(&labelSelector)}

	// Call the Kubernetes client to get the list of Pods that match the label selector
	pods, err := client.CoreV1().Pods(namespace).List(context.TODO(), podListOptions)
	if err != nil {
		logger.Error(err, "Failed in getting podList for redis")
		return nil, err

	}

	return pods, nil
}

// AddLabelToPod adds a new label to the specified Pod.
func AddLabelToPod(client kubernetes.Interface, namespace, podName, key, value string) error {
	logger := podLogger(namespace, podName)
	return retry.RetryOnConflict(retry.DefaultRetry, func() error {
		// Fetch the latest Pod object
		currentPod, err := client.CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
		if err != nil {
			logger.Error(err, "Failed to get pod")
			return err
		}

		// Add or update the label
		labels := currentPod.Labels
		if labels == nil {
			labels = make(map[string]string)
		}
		labels[key] = value
		patchData := fmt.Sprintf(`{"metadata":{"labels":%s}}`, toJSON(labels))

		// Update the Pod with the new label
		_, err = client.CoreV1().Pods(namespace).Patch(context.TODO(), podName, types.StrategicMergePatchType,
			[]byte(patchData), metav1.PatchOptions{})
		return err
	})
}

func toJSON(data map[string]string) string {
	result := "{"
	for key, value := range data {
		result += fmt.Sprintf(`"%s":"%s",`, key, value)
	}
	if len(data) > 0 {
		result = result[:len(result)-1]
	}
	result += "}"
	return result
}
