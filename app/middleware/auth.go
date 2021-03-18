package middleware

import (
	"fmt"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"log"
	UserUseCase "ocr.service.authorization/app/user/usecase"
	"ocr.service.authorization/config"
	"ocr.service.authorization/model"
	"ocr.service.authorization/module/salt"
	"os"
	"time"
)

// @tags User
// @Summary user login
// @Description create user
// @start_time default
// @Param body body model.UserLogin true "json"
// @Success 200 {object} model.LoginResponse ""
// @Router /api/v1/login [post]
func NewAuth() *jwt.GinJWTMiddleware {
	CONFIG, _ := config.NewConfig(nil)
	var identityKey = CONFIG.GetString("IDENTITY_KEY")
	userUseCase, err := UserUseCase.NewUserUseCase()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// the jwt middleware
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "test zone",
		Key:         []byte("secret key"),
		Timeout:     time.Hour,
		MaxRefresh:  2 * time.Hour,
		IdentityKey: identityKey,
		Authenticator: func(c *gin.Context) (interface{}, error) {
			fmt.Println("Authenticator")
			var userLogin model.UserLogin
			if err := c.ShouldBind(&userLogin); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			user, err := userUseCase.GetFull(model.User{Username: userLogin.Username})
			if err == nil && user.Status != "block" {
				isRightPassword := salt.ComparePasswords(user.Password, []byte(userLogin.Password))
				if isRightPassword {
					user := model.User{
						Id:       user.Id,
						Username: user.Username,
						Role:    user.Role,
					}
					c.Set("user", user)
					return &user, nil
				}
			}
			return nil, jwt.ErrFailedAuthentication
		},
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			fmt.Println("PayloadFunc", data.(*model.User))
			if v, ok := data.(*model.User); ok {
				return jwt.MapClaims{
					identityKey: v.Id,
					"role":     v.Role,
				}
			}
			return jwt.MapClaims{}
		},
		LoginResponse: func(c *gin.Context, i int, s string, t time.Time) {
			//claims := jwt.ExtractClaims(c)
			//fmt.Println(claims)
			user, _ := c.Get("user")
			var loginResponse = model.LoginResponse{
				Code:   i,
				Expire: t.Format(time.RFC3339),
				Token:  s,
				UserId: user.(model.User).Id,
			}
			c.JSON(i, loginResponse)
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			fmt.Println("IdentityHandler")
			claims := jwt.ExtractClaims(c)
			return &model.User{
				Id:    claims[identityKey].(string),
				Role: claims["role"].(string),
			}
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			fmt.Println("Authorizator")
			if v, ok := data.(*model.User); ok && (v.Role == "user" || v.Role == "admin") {
				return true
			}
			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			fmt.Println("Unauthorized")
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		// - "param:<name>"
		TokenLookup: "header: Authorization, query: token, cookie: jwt",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	})

	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	// When you use jwt.New(), the function is already automatically called for checking,
	// which means you don't need to call it again.
	errInit := authMiddleware.MiddlewareInit()

	if errInit != nil {
		log.Fatal("authMiddleware.MiddlewareInit() Error:" + errInit.Error())
	}
	return authMiddleware
}
