package conferences

import (
	"context"
	"testing"
)

func TestDeletePaper(t *testing.T) {

	t.Run("user can delete a specific paper that they have submitted", func(t *testing.T) {

		paper := &Paper{
			UserID:        "test_user_1",
			ConferenceID:  1,
			Title:         "Test title",
			ElevatorPitch: "Elevating elevator pitch",
			Description:   "Descriptive description",
			Notes:         "Notable Notes",
		}

		ctx := context.Background()
		response, err := AddPaper(ctx, &AddPaperParams{
			Paper: paper,
		},
		)
		if err != nil {
			t.Fatalf("unexpected database error: %v", err)
		}

		err = DeletePaper(ctx, &DeletePaperParams{
			PaperID: response.PaperID,
		})

		if err != nil {
			t.Errorf("paper was not deleted: %w", err)
		}
	})

	t.Run("checks that an error is returned if a delete is attempted on a paperid that does not exist", func(t *testing.T) {

		ctx := context.Background()
		err := DeletePaper(ctx, &DeletePaperParams{PaperID: 0})
		if err == nil {
			t.Errorf("paper that does not exist was deleted")
		}
	})
}
