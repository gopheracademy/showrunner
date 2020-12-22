package conferences

import (
	"context"
	"fmt"

	"encore.dev/storage/sqldb"
)

// ListApprovedJobsParams defines the inputs used by
// the ListJobs API method
type ListApprovedJobsParams struct {
}

// ListApprovedJobsResponse defines the output returned
// by the ListApprovedJobs API method
type ListApprovedJobsResponse struct {
	Jobs []Job
}

// ListApprovedJobs retrieves all jobs (approved or not) from
// the job_board table
// encore:api public
func ListApprovedJobs(ctx context.Context) (*ListApprovedJobsResponse, error) {

	rows, err := sqldb.Query(ctx,
		`
		SELECT id,
			company_name,
			title,
			description,
			link,
			discord,
			rank,
			approved
		FROM job_board
		WHERE approved = true
		ORDER BY rank
`)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve all jobs: %w", err)
	}

	defer rows.Close()

	var jobs []Job

	for rows.Next() {

		var job Job

		err := rows.Scan(
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
			return nil, fmt.Errorf("failed to scan rows: %w", err)
		}

		jobs = append(jobs, job)
	}

	return &ListApprovedJobsResponse{Jobs: jobs}, nil

}
