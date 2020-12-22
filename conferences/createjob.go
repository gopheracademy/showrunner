package conferences

import (
	"context"
	"fmt"

	"encore.dev/storage/sqldb"
)

// CreateJobParams defines the inputs used by the CreateJob API method
type CreateJobParams struct {
	Job *Job
}

// CreateJobResponse defines the output returned by the CreateJob API method
type CreateJobResponse struct {
	Job *Job
}

// CreateJob inserts a job into the job_board table
// encore:api public
func CreateJob(ctx context.Context, params *CreateJobParams) (*CreateJobResponse, error) {

	row := sqldb.QueryRow(ctx,
		`INSERT INTO job_board (
			company_name,
			title,
			description,
			link,
			discord,
			rank
		) VALUES (
			$1,
			$2,
			$3,
			$4,
			$5,
			$6
			) RETURNING id,
			company_name,
			title,
			description,
			link,
			discord,
			rank`,
		params.Job.CompanyName,
		params.Job.Title,
		params.Job.Description,
		params.Job.Link,
		params.Job.Discord,
		params.Job.Rank,
	)

	var job Job
	err := row.Scan(
		&job.ID,
		&job.CompanyName,
		&job.Title,
		&job.Description,
		&job.Link,
		&job.Discord,
		&job.Rank,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to add job: %w", err)
	}

	return &CreateJobResponse{Job: &job}, nil
}
