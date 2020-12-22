package conferences

import (
	"context"
	"testing"
)

func TestUpdateJob(t *testing.T) {

	job := &Job{
		CompanyName: "Unicorn",
		Title:       "Entry-level Software Engineer",
		Description: "At least 12 years experience with Go, You must hold at least 3 PhDs",
		Link:        "Uni.corn/Job",
		Discord:     "https://discord.gg/unicorn",
		Rank:        3,
		Approved:    true,
	}

	ctx := context.Background()
	response, err := CreateJob(ctx, &CreateJobParams{Job: job})

	if err != nil {
		t.Fatalf("failed to create job: %v", err)
	}

	job.ID = response.Job.ID

	updatedJob := &Job{
		ID:          response.Job.ID,
		CompanyName: "Unicorn Enterprises",
		Title:       "Entry-level Software Engineer",
		Description: "New graduates welcome to apply!",
		Link:        "Uni.corn/Job",
		Discord:     "https://discord.gg/unicorn",
		Rank:        1,
		Approved:    true,
	}

	result, err := UpdateJob(ctx, &UpdateJobParams{Job: updatedJob})
	if err != nil {
		t.Fatalf("failed to update job: %v", err)
	}

	if result.Job.ID != job.ID {
		t.Errorf("job id does not match got %v, want %v", result.Job.ID, job.ID)
	}

	if result.Job.CompanyName == job.CompanyName {
		t.Errorf("incorrect company name retrieved got %v want %v", result.Job.CompanyName, job.CompanyName)
	}

	if result.Job.Title != job.Title {
		t.Errorf("incorrect title retrieved got %v want %v", result.Job.Title, job.Title)
	}

	if result.Job.Description == job.Description {
		t.Errorf("incorrect description retrieved got %v want %v", result.Job.Description, job.Description)
	}

	if result.Job.Link != job.Link {
		t.Errorf("incorrect link retrieved got %v want %v", result.Job.Link, job.Link)
	}

	if result.Job.Discord != job.Discord {
		t.Errorf("incorrect discord retrieved got %v want %v", result.Job.Discord, job.Discord)
	}

	if result.Job.Rank == job.Rank {
		t.Errorf("incorrect rank retrieved got %v want %v", result.Job.Rank, job.Rank)
	}

}
