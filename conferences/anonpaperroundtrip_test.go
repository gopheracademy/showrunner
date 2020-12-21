package conferences

import (
	"context"
	"testing"
)

func TestGetAnonPaperRoundTrip(t *testing.T) {

	t.Run("adds a paper for a specific conference", func(t *testing.T) {

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

		result, err := GetAnonPaper(ctx, &GetAnonPaperParams{PaperID: response.PaperID})

		if err != nil {
			t.Fatalf("unexpected database error: %v", err)
		}

		if result.AnonPaper.Title != paper.Title {
			t.Errorf("incorrect title returned got %v want %v", result.AnonPaper.Title, paper.UserID)
		}

		if result.AnonPaper.ElevatorPitch != paper.ElevatorPitch {
			t.Errorf("incorrect elevator pitch returned got %v want %v", result.AnonPaper.ElevatorPitch, paper.ElevatorPitch)
		}

		if result.AnonPaper.Description != paper.Description {
			t.Errorf("incorrect description returned got %v want %v", result.AnonPaper.Description, paper.Description)
		}

		if result.AnonPaper.Notes != paper.Notes {
			t.Errorf("incorrect notes returned got %v want %v", result.AnonPaper.Notes, paper.Notes)
		}

	})

}
