package conferences

import (
	"context"
	"fmt"

	"encore.dev/storage/sqldb"
)

// GetJobParams defines the inputs used by the GetJob API method
type GetJobParams struct {
	JobID uint32
}

// GetJobResponse defines the output returned by the GetJob API method
type GetJobResponse struct {
	Job *Job
}

// GetJob retrieves a job posting by JobID
func GetJob(ctx context.Context, params *GetJobParams) (*GetJobResponse, error) {

	row := sqldb.QueryRow(
		ctx,
		`
		SELECT id,
		company_name,
		title,
		description,
		link,
		discord,
		rank
		FROM job_board
		WHERE id = $1
		`, params.JobID,
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
		return nil, fmt.Errorf("failed to retrieve job: %w", err)
	}

	return &GetJobResponse{Job: &job}, nil
}
