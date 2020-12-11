package conferences

import (
	"context"
	"encoding/json"
	"log"

	"encore.dev/beta/auth"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/oauth2"
)

var secrets struct {
	JWTKey     string
	OKTAPUBKEY string
}

// Data is the structure that Encore returns with information about the authenticated user
type Data struct {
	Exp              float64  `json:"exp"`
	Iat              float64  `json:"iat"`
	IdentityProvider string   `json:"identityProvider"`
	UserDetails      string   `json:"userDetails"`
	UserID           string   `json:"userId"`
	UserRoles        []string `json:"userRoles"`
}
type info struct {
	IDToken   string          `json:"id"`
	AuthToken string          `json:"auth"`
	UserInfo  json.RawMessage `json:"ui"`
}

// VerifyToken accepts a JWT token and returns a UserID and User Data, or an error.
// Return a zero-value UID for Unauthorized, return a non-nil error for a 500 error
// encore:authhandler
func VerifyToken(ctx context.Context, token string) (auth.UID, *Data, error) {
	provider, err := oidc.NewProvider(ctx, "https://dev-7217861.okta.com")
	if err != nil {
		log.Println("provider create error", err)
		return "", nil, err
	}
	log.Println(token)
	var i info
	err = json.Unmarshal([]byte(token), &i)

	if err != nil {
		log.Println("unmarshal token data struct", err)
		return "", nil, err
	}
	var verifier = provider.Verifier(&oidc.Config{ClientID: "0oa26dc0cgcjzHwsJ5d6"})
	idtok, err := verifier.Verify(ctx, i.IDToken)
	if err != nil {
		log.Println("verify token error: ", err)
		return "", nil, err
	}
	log.Println("idtok", idtok)
	var oidcToken oauth2.Token
	oidcToken.AccessToken = string(i.AuthToken)
	oidcToken.TokenType = "Bearer"

	ui, err := provider.UserInfo(ctx, oauth2.StaticTokenSource(&oidcToken))
	if err != nil {
		log.Println("user info error:", err)
		return "", nil, err
	}

	log.Println("ui", ui)
	var d Data
	ui.Claims(&d)
	//TODO: actually validate things from the token/claims
	//return auth.UID(d.UserID), d, nil
	return "", &d, nil
}

func mapClaims(values jwt.MapClaims) *Data {
	d := &Data{}
	d.Exp = values["exp"].(float64)
	d.Iat = values["iat"].(float64)
	d.IdentityProvider = values["identityProvider"].(string)
	d.UserID = values["userId"].(string)
	d.UserDetails = values["userDetails"].(string)
	for _, role := range values["userRoles"].([]interface{}) {
		rstr := role.(string)
		d.UserRoles = append(d.UserRoles, rstr)
	}
	return d

}
