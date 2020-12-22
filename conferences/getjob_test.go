package conferences

import (
	"context"
	"testing"
)

func TestCreateandGetJob(t *testing.T) {

	t.Run("create a job and retrieve by JobID", func(t *testing.T) {

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

		result, err := GetJob(ctx, &GetJobParams{JobID: response.Job.ID})

		if err != nil {
			t.Fatalf("job was not retrived: %v", err)
		}

		if result.Job.CompanyName != job.CompanyName {
			t.Errorf("incorrect company name retrieved got %v want %v", result.Job.CompanyName, job.CompanyName)
		}

		if result.Job.Title != job.Title {
			t.Errorf("incorrect title retrieved got %v want %v", result.Job.Title, job.Title)
		}

		if result.Job.Description != job.Description {
			t.Errorf("incorrect description retrieved got %v want %v", result.Job.Description, job.Description)
		}

		if result.Job.Link != job.Link {
			t.Errorf("incorrect link retrieved got %v want %v", result.Job.Link, job.Link)
		}

		if result.Job.Discord != job.Discord {
			t.Errorf("incorrect discord retrieved got %v want %v", result.Job.Discord, job.Discord)
		}

		if result.Job.Rank != job.Rank {
			t.Errorf("incorrect rank retrieved got %v want %v", result.Job.Rank, job.Rank)
		}
	})

	t.Run("behaves correctly when a job that does not exist is requested", func(t *testing.T) {

		ctx := context.Background()

		_, err := GetJob(ctx, &GetJobParams{JobID: 0})

		if err == nil {
			t.Fatalf("did not get an error when retrieving non existant job: %v", err)
		}
	})
}
