package utils

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"time"

	errs "github.com/ElfAstAhe/goph-keeper/pkg/error"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"
)

const (
	DefaultJWTSigningMethodName  = "HS256"
	DefaultJWTExpirationDuration = 30 * time.Minute
)

var (
	DefaultJWTSigningMethod = jwt.GetSigningMethod(DefaultJWTSigningMethodName)
)

type TokenIDBuilder func() string

type AppClaims struct {
	jwt.RegisteredClaims
	Admin  bool   `json:"admin,omitempty"`
	UserID string `json:"user_id,omitempty"`
	Roles  Roles  `json:"roles,omitempty"`
}

func NewAppClaims(
	userID string,
	user string,
	admin bool,
	tokenIDBuilder TokenIDBuilder,
	issuer string,
	expirationDuration time.Duration,
	roles ...string,
) *AppClaims {
	return &AppClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        tokenIDBuilder(),
			Issuer:    issuer,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expirationDuration)),
			Subject:   user,
		},
		Admin:  admin,
		UserID: userID,
		Roles:  roles,
	}
}

type JWTHelper struct {
	signingMethod      jwt.SigningMethod
	secretKey          string
	expirationDuration time.Duration
	tokenIDBuilder     TokenIDBuilder
}

func NewJWTHelper(signingMethod jwt.SigningMethod, secretKey string, expirationDuration time.Duration, tokenIDBuilder TokenIDBuilder) *JWTHelper {
	return &JWTHelper{
		signingMethod:      signingMethod,
		secretKey:          secretKey,
		expirationDuration: expirationDuration,
		tokenIDBuilder:     tokenIDBuilder,
	}
}

func NewDefaultJWTHelper(secretKey string) *JWTHelper {
	return NewJWTHelper(DefaultJWTSigningMethod, secretKey, DefaultJWTExpirationDuration, defaultTokenIDBuilder)
}

func (h *JWTHelper) ExtractTokenStringFromCookie(cookieName string, r *http.Request) (string, error) {
	if strings.TrimSpace(cookieName) == "" {
		return "", errs.NewUtlJWTError("empty cookie name", nil)
	}
	if r == nil {
		return "", errs.NewUtlJWTError("nil HTTP Request", nil)
	}

	cookie, err := r.Cookie(cookieName)
	if err != nil {
		return "", errs.NewUtlJWTError(fmt.Sprintf("cookie [%s] extraction", cookieName), err)
	}
	if cookie == nil {
		return "", errs.NewUtlJWTError(fmt.Sprintf("cookie not found [%s]", cookieName), err)
	}
	if err = cookie.Valid(); err != nil {
		return "", errs.NewUtlJWTError(fmt.Sprintf("cookie [%s] is invalid", cookieName), err)
	}

	return cookie.Value, nil
}

func (h *JWTHelper) ExtractTokenStringFromGRPCMetadata(metadataName string, md metadata.MD) (string, error) {
	if strings.TrimSpace(metadataName) == "" {
		return "", errs.NewUtlJWTError("empty metadata name", nil)
	}
	if md == nil {
		return "", errs.NewUtlJWTError("nil metadata", nil)
	}

	values := md.Get(metadataName)
	if len(values) == 0 {
		return "", nil
	}

	return values[0], nil
}

func (h *JWTHelper) ExtractClaims(token *jwt.Token) (*AppClaims, error) {
	res, ok := token.Claims.(*AppClaims)
	if !ok {
		return nil, errs.NewUtlJWTError("invalid claims", nil)
	}

	return res, nil
}

func (h *JWTHelper) ExtractTokenFromString(tokenString string) (*jwt.Token, error) {
	claims := new(AppClaims)
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if reflect.TypeOf(h.signingMethod) != reflect.TypeOf(token.Method) {
			return nil, errs.NewUtlJWTError("invalid signing method", nil)
		}

		return []byte(h.secretKey), nil
	})

	if err != nil {
		if !errors.As(err, &errs.ErrUtlJWT) {
			return nil, errs.NewUtlJWTError("error parse jwt", err)
		}

		return nil, err
	}

	if !token.Valid {
		return nil, errs.NewUtlJWTError("token validation failed", nil)
	}

	return token, nil
}

func (h *JWTHelper) buildTokenID() string {
	if h.tokenIDBuilder != nil {
		return h.tokenIDBuilder()
	}

	return defaultTokenIDBuilder()
}

func defaultTokenIDBuilder() string {
	template := "undef-%v"
	res, err := uuid.NewRandom()
	if err != nil {
		return fmt.Sprintf(template, time.Now().Nanosecond())
	}

	return fmt.Sprintf(template, res.String())
}
