package conferences

import (
	"context"
	"testing"
)

func TestApproveJob(t *testing.T) {

	t.Run("job can have approval set to true", func(t *testing.T) {

		job := &Job{
			CompanyName: "Unicorn",
			Title:       "Entry-level Software Engineer",
			Description: "At least 12 years experience with Go, You must hold at least 3 PhDs",
			Link:        "Uni.corn/Job",
			Discord:     "https://discord.gg/unicorn",
			Rank:        3,
		}

		ctx := context.Background()
		response, err := CreateJob(ctx, &CreateJobParams{Job: job})

		if err != nil {
			t.Fatalf("failed to create job: %v", err)
		}

		approvedStatus := true

		result, err := UpdateApproveJob(ctx, &UpdateApproveJobParams{JobID: response.Job.ID, ApprovedStatus: approvedStatus})

		if err != nil {
			t.Fatalf("failed to set approved flag to true: %v", err)
		}

		if result.Job.Approved != true {
			t.Errorf("flag was not set to true got %v want %v", result.Job.Approved, approvedStatus)
		}
	})

	t.Run("job can have approval set to true", func(t *testing.T) {

		job := &Job{
			CompanyName: "Big Cake",
			Title:       "Software Engineer",
			Description: "Help selling large cakes online",
			Link:        "Ca.Ke/Job",
			Discord:     "https://discord.gg/cake",
			Rank:        5,
		}

		ctx := context.Background()
		response, err := CreateJob(ctx, &CreateJobParams{Job: job})

		if err != nil {
			t.Fatalf("failed to create job: %v", err)
		}

		approvedStatus := true

		result, err := UpdateApproveJob(ctx, &UpdateApproveJobParams{JobID: response.Job.ID, ApprovedStatus: approvedStatus})

		if err != nil {
			t.Fatalf("failed to set approved flag to true: %v", err)
		}

		if result.Job.Approved != true {
			t.Fatalf("job approved flag should be true got %v want %v", response.Job.Approved, true)
		}

		falseApprovedStatus := false

		result, err = UpdateApproveJob(ctx, &UpdateApproveJobParams{JobID: response.Job.ID, ApprovedStatus: falseApprovedStatus})

		if err != nil {
			t.Fatalf("failed to set approved flag to true: %v", err)
		}

		if result.Job.Approved != false {
			t.Errorf("flag was not set to true got %v want %v", result.Job.Approved, falseApprovedStatus)
		}
	})

}
