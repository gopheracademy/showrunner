package conferences

import (
	"context"
	"fmt"

	"encore.dev/storage/sqldb"
)

// DeletePaperParams defines the inputs used by the DeletePaper API method
type DeletePaperParams struct {
	PaperID uint32
}

// DeletePaper removes a specific paper from the db
//encore:api public
func DeletePaper(ctx context.Context, params *DeletePaperParams) error {

	result, err := sqldb.Exec(ctx,
		`DELETE FROM paper_submission WHERE id = $1`, params.PaperID)
	if err != nil {
		return fmt.Errorf("failed to delete paper: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("could not find paper with id %v", params.PaperID)
	}
	return nil
}
