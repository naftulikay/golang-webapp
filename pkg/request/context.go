package request

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/naftulikay/golang-webapp/pkg/auth"
	"github.com/naftulikay/golang-webapp/pkg/constants"
	"github.com/naftulikay/golang-webapp/pkg/interfaces"
	"github.com/naftulikay/golang-webapp/pkg/logging"
	"go.uber.org/zap"
	"net/http"
)

var _ interfaces.RequestContext = (*RequestContextImpl)(nil)

var V validator.Validate

type RequestContextImpl struct {
	id            *uuid.UUID
	app           interfaces.App
	logger        *zap.Logger
	req           *http.Request
	authenticated bool
	token         *jwt.Token
	claims        *auth.JWTClaims
}

func (r RequestContextImpl) ID() uuid.UUID {
	return *r.id
}

func (r RequestContextImpl) App() interfaces.App {
	return r.app
}

func (r RequestContextImpl) Config() interfaces.AppConfig {
	return r.app.Config()
}

func (r RequestContextImpl) Services() interfaces.AppServices {
	return r.app.Service()
}

func (r RequestContextImpl) Dao() interfaces.AppDaos {
	return r.app.Dao()
}

func (r RequestContextImpl) Logger() *zap.Logger {
	return r.logger
}

func (r RequestContextImpl) Req() *http.Request {
	return r.req
}

func (r RequestContextImpl) Context() context.Context {
	return r.req.Context()
}

func (r RequestContextImpl) Authenticated() bool {
	return r.authenticated
}

func (r RequestContextImpl) IsUser() bool {
	return r.authenticated && r.claims.Role == "user" || r.claims.Role == "admin"
}

func (r RequestContextImpl) IsAdmin() bool {
	return r.authenticated && r.claims.Role == "admin"
}

func (r RequestContextImpl) Token() *jwt.Token {
	return r.token
}

func (r RequestContextImpl) Claims() *auth.JWTClaims {
	return r.claims
}

func (r RequestContextImpl) JSON(dest interface{}) error {
	if err := json.NewDecoder(r.Req().Body).Decode(dest); err != nil {
		return err
	}

	return nil
}

func (r RequestContextImpl) JSONV(dest interface{}) error {
	if err := r.JSON(dest); err != nil {
		return err
	}

	if err := V.Struct(dest); err != nil {
		return err
	}

	return nil
}

func NewRequestContext(app interfaces.App, logger *zap.Logger, r *http.Request) (*RequestContextImpl, error) {
	if r.Context().Value(constants.ContextKeyRequestID) == nil {
		return nil, fmt.Errorf("request id not present in the request context, middleware not configured")
	}

	reqID, ok := r.Context().Value(constants.ContextKeyRequestID).(uuid.UUID)

	if !ok {
		return nil, fmt.Errorf("unable to extract request id as a UUID")
	}

	authenticated, ok := r.Context().Value(constants.ContextKeyAuthenticated).(bool)

	if !ok {
		logger.Error("Unable to extract 'authenticated' bool field from the request context.")
	}

	var token *jwt.Token
	var claims *auth.JWTClaims

	if authenticated {
		t, ok := r.Context().Value(constants.ContextKeyJWTToken).(jwt.Token)

		if !ok {
			logger.Error("Unable to extract 'jwt_token' from the request context.")
		} else {
			token = &t
		}

		c, ok := r.Context().Value(constants.ContextKeyJWTClaims).(auth.JWTClaims)

		if !ok {
			logger.Error("Unable to extract 'jwt_claims' from the request context.")
		} else {
			claims = &c
		}
	}

	fields := logging.ZapFieldsFromContext(r.Context())

	return &RequestContextImpl{
		id:            &reqID,
		app:           app,
		logger:        logger.With(fields...),
		req:           r,
		authenticated: authenticated,
		token:         token,
		claims:        claims,
	}, nil
}

func init() {
	V = *validator.New()
}
