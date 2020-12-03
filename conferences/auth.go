package conferences

import (
	"context"

	"encore.dev/beta/auth"
	"encore.dev/rlog"
	"github.com/dgrijalva/jwt-go"
)

var secrets struct {
	JWTKey string
}

type Data struct {
	Exp              float64  `json:"exp"`
	Iat              float64  `json:"iat"`
	IdentityProvider string   `json:"identityProvider"`
	UserDetails      string   `json:"userDetails"`
	UserID           string   `json:"userId"`
	UserRoles        []string `json:"userRoles"`
}

// encore:authhandler
func VerifyToken(ctx context.Context, token string) (auth.UID, *Data, error) {
	// Look up in database etc...
	rlog.Info("decrypting token", "token", token)
	tok, err := jwt.Parse(token, nil)
	if tok == nil {
		return auth.UID("ERROR"), nil, err
	}
	claims, _ := tok.Claims.(jwt.MapClaims)
	d := mapClaims(claims)
	return auth.UID(d.UserID), d, nil
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
