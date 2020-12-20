package conferences

import (
	"context"
	"testing"
)

func TestUpdatePaperSubmission(t *testing.T) {

	t.Run("checks that a user can update all fields in their proposal", func(t *testing.T) {

		originalPaper := &Paper{
			UserID:        "test_user_1",
			ConferenceID:  1,
			Title:         "Test title",
			ElevatorPitch: "Elevating elevator pitch",
			Description:   "Descriptive description",
			Notes:         "Notable Notes",
		}

		ctx := context.Background()
		response, err := AddPaper(ctx, &AddPaperParams{
			Paper: originalPaper,
		},
		)
		if err != nil {
			t.Fatalf("unexpected database error: %v", err)
		}

		updatedPaper := &Paper{
			ID:            response.PaperID,
			UserID:        "test_user_1",
			ConferenceID:  1,
			Title:         "Can anyone code?",
			ElevatorPitch: "Is anyone capeable of coding? Lets discuss",
			Description:   "What does it require for someone to learn to code? Intelligence? Mindset?",
			Notes:         "Target Audience: Anyone!",
		}

		result, err := UpdatePaper(ctx, &UpdatePaperParams{Paper: updatedPaper})

		if err != nil {
			t.Fatalf("unexpected database error: %v", err)
		}

		if result.Paper.UserID != originalPaper.UserID {
			t.Errorf("UserID was unexpectedly updated got %v want %v", result.Paper.UserID, originalPaper.UserID)
		}

		if result.Paper.Title == originalPaper.Title {
			t.Errorf("title was not updated got %v want %v", result.Paper.UserID, originalPaper.UserID)
		}

		if result.Paper.ElevatorPitch == originalPaper.ElevatorPitch {
			t.Errorf("elevator pitch was not updated got %v want %v", result.Paper.ElevatorPitch, originalPaper.ElevatorPitch)
		}

		if result.Paper.Description == originalPaper.Description {
			t.Errorf("description was not updated got %v want %v", result.Paper.Description, originalPaper.Description)
		}

		if result.Paper.Notes == originalPaper.Notes {
			t.Errorf("notes was not updated got %v want %v", result.Paper.Notes, originalPaper.Notes)
		}
	})
}
