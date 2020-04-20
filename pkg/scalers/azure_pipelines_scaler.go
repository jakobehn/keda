package scalers

import (
	"context"
	"fmt"

	// "github.com/microsoft/azure-devops-go-api/azuredevops"
	// "github.com/microsoft/azure-devops-go-api/azuredevops/core"

	v2beta1 "k8s.io/api/autoscaling/v2beta1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/metrics/pkg/apis/external_metrics"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

const (
//	queueLengthMetricName    = "queueLength"
	defaultMaxJobsPerAgent = 1
	defaultTargetPipelinesQueueLength = 0
)

type azurePipelinesScaler struct {
	metadata    *azurePipelinesMetadata
}

type azurePipelinesMetadata struct {
	url				 string
	pat	        	 string
	queueName		 string
}

var azurePipelinesLog = logf.Log.WithName("azure_pipelines_scaler")


func (s *azurePipelinesScaler) GetMetrics(ctx context.Context, metricName string, metricSelector labels.Selector) ([]external_metrics.ExternalMetricValue, error) {

	queuelen, err := s.GetAzurePipelinesQueueLength(ctx)

	if err != nil {
		azurePipelinesLog.Error(err, "error getting pipelines queue length")
		return []external_metrics.ExternalMetricValue{}, err
	}

	metric := external_metrics.ExternalMetricValue{
		MetricName: metricName,
		Value:      *resource.NewQuantity(int64(queuelen), resource.DecimalSI),
		Timestamp:  metav1.Now(),
	}

	return append([]external_metrics.ExternalMetricValue{}, metric), nil
}

func (s *azurePipelinesScaler) GetAzurePipelinesQueueLength(ctx context.Context) (int32, error) {
	var err error
	return 3,err
}


// Returns the metric spec to be used by the HPA
func (s *azurePipelinesScaler) GetMetricSpecForScaling() []v2beta1.MetricSpec {
	maxJobsPerAgentQty := resource.NewQuantity(int64(defaultMaxJobsPerAgent), resource.DecimalSI)
	targetPipelinesQueueLengthQty := resource.NewQuantity(int64(defaultTargetPipelinesQueueLength), resource.DecimalSI)
	externalMetric := &v2beta1.ExternalMetricSource{MetricName: queueLengthMetricName, TargetAverageValue: maxJobsPerAgentQty, TargetValue: targetPipelinesQueueLengthQty}
	metricSpec := v2beta1.MetricSpec{External: externalMetric, Type: externalMetricType}
	return []v2beta1.MetricSpec{metricSpec}
}

func NewAzurePipelinesScaler(resolvedEnv, metadata, authParams map[string]string) (Scaler, error) {
	meta, err := parseAzurePipelinesMetadata(metadata, resolvedEnv, authParams)
	if err != nil {
		return nil, fmt.Errorf("error parsing azure Pipelines metadata: %s", err)
	}

	return &azurePipelinesScaler{
		metadata:    meta,
	}, nil
}

func parseAzurePipelinesMetadata(metadata, resolvedEnv, authParams map[string]string) (*azurePipelinesMetadata, error) {
	meta := azurePipelinesMetadata{}

	// if val, ok := metadata[queueLengthMetricName]; ok {
	// 	queueLength, err := strconv.Atoi(val)
	// 	if err != nil {
	// 		azurePipelinesLog.Error(err, "Error parsing azure Pipelines metadata", "queueLengthMetricName", queueLengthMetricName)
	// 		return nil, "", fmt.Errorf("Error parsing azure Pipelines metadata %s: %s", queueLengthMetricName, err.Error())
	// 	}

	// 	meta.targetQueueLength = queueLength
	// }

	if val, ok := metadata["queueName"]; ok && val != "" {
		meta.queueName = val
	} else {
		return nil, fmt.Errorf("no queueName given")
	}

	if val, ok := metadata["pat"]; ok && val != "" {
		meta.pat = val
	} else {
		return nil, fmt.Errorf("no pat given")
	}

	if val, ok := metadata["url"]; ok && val != "" {
		meta.url = val
	} else {
		return nil, fmt.Errorf("no url given")
	}	

	return &meta, nil
}


// GetScaleDecision is a func
func (s *azurePipelinesScaler) IsActive(ctx context.Context) (bool, error) {

	return true, nil
}

func (s *azurePipelinesScaler) Close() error {
	return nil
}
