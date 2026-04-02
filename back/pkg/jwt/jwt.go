package jwt

import (
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/vkotsiuba99/NeoHome/back/internal/domain"
)

const (
	claimKeyUserID = "user_id"
	claimKeyRole   = "role"
	claimKeyExp    = "exp"
	claimKeyIat    = "iat"
)

var (
	ErrInvalidToken  = errors.New("invalid token")
	ErrInvalidClaims = errors.New("invalid token claims")
)

func NewToken(user domain.User, duration time.Duration, secret []byte) (string, error) {
	if len(secret) == 0 {
		return "", ErrInvalidToken
	}
	if duration <= 0 {
		return "", ErrInvalidToken
	}

	now := time.Now().UTC()
	expirationTime := now.Add(duration)

	claims := jwt.MapClaims{
		claimKeyUserID: strconv.FormatInt(user.UserID, 10),
		claimKeyRole:   user.Role,
		claimKeyExp:    expirationTime.Unix(),
		claimKeyIat:    now.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ParseToken(tokenString string, secret []byte) (domain.TokenClaims, error) {
	if len(tokenString) == 0 || len(secret) == 0 {
		return domain.TokenClaims{}, ErrInvalidToken
	}

	parser := jwt.NewParser(jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	token, err := parser.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil || token == nil || !token.Valid {
		return domain.TokenClaims{}, ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return domain.TokenClaims{}, ErrInvalidClaims
	}

	userID, okUserID := toInt64(claims[claimKeyUserID])
	role, okRole := claims[claimKeyRole].(string)
	expUnix, okExp := toInt64(claims[claimKeyExp])
	iatUnix, okIat := toInt64(claims[claimKeyIat])
	if !okUserID || userID <= 0 || !okRole || len(role) == 0 || !okExp || !okIat {
		return domain.TokenClaims{}, ErrInvalidClaims
	}

	return domain.TokenClaims{
		UserID:    userID,
		Role:      role,
		ExpiresAt: expUnix,
		IssuedAt:  iatUnix,
	}, nil
}

func ParseUserClaims(claims domain.TokenClaims) (domain.User, error) {
	if claims.UserID <= 0 || len(claims.Role) == 0 {
		return domain.User{}, ErrInvalidClaims
	}

	return domain.User{
		UserID: claims.UserID,
		Role:   claims.Role,
	}, nil
}

func toInt64(value interface{}) (int64, bool) {
	switch typed := value.(type) {
	case float64:
		return int64(typed), true
	case int64:
		return typed, true
	case int:
		return int64(typed), true
	case string:
		parsed, err := strconv.ParseInt(typed, 10, 64)
		if err != nil {
			return 0, false
		}
		return parsed, true
	default:
		return 0, false
	}
}
