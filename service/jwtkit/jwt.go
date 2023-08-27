package jwtkit

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/mizuki1412/go-core-kit/class"
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/mizuki1412/go-core-kit/cli/configkey"
	"github.com/mizuki1412/go-core-kit/service/configkit"
	"github.com/spf13/cast"
	"time"
)

type Claims struct {
	Id    any             `json:"id"`
	Ext   class.MapString `json:"ext"`
	Valid bool            `json:"valid"`
	jwt.RegisteredClaims
}

// New id: string or int
func New(id any) Claims {
	return Claims{
		Id:    id,
		Valid: true,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(configkit.GetInt(configkey.JwtExpire)) * time.Hour)), // 过期时间
			IssuedAt:  jwt.NewNumericDate(time.Now()),                                                                       // 签发时间
			NotBefore: jwt.NewNumericDate(time.Now()),                                                                       // 生效时间
		},
	}
}

func (c Claims) Token() string {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	s, err := t.SignedString([]byte(configkit.GetString(configkey.JwtSecretKey)))
	if err != nil {
		panic(exception.New("jwt token err: " + err.Error()))
	}
	return s
}

func (c Claims) IdInt() int {
	return cast.ToInt(c.Id)
}
func (c Claims) IdInt32() int32 {
	return cast.ToInt32(c.Id)
}
func (c Claims) IdInt64() int64 {
	return cast.ToInt64(c.Id)
}
func (c Claims) IdStr() string {
	return cast.ToString(c.Id)
}

func (c Claims) IsValid() bool {
	if !c.Valid {
		return false
	}
	if c.ExpiresAt.Unix() < time.Now().Unix() {
		return false
	}
	return true
}

func Parse(token string) Claims {
	t, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (any, error) {
		return []byte(configkit.GetString(configkey.JwtSecretKey)), nil
	})

	if claims, ok := t.Claims.(*Claims); ok && t.Valid {
		return *claims
	} else {
		panic(exception.New("jwt parse err: " + err.Error()))
	}
}
