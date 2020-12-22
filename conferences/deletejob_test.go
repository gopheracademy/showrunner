package conferences

import (
	"context"
	"testing"
)

func TestDeletePaper(t *testing.T) {

	t.Run("checks a paper can be deleted by a specific id", func(t *testing.T) {

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

		err = DeleteJob(ctx, &DeleteJobParams{JobID: response.Job.ID})

		if err != nil {
			t.Errorf("failed to delete paper: %v", err)
		}
	})

	t.Run("attempts to delete a job that does not exist", func(t *testing.T) {

		ctx := context.Background()
		err := DeleteJob(ctx, &DeleteJobParams{JobID: 0})

		if err == nil {
			t.Errorf("failed to delete paper: %v", err)
		}
	})
}
