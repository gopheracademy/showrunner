package conferences

import (
	"context"
	"testing"

	"encore.dev/storage/sqldb"
)

func TestUpdateSponsorContactInformation(t *testing.T) {
	ctx := context.Background()

	t.Run("update a sponsor contact", func(t *testing.T) {

		row := sqldb.QueryRow(ctx, `INSERT INTO sponsor (
		name, 
		address, 
		website, 
		sponsorship_level
	) VALUES (
		'Crowdstrike', 
		'Crow Tower, Strike City, 911 911', 
		'https://www.crowdstrike.com', 
		'platinum'
	)
	RETURNING id;`)

		var sponsorID uint32

		err := row.Scan(&sponsorID)
		assertDatabaseError(t, err)

		sponsorContactInformation := SponsorContactInformation{
			Name:  "Corey Crow",
			Role:  ContactRoleSoleContact,
			Email: "corey@crowdstrike.com",
			Phone: "555666777",
		}

		row = sqldb.QueryRow(ctx, `INSERT INTO sponsor_contact_information (
		name, 
		role, 
		email, 
		phone, 
		sponsor_id
	) VALUES (
		$1, 
		$2, 
		$3, 
		$4,
		$5
	)
	RETURNING id;`,
			sponsorContactInformation.Name,
			sponsorContactInformation.Role.String(),
			sponsorContactInformation.Email,
			sponsorContactInformation.Phone,
			sponsorID,
		)

		err = row.Scan(&sponsorContactInformation.ID)
		assertDatabaseError(t, err)

		updateSponsorContactInformation := SponsorContactInformation{
			ID:    sponsorContactInformation.ID,
			Name:  "Cory Crow",
			Role:  ContactRoleMarketing,
			Email: "cory@strikecrowd.com",
			Phone: sponsorContactInformation.Phone,
		}

		_, err = UpdateSponsorContact(ctx, &UpdateSponsorContactParams{
			SponsorContactInformation: &updateSponsorContactInformation,
		})
		assertDatabaseError(t, err)

		row = sqldb.QueryRow(ctx, `SELECT name, role, email, phone  FROM sponsor_contact_information WHERE id = $1;`, sponsorContactInformation.ID)

		var retrievedSponsorContactInformation SponsorContactInformation

		err = row.Scan(
			&retrievedSponsorContactInformation.Name,
			&retrievedSponsorContactInformation.Role,
			&retrievedSponsorContactInformation.Email,
			&retrievedSponsorContactInformation.Phone,
		)

		if retrievedSponsorContactInformation.Name == sponsorContactInformation.Name {
			t.Errorf(
				"name was not updated for contact: got %v want %v",
				retrievedSponsorContactInformation.Name,
				sponsorContactInformation.Name,
			)
		}

		if retrievedSponsorContactInformation.Role == sponsorContactInformation.Role {
			t.Errorf(
				"role was not updated for contact: got %v want %v",
				retrievedSponsorContactInformation.Role,
				sponsorContactInformation.Role,
			)
		}

		if retrievedSponsorContactInformation.Email == sponsorContactInformation.Email {
			t.Errorf(
				"email address was not updated for contact: got %v want %v",
				retrievedSponsorContactInformation.Email,
				sponsorContactInformation.Email,
			)
		}
	})

	t.Run("rejects an invalid role which is a positive int", func(t *testing.T) {

		sponsorContactInformation := SponsorContactInformation{
			Name:  "Corey Crow",
			Role:  5,
			Email: "corey@crowdstrike.com",
			Phone: "555666777",
		}

		_, err := UpdateSponsorContact(ctx, &UpdateSponsorContactParams{
			SponsorContactInformation: &sponsorContactInformation,
		})

		if err == nil {
			t.Fatalf("invalid role did not cause an error")
		}

	})

	t.Run("rejects an invalid role which is a negative int", func(t *testing.T) {

		sponsorContactInformation := SponsorContactInformation{
			Name:  "Corey Crow",
			Role:  -1,
			Email: "corey@crowdstrike.com",
			Phone: "555666777",
		}

		_, err := UpdateSponsorContact(ctx, &UpdateSponsorContactParams{
			SponsorContactInformation: &sponsorContactInformation,
		})

		if err == nil {
			t.Fatalf("invalid role did not cause an error")
		}

	})

	t.Run("returns error if no sponsor contact information is provided", func(t *testing.T) {

		_, err := UpdateSponsorContact(ctx, &UpdateSponsorContactParams{
			SponsorContactInformation: nil,
		})

		if err == nil {
			t.Fatalf("nil entry did not cause an error")
		}

	})

}

func assertDatabaseError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected database error: %v", err)
	}
}
