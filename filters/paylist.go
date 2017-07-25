package filters

import (
	"github.com/astaxie/beego/context"
	"strings"
	"github.com/dgrijalva/jwt-go"
)

const SecretKey = "jwtprime"

func PaylistFilter(ctx * context.Context) {
	bytoken := strings.TrimSpace(ctx.Input.Header("Authorization"))
	token, err := jwt.Parse(bytoken, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})
	if err != nil {
		ctx.Redirect(302, "/v1/user/deny")
	}else {
		if !token.Valid 	{
			ctx.Redirect(302, "/v1/user/deny")
		}
	}
}

