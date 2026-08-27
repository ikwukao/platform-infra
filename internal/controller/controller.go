// Package controller contains the Platform-Infra reconciliation logic.
package controller

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/ikwukao/platform-infra/internal/deployments"
)

// Controller reconciles platform resources toward their desired state.
type Controller struct {
	deployments deployments.Repository
}

// New creates a deployment controller.
func New(
	deploymentRepository deployments.Repository,
) *Controller {
	return &Controller{
		deployments: deploymentRepository,
	}
}

// UUID Helper
func parseUUID(value string) (uuid.UUID, error) {
	return uuid.Parse(value)
}

// Reconcile processes a deployment and advances its lifecycle.
//
// Execution of the actual workload will be added later. For now,
// reconciliation validates that a deployment can move from pending
// to running.
func (c *Controller) Reconcile(
	ctx context.Context,
	deploymentID string,
) error {
	id, err := parseUUID(deploymentID)
	if err != nil {
		return fmt.Errorf("parse deployment id: %w", err)
	}

	deployment, err := c.deployments.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("get deployment: %w", err)
	}

	if deployment.Status != deployments.StatusPending {
		return nil
	}

	if err := c.deployments.UpdateStatus(
		ctx,
		deployment.ID,
		deployments.StatusRunning,
	); err != nil {
		return fmt.Errorf("mark deployment running: %w", err)
	}

	return nil
}

// Run continuously reconciles pending deployments until the context is cancelled.
func (c *Controller) Run(ctx context.Context, interval time.Duration) error {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		if err := c.reconcilePending(ctx); err != nil {
			return err
		}

		select {
		case <-ctx.Done():
			return ctx.Err()

		case <-ticker.C:
		}
	}
}

func (c *Controller) reconcilePending(ctx context.Context) error {
	items, err := c.deployments.ListPending(ctx)
	if err != nil {
		return fmt.Errorf("list pending deployments: %w", err)
	}

	for _, deployment := range items {
		if err := c.Reconcile(
			ctx,
			deployment.ID.String(),
		); err != nil {
			return fmt.Errorf(
				"reconcile deployment %s: %w",
				deployment.ID,
				err,
			)
		}
	}

	return nil
}
