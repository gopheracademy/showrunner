package conferences

import (
	"context"
	"testing"

	"encore.dev/storage/sqldb"
)

func TestAddPaperRoundTrip(t *testing.T) {

	t.Run("adds a paper for a specific conference", func(t *testing.T) {

		tx, err := sqldb.Begin(context.TODO())
		if err != nil {
			t.Fatalf("beginning transaction: %v", err)
		}
		att01 := &User{
			Email:       "testmail01@gophercon.com",
			CoCAccepted: true,
		}
		savedAttendee01, err := createAttendee(context.TODO(), tx, att01)
		if err := sqldb.Commit(tx); err != nil {
			t.Fatalf("committing test setup transaction: %v", err)
		}
		paper := &Paper{
			UserID:        savedAttendee01.ID,
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

		result, err := GetPaper(ctx, &GetPaperParams{PaperID: response.PaperID})

		if err != nil {
			t.Fatalf("unexpected database error: %v", err)
		}

		if result.Paper.UserID != paper.UserID {
			t.Errorf("incorrect UserID returned got %v want %v", result.Paper.UserID, paper.UserID)
		}

		if result.Paper.Title != paper.Title {
			t.Errorf("incorrect title returned got %v want %v", result.Paper.UserID, paper.UserID)
		}

		if result.Paper.ElevatorPitch != paper.ElevatorPitch {
			t.Errorf("incorrect elevator pitch returned got %v want %v", result.Paper.ElevatorPitch, paper.ElevatorPitch)
		}

		if result.Paper.Description != paper.Description {
			t.Errorf("incorrect description returned got %v want %v", result.Paper.Description, paper.Description)
		}

		if result.Paper.Notes != paper.Notes {
			t.Errorf("incorrect notes returned got %v want %v", result.Paper.Notes, paper.Notes)
		}

	})

}
