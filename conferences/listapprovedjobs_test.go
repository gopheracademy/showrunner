package conferences

import (
	"context"
	"testing"
)

func TestListApprovedJobs(t *testing.T) {

	t.Run("expects a known approved job to be returned", func(t *testing.T) {
		ctx := context.Background()
		result, err := ListApprovedJobs(ctx)

		if err != nil {
			t.Fatalf("failed to retrieve jobs: %v", err)
		}

		resultsLength := len(result.Jobs)

		if resultsLength != 1 {
			t.Errorf("unexpected number of results returned got %v want %v", resultsLength, 1)
		}

		if result.Jobs[0].Approved != true {
			t.Errorf("unapproved job was returned")
		}
	})
}
