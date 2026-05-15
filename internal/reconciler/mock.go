package reconciler

import (
	"context"
)

// MockCoreClient is a placeholder implementation for development
type MockCoreClient struct{}

func (m *MockCoreClient) DeployWorkload(ctx context.Context, spec *WorkloadSpec) (*WorkloadStatus, error) {
	return &WorkloadStatus{
		Phase:         "Running",
		ReadyReplicas: spec.Replicas,
	}, nil
}

func (m *MockCoreClient) GetWorkloadStatus(ctx context.Context, name, namespace string) (*WorkloadStatus, error) {
	return &WorkloadStatus{
		Phase:         "Running",
		ReadyReplicas: 1,
	}, nil
}

func (m *MockCoreClient) CreateNamespace(ctx context.Context, name string, labels map[string]string, quota *ResourceQuota) error {
	return nil
}

func (m *MockCoreClient) ApplyPolicy(ctx context.Context, namespace string, policy *PolicySpec) error {
	return nil
}
