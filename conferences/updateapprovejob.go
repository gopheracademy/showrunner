package conferences

import (
	"context"
	"database/sql"
	"fmt"

	"encore.dev/storage/sqldb"
)

// UpdateApproveJobParams defines the inputs used by the
// UpdateApproveJob API method
type UpdateApproveJobParams struct {
	JobID          uint32
	ApprovedStatus bool
}

// UpdateApproveJobResponse defines the output the returned
// by the ApproveJob API method
type UpdateApproveJobResponse struct {
	Job *Job
}

// UpdateApproveJob sets the approval of a job to true
// or false depending on input
// encore:api public
func UpdateApproveJob(ctx context.Context, params *UpdateApproveJobParams) (*UpdateApproveJobResponse, error) {

	row := sqldb.QueryRow(ctx, `
	UPDATE job_board
	SET approved = $1
	WHERE id = $2
	RETURNING id,
		company_name,
		title,
		description,
		link,
		discord,
		rank,
		approved
	`, params.ApprovedStatus, params.JobID)

	var job Job
	err := row.Scan(
		&job.ID,
		&job.CompanyName,
		&job.Title,
		&job.Description,
		&job.Link,
		&job.Discord,
		&job.Rank,
		&job.Approved,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no such job found")
		}
		return nil, fmt.Errorf("failed to update job: %w", err)
	}

	return &UpdateApproveJobResponse{Job: &job}, nil

}
