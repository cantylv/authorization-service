package tokens

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	ent "github.com/cantylv/authorization-service/internal/entity"
	"github.com/cantylv/authorization-service/internal/entity/dto"
	f "github.com/cantylv/authorization-service/internal/utils/functions"
	mc "github.com/cantylv/authorization-service/internal/utils/myconstants"
	me "github.com/cantylv/authorization-service/internal/utils/myerrors"
	"github.com/spf13/viper"
)

func getMetadataFromConnection(r *http.Request) (*ent.Session, error) {
	userAgent := r.UserAgent()
	if userAgent == "" {
		return nil, me.ErrInvalidUserAgent
	}
	userIpAddress := r.RemoteAddr
	if userAgent == "" {
		return nil, me.ErrInvalidRemoteIp
	}
	tokenCookie, err := r.Cookie(mc.RefreshToken)
	if err != nil && !errors.Is(err, http.ErrNoCookie) {
		return nil, err
	}
	refreshToken := ""
	if tokenCookie != nil {
		refreshToken = tokenCookie.Value
	}
	// best practise is to send js fingerprint after user authorization
	return &ent.Session{
		Fingerprint:   fmt.Sprintf("%s.%s", userIpAddress, userAgent),
		UserIpAddress: userIpAddress,
		RefreshToken:  refreshToken,
	}, nil
}

func getUserResponse(jwtToken string, user *ent.User) *dto.UserResponse {
	return &dto.UserResponse{
		Payload:  *user,
		JwtToken: jwtToken,
	}
}

func createJwtToken(token *ent.Session) (string, error) {
	jwtHeader := dto.JwtHeader{
		Alg:  "hc512",
		Type: "jwt",
		Exp:  token.ExpiresAt,
	}
	jwtHeaderJSON, err := json.Marshal(jwtHeader)
	if err != nil {
		return "", err
	}
	jwtHeaderEncoded := base64.StdEncoding.EncodeToString(jwtHeaderJSON)

	jwtPayload := dto.JwtPayload{
		UserIpAddress: token.UserIpAddress,
		UserId:        token.UserID,
	}
	jwtPayloadJSON, err := json.Marshal(jwtPayload)
	if err != nil {
		return "", err
	}
	jwtPayloadEncoded := base64.StdEncoding.EncodeToString(jwtPayloadJSON)

	signatureHash := f.HashData([]byte(jwtHeaderEncoded+"."+jwtPayloadEncoded), []byte(viper.GetString("secret_key")))
	signatureHashEncoded := base64.StdEncoding.EncodeToString([]byte(signatureHash))
	return strings.Join([]string{jwtHeaderEncoded, jwtPayloadEncoded, signatureHashEncoded}, "."), nil
}
