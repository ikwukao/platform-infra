package controller

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/ikwukao/platform-infra/internal/deployments"
)

type mockDeploymentRepository struct {
	items []deployments.Deployment
}

func (m *mockDeploymentRepository) Create(
	_ context.Context,
	deployment *deployments.Deployment,
) error {
	m.items = append(m.items, *deployment)

	return nil
}

func (m *mockDeploymentRepository) Get(
	_ context.Context,
	id uuid.UUID,
) (*deployments.Deployment, error) {
	for _, deployment := range m.items {
		if deployment.ID == id {
			result := deployment
			return &result, nil
		}
	}

	return nil, deployments.ErrNotFound
}

func (m *mockDeploymentRepository) ListByService(
	_ context.Context,
	serviceID uuid.UUID,
) ([]deployments.Deployment, error) {
	var result []deployments.Deployment

	for _, deployment := range m.items {
		if deployment.ServiceID == serviceID {
			result = append(result, deployment)
		}
	}

	return result, nil
}

func (m *mockDeploymentRepository) UpdateStatus(
	_ context.Context,
	id uuid.UUID,
	status string,
) error {
	for i := range m.items {
		if m.items[i].ID == id {
			m.items[i].Status = status
			return nil
		}
	}

	return deployments.ErrNotFound
}

func TestReconcilePendingDeployment(t *testing.T) {
	deploymentID := uuid.New()

	repository := &mockDeploymentRepository{
		items: []deployments.Deployment{
			{
				ID:      deploymentID,
				Version: "v1.0.0",
				Status:  deployments.StatusPending,
			},
		},
	}

	reconciler := New(repository)

	if err := reconciler.Reconcile(
		context.Background(),
		deploymentID.String(),
	); err != nil {
		t.Fatalf("reconcile failed: %v", err)
	}

	deployment, err := repository.Get(
		context.Background(),
		deploymentID,
	)
	if err != nil {
		t.Fatalf("get deployment: %v", err)
	}

	if deployment.Status != deployments.StatusRunning {
		t.Fatalf(
			"expected status %q, got %q",
			deployments.StatusRunning,
			deployment.Status,
		)
	}
}

func TestReconcileIgnoresCompletedDeployment(t *testing.T) {
	deploymentID := uuid.New()

	repository := &mockDeploymentRepository{
		items: []deployments.Deployment{
			{
				ID:      deploymentID,
				Version: "v1.0.0",
				Status:  deployments.StatusSucceeded,
			},
		},
	}

	reconciler := New(repository)

	if err := reconciler.Reconcile(
		context.Background(),
		deploymentID.String(),
	); err != nil {
		t.Fatalf("reconcile failed: %v", err)
	}

	deployment, err := repository.Get(
		context.Background(),
		deploymentID,
	)
	if err != nil {
		t.Fatalf("get deployment: %v", err)
	}

	if deployment.Status != deployments.StatusSucceeded {
		t.Fatalf(
			"expected status to remain %q, got %q",
			deployments.StatusSucceeded,
			deployment.Status,
		)
	}
}

func TestReconcileRejectsInvalidDeploymentID(t *testing.T) {
	repository := &mockDeploymentRepository{}
	reconciler := New(repository)

	if err := reconciler.Reconcile(
		context.Background(),
		"not-a-uuid",
	); err == nil {
		t.Fatal("expected invalid deployment ID error")
	}
}

func TestReconcileMissingDeployment(t *testing.T) {
	repository := &mockDeploymentRepository{}
	reconciler := New(repository)

	err := reconciler.Reconcile(
		context.Background(),
		uuid.New().String(),
	)

	if err == nil {
		t.Fatal("expected missing deployment error")
	}
}

func (m *mockDeploymentRepository) ListPending(
	_ context.Context,
) ([]deployments.Deployment, error) {
	var result []deployments.Deployment

	for _, deployment := range m.items {
		if deployment.Status == deployments.StatusPending {
			result = append(result, deployment)
		}
	}

	return result, nil
}

func TestRunReconcilesPendingDeployments(t *testing.T) {
	firstID := uuid.New()
	secondID := uuid.New()

	repository := &mockDeploymentRepository{
		items: []deployments.Deployment{
			{
				ID:      firstID,
				Version: "v1.0.0",
				Status:  deployments.StatusPending,
			},
			{
				ID:      secondID,
				Version: "v2.0.0",
				Status:  deployments.StatusPending,
			},
		},
	}

	reconciler := New(repository)

	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan error, 1)

	go func() {
		done <- reconciler.Run(
			ctx,
			time.Hour,
		)
	}()

	// The first reconciliation happens immediately.
	cancel()

	err := <-done
	if err != context.Canceled {
		t.Fatalf(
			"expected context cancellation, got %v",
			err,
		)
	}

	for _, deployment := range repository.items {
		if deployment.Status != deployments.StatusRunning {
			t.Fatalf(
				"deployment %s: expected status %q, got %q",
				deployment.ID,
				deployments.StatusRunning,
				deployment.Status,
			)
		}
	}
}
