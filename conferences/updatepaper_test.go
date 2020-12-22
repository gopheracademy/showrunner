package conferences

import (
	"context"
	"testing"

	"encore.dev/storage/sqldb"
)

func TestUpdatePaperSubmission(t *testing.T) {

	t.Run("checks that a user can update all fields in their proposal", func(t *testing.T) {
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
		originalPaper := &Paper{
			UserID:        savedAttendee01.ID,
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
			UserID:        savedAttendee01.ID,
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

	t.Run("checks a user can update some of their proposal", func(t *testing.T) {
		tx, err := sqldb.Begin(context.TODO())
		if err != nil {
			t.Fatalf("beginning transaction: %v", err)
		}
		att02 := &User{
			Email:       "testmail02@gophercon.com",
			CoCAccepted: true,
		}
		savedAttendee02, err := createAttendee(context.TODO(), tx, att02)
		if err := sqldb.Commit(tx); err != nil {
			t.Fatalf("committing test setup transaction: %v", err)
		}
		originalPaper := &Paper{
			UserID:        savedAttendee02.ID,
			ConferenceID:  2,
			Title:         "Go ahead with Go",
			ElevatorPitch: "Why you should code in Go",
			Description:   "Because the mascot is the best",
			Notes:         "It comes in all shapes and sizes",
		}

		ctx := context.Background()
		response, err := AddPaper(ctx, &AddPaperParams{
			Paper: originalPaper,
		},
		)
		if err != nil {
			t.Fatalf("unexpected database error: %v", err)
		}

		updatedPaper := *originalPaper

		updatedPaper.ID = response.PaperID
		updatedPaper.Title = "Get Great with Go"
		updatedPaper.Notes = "Target audience: New Gophers"

		result, err := UpdatePaper(ctx, &UpdatePaperParams{Paper: &updatedPaper})

		if err != nil {
			t.Fatalf("unexpected database error: %v", err)
		}

		if result.Paper.UserID != originalPaper.UserID {
			t.Errorf("UserID was unexpectedly updated got %v want %v", result.Paper.UserID, originalPaper.UserID)
		}

		if result.Paper.Title == originalPaper.Title {
			t.Errorf("title was not updated got %v want %v", result.Paper.UserID, originalPaper.UserID)
		}

		if result.Paper.ElevatorPitch != originalPaper.ElevatorPitch {
			t.Errorf("elevator pitch was not updated got %v want %v", result.Paper.ElevatorPitch, originalPaper.ElevatorPitch)
		}

		if result.Paper.Description != originalPaper.Description {
			t.Errorf("description was not updated got %v want %v", result.Paper.Description, originalPaper.Description)
		}

		if result.Paper.Notes == originalPaper.Notes {
			t.Errorf("notes was not updated got %v want %v", result.Paper.Notes, originalPaper.Notes)
		}
	})
}
