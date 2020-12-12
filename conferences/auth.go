package conferences

import (
	"context"
	"fmt"
	"log"

	"encore.dev/beta/auth"
	"encore.dev/storage/sqldb"
	"github.com/coreos/go-oidc/v3/oidc"
)

// VerifyToken accepts a JWT token and returns a UserID, or an error.
// Return a zero-value UID for Unauthorized, return a non-nil error for a 500 error
// encore:authhandler
func VerifyToken(ctx context.Context, token string) (auth.UID, error) {
	provider, err := oidc.NewProvider(ctx, "https://dev-7217861.okta.com")
	if err != nil {
		log.Println("provider create error", err)
		// return nil error and zero value id to trigger unauthorized response
		return "", nil
	}
	var verifier = provider.Verifier(&oidc.Config{ClientID: "0oa26dc0cgcjzHwsJ5d6"})
	idt, err := verifier.Verify(ctx, token)
	if err != nil {
		log.Println("verify token error: ", err)
		// return nil error and zero value id to trigger unauthorized response
		return "", nil
	}

	var sqlStatement = `
	INSERT INTO users (id)
	VALUES $1
	ON CONFLICT (id) DO NOTHING`

	_, err = sqldb.Exec(ctx, sqlStatement, idt.Subject)
	if err != nil {
		return "", fmt.Errorf("failed to insert user: %w", err)
	}

	return auth.UID(idt.Subject), nil
}
