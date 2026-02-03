package utils

import (
	"context"
	"maps"
	"net/http"
	"slices"
	"strings"

	errs "github.com/ElfAstAhe/goph-keeper/pkg/error"

	"github.com/golang-jwt/jwt/v4"
	"google.golang.org/grpc/metadata"
)

type ContextUserInfoType string

const (
	DefaultContextUserInfo ContextUserInfoType = "UserInfo"
)

const (
	DefaultCookieName   string = "Authorization"
	DefaultMetadataName string = "Authorization"
)

type Roles []string

// UserInfo - содержит информацию о правах доступа пользователя
type UserInfo struct {
	user   string
	userID string
	admin  bool
	roles  map[string]struct{}
}

func NewUserInfo(userID string, user string, admin bool, roles Roles) *UserInfo {
	return &UserInfo{
		user:   user,
		userID: userID,
		admin:  admin,
		roles:  toInternalRoles(roles),
	}
}

func (ui *UserInfo) User() string {
	return ui.user
}

func (ui *UserInfo) UserID() string {
	return ui.userID
}

// IsAdmin - признак администратора системы
func (ui *UserInfo) IsAdmin() bool {
	return ui.admin
}

func (ui *UserInfo) Roles() Roles {
	return slices.Collect(maps.Keys(ui.roles))
}

// InRole выполняет проверку прав.
//
// Если срез Roles пуст, метод всегда возвращает false.
// Сравнение ролей чувствительно к регистру.
//
// Пример использования:
//
//	if user.InRole("admin") {
//	    fmt.Println("Доступ разрешен")
//	}
func (ui *UserInfo) InRole(role string) bool {
	_, ok := ui.roles[role]

	return ok
}

type AuthHelper struct {
	contextUserInfo ContextUserInfoType
	cookieName      string
	metadataName    string
	jwtHelper       *JWTHelper
}

func NewAuthHelper(contextUserInfo ContextUserInfoType, cookieName, metadataName string, jwtHelper *JWTHelper) *AuthHelper {
	return &AuthHelper{
		contextUserInfo: contextUserInfo,
		cookieName:      cookieName,
		metadataName:    metadataName,
		jwtHelper:       jwtHelper,
	}
}

func NewDefaultAuthHelper(secretKey string) *AuthHelper {
	return NewDefaultAuthHelperEx(NewDefaultJWTHelper(secretKey))
}

func NewDefaultAuthHelperEx(jwtHelper *JWTHelper) *AuthHelper {
	return NewAuthHelper(DefaultContextUserInfo, DefaultCookieName, DefaultMetadataName, jwtHelper)
}

func (ah *AuthHelper) UserInfoFromToken(token *jwt.Token) (*UserInfo, error) {
	if token == nil {
		return nil, errs.NewUtlAuthError("nil jwt token", nil)
	}

	claims, err := ah.jwtHelper.ExtractClaims(token)
	if err != nil {
		return nil, errs.NewUtlAuthError("extract claims", err)
	}

	return NewUserInfo(claims.UserID, claims.Subject, claims.Admin, claims.Roles), nil
}

func (ah *AuthHelper) UserInfoFromTokenString(tokenString string) (*UserInfo, error) {
	token, err := ah.jwtHelper.ExtractTokenFromString(tokenString)
	if err != nil {
		return nil, errs.NewUtlAuthError("extract token", err)
	}

	return ah.UserInfoFromToken(token)
}

func (ah *AuthHelper) UserInfoFromContext(ctx context.Context) (*UserInfo, error) {
	if res, ok := ctx.Value(ah.contextUserInfo).(*UserInfo); ok {
		return res, nil
	}

	return nil, errs.NewUtlAuthError("user info not found", nil)
}

func (ah *AuthHelper) HasUserInfoInContext(ctx context.Context) bool {
	userInfo, err := ah.UserInfoFromContext(ctx)
	if err != nil {
		return false
	}

	return userInfo != nil
}

func (ah *AuthHelper) UserInfoFromHTTPRequest(r *http.Request) (*UserInfo, error) {
	tokenString, err := ah.jwtHelper.ExtractTokenStringFromCookie(ah.cookieName, r)
	if err != nil {
		return nil, errs.NewUtlAuthError("extract token string", err)
	}

	userInfo, err := ah.UserInfoFromTokenString(tokenString)
	if err != nil {
		return nil, errs.NewUtlAuthError("extract user info", err)
	}

	return userInfo, nil
}

func (ah *AuthHelper) UserInfoFromGRPCMetadata(md metadata.MD) (*UserInfo, error) {
	tokenString, err := ah.jwtHelper.ExtractTokenStringFromGRPCMetadata(ah.metadataName, md)
	if err != nil {
		return nil, errs.NewUtlAuthError("extract token string", err)
	}

	userInfo, err := ah.UserInfoFromTokenString(tokenString)
	if err != nil {
		return nil, errs.NewUtlAuthError("extract user info", err)
	}

	return userInfo, nil
}

func (ah *AuthHelper) UserInfoFromGRPCContext(gRPCCtx context.Context) (*UserInfo, error) {
	md, ok := metadata.FromIncomingContext(gRPCCtx)
	if !ok {
		return nil, errs.NewUtlAuthError("gRPC metadata not found", nil)
	}

	userInfo, err := ah.UserInfoFromGRPCMetadata(md)
	if err != nil {
		return nil, errs.NewUtlAuthError("gRPC metadata", err)
	}

	return userInfo, nil
}

func toInternalRoles(roles Roles) map[string]struct{} {
	res := make(map[string]struct{})
	for _, role := range roles {
		if strings.TrimSpace(strings.ToLower(role)) != "" {
			res[role] = struct{}{}
		}
	}

	return res
}
