// Copyright © ivanlobanov. All rights reserved.
package middlewares

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/cantylv/authorization-service/internal/entity/dto"
	f "github.com/cantylv/authorization-service/internal/utils/functions"
	mc "github.com/cantylv/authorization-service/internal/utils/myconstants"
	me "github.com/cantylv/authorization-service/internal/utils/myerrors"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// JWT --> header.payload.signature
// header --> base64(meta_information)
// payload --> base64(payload_data)
// signature --> base64(hmacsha512('header.payload' + secret))

//// e.g. header
// {
// 	"exp": "28.09.2024 15:04:05 UTC+03",
//  "type": "jwt",
//  "alg": "sha512"
// }
//// e.g. payload
// {
// 	"user_id": "c9b299f4-a56a-4ea4-bdd3-bdd222a789e2",
//  "user_ip_address": "10.25.232.7"
// }

// JwtVerification Needed for authdtoication.
func JwtVerification(h http.Handler, logger *zap.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID, ok := r.Context().Value(mc.AccessKey(mc.RequestID)).(string)
		if !ok {
			requestID = r.RemoteAddr
			logger.Error(me.ErrNoRequestIdInContext.Error(), zap.String(mc.RequestID, requestID))
			ctx := context.WithValue(r.Context(), mc.AccessKey(mc.RequestID), requestID)
			r = r.WithContext(ctx)
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Error(err.Error(), zap.String(mc.RequestID, requestID))
			f.Response(w, dto.ResponseError{Error: me.ErrInternal.Error()}, http.StatusInternalServerError)
			return
		}
		if len(body) != 0 {
			var jwtTokenBody map[string]string
			err = json.Unmarshal(body, &jwtTokenBody)
			if err != nil {
				logger.Error(err.Error(), zap.String(mc.RequestID, requestID))
				f.Response(w, dto.ResponseError{Error: me.ErrInternal.Error()}, http.StatusInternalServerError)
				return
			}
			jwtToken := jwtTokenBody["jwt_token"]
			if jwtToken != "" {
				payload, err := jwtTokenIsValid(jwtToken)
				if err != nil {
					logger.Error(err.Error(), zap.String(mc.RequestID, requestID))
					if errors.Is(err, me.ErrInvalidJwtToken) || errors.Is(err, me.ErrJwtAlreadyExpired) {
						f.Response(w, dto.ResponseError{Error: err.Error()}, http.StatusBadRequest)
						return
					}
					f.Response(w, dto.ResponseError{Error: me.ErrInternal.Error()}, http.StatusInternalServerError)
					return
				}
				ctx := context.WithValue(r.Context(), mc.AccessKey(mc.JwtPayload), payload)
				r = r.WithContext(ctx)
			}
		}
		h.ServeHTTP(w, r)
	})
}

// jwtTokenIsValid Needed for validation jwt-token.
func jwtTokenIsValid(token string) (*dto.JwtPayload, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, me.ErrInvalidJwtToken
	}
	signatureHash := f.HashData([]byte(parts[0]+"."+parts[1]), []byte(viper.GetString("secret_key")))
	signature := base64.StdEncoding.EncodeToString([]byte(signatureHash))
	if signature != parts[2] {
		return nil, me.ErrInvalidJwtToken
	}

	dataHeader, err := base64.StdEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, err
	}
	var h dto.JwtHeader
	err = json.Unmarshal(dataHeader, &h)
	if err != nil {
		return nil, err
	}
	// check date experation
	dateNow := time.Now()
	if h.Exp.Equal(dateNow) || dateNow.After(h.Exp) {
		return nil, me.ErrJwtAlreadyExpired
	}

	dataPayload, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, err
	}
	var p dto.JwtPayload
	err = json.Unmarshal(dataPayload, &p)
	if err != nil {
		return nil, err
	}
	return &p, nil
}
