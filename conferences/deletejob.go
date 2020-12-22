package conferences

import (
	"context"
	"fmt"

	"encore.dev/storage/sqldb"
)

// DeleteJobParams defines the input used by
// the DeleteJob API method
type DeleteJobParams struct {
	JobID uint32
}

// DeleteJob deletes a job by id from the
// job_board table
// encore:api public
func DeleteJob(ctx context.Context, params *DeleteJobParams) error {

	result, err := sqldb.Exec(ctx,
		`
		DELETE FROM job_board
		WHERE id = $1
		`,
		params.JobID,
	)
	if err != nil {
		return fmt.Errorf("failed to delete job: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to retrieve rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no such job found")
	}

	return nil

}
