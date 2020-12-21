package conferences

import (
	"context"
	"testing"
)

func TestListAnonPapers(t *testing.T) {

	t.Run("returns all papers without identifying information", func(t *testing.T) {

		paper := &Paper{
			UserID:        "test_user_1",
			ConferenceID:  1,
			Title:         "Test title",
			ElevatorPitch: "Elevating elevator pitch",
			Description:   "Descriptive description",
			Notes:         "Notable Notes",
		}

		ctx := context.Background()
		_, err := AddPaper(ctx, &AddPaperParams{
			Paper: paper,
		},
		)
		if err != nil {
			t.Fatalf("unexpected database error: %v", err)
		}

		result, err := ListAnonPapers(ctx,
			&ListAnonPapersParams{ConferenceID: paper.ConferenceID},
		)
		if err != nil {
			t.Fatalf("failed retrieve papers: %v", err)
		}

		resultLength := len(result.AnonPapers) - 1

		if result.AnonPapers[resultLength].Title != paper.Title {
			t.Errorf("title was not as expected got %v want %v", result.AnonPapers[resultLength].Title, paper.Title)
		}

		if result.AnonPapers[resultLength].ElevatorPitch != paper.ElevatorPitch {
			t.Errorf("elevator pitch was not as expected got %v want %v", result.AnonPapers[resultLength].ElevatorPitch, paper.ElevatorPitch)
		}

		if result.AnonPapers[resultLength].Description != paper.Description {
			t.Errorf("elevator pitch was not as expected got %v want %v", result.AnonPapers[resultLength].Description, paper.Description)
		}

		if result.AnonPapers[resultLength].Notes != paper.Notes {
			t.Errorf("elevator pitch was not as expected got %v want %v", result.AnonPapers[resultLength].Notes, paper.Notes)
		}

	})
}
