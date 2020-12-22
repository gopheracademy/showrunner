package conferences

import (
	"context"
	"database/sql"
	"fmt"

	"encore.dev/storage/sqldb"
)

// UpdateJobParams defines the input used by
// the UpdateJob API method
type UpdateJobParams struct {
	Job *Job
}

// UpdateJobResponse defines the output returned
// by the UpdateJob API method
type UpdateJobResponse struct {
	Job *Job
}

// UpdateJob updates a job entry based on id in the
// job_board table
// encore:api public
func UpdateJob(ctx context.Context, params *UpdateJobParams) (*UpdateJobResponse, error) {

	row := sqldb.QueryRow(ctx, `
	UPDATE job_board
	SET company_name = $1,
		title = $2,
		description = $3,
		link = $4,
		discord = $5,
		rank = $6
	WHERE id = $7
	RETURNING id,
		company_name,
		title,
		description,
		link,
		discord,
		rank
	`, params.Job.CompanyName, params.Job.Title, params.Job.Description, params.Job.Link, params.Job.Discord, params.Job.Rank, params.Job.ID)

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
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no such job found")
		}
		return nil, fmt.Errorf("failed to update job: %w", err)
	}

	return &UpdateJobResponse{Job: &job}, nil
}
