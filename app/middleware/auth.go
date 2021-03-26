package middleware

import (
	"encoding/json"
	"fmt"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"log"
	UserUseCase "ocr.service.authorization/app/user/usecase"
	UserLogUseCase "ocr.service.authorization/app/user_log/usecase"
	"ocr.service.authorization/config"
	"ocr.service.authorization/enum"
	"ocr.service.authorization/model"
	"os"
	"time"
)

// @tags User
// @Summary user login
// @Description user login
// @start_time default
// @Param body body model.UserLogin true "json"
// @Success 200 {object} model.LoginResponse ""
// @Router /api/v1/login [post]
func NewAuth() *jwt.GinJWTMiddleware {
	CONFIG, _ := config.NewConfig(nil)
	var identityKey = CONFIG.GetString("IDENTITY_KEY")
	var secret = CONFIG.GetString("SECRET")
	userUseCase, err := UserUseCase.NewUserUseCase()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	userLogUseCase, err := UserLogUseCase.NewUserLogUseCase()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// the jwt middleware
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "test zone",
		Key:         []byte(secret),
		Timeout:     CONFIG.GetDuration("TOKEN_EXPIRE_TIME"),
		MaxRefresh:  time.Second,
		IdentityKey: identityKey,
		Authenticator: func(c *gin.Context) (interface{}, error) {
			fmt.Println("Authenticator")
			var userLogin model.UserLogin
			if err := c.ShouldBind(&userLogin); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			user, err := userUseCase.Login(userLogin)
			if err == nil {
				user := model.User{
					Id:       user.Id,
					Username: user.Username,
					Role:     user.Role,
				}
				c.Set("user", user)
				return &user, nil
			}
			return nil, jwt.ErrFailedAuthentication
		},
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			fmt.Println("PayloadFunc", data.(*model.User))
			if v, ok := data.(*model.User); ok {
				var jwtClaim jwt.MapClaims
				byteClaim, _ := json.Marshal(model.Claim{UserId: v.Id, Role: v.Role})
				json.Unmarshal(byteClaim, &jwtClaim)
				return jwtClaim
			}
			return jwt.MapClaims{}
		},
		LoginResponse: func(c *gin.Context, i int, s string, t time.Time) {
			//claims := jwt.ExtractClaims(c)
			//fmt.Println(claims)
			iUser, _ := c.Get("user")
			user := iUser.(model.User)
			var loginResponse = model.LoginResponse{
				Code:   i,
				Expire: t.Format(time.RFC3339),
				Token:  s,
				UserId: user.Id,
			}
			userLog := model.UserLog{
				UserId:      user.Id,
				CreateAt:    time.Now().Format(time.RFC3339),
				ExpiredTime: t.Format(time.RFC3339),
				Ip:          "",
				Mac:         "",
				Status:      enum.UserLogStatusInit,
			}
			isAllow, err := userLogUseCase.IsAllowLogin(user.Id)
			if err == nil {
				if isAllow {
					err = userLogUseCase.Add(userLog)
					c.JSON(i, loginResponse)
					return
				} else {
					c.String(401, "limit concurrent user login")
					return
				}
			} else {
				c.String(500, "check failed")
				return
			}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			fmt.Println("IdentityHandler")
			claims := jwt.ExtractClaims(c)
			return &model.Claim{
				UserId: claims[identityKey].(string),
				Role:   claims["role"].(string),
			}
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			fmt.Println("Authorizator")
			if v, ok := data.(*model.Claim); ok && (v.Role == "user" || v.Role == "admin") {
				return true
			}
			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			fmt.Println("Unauthorized")
			switch message {
			case "server error":
				c.JSON(code, gin.H{
					"code":    500,
					"message": message,
				})
			default:
				c.JSON(code, gin.H{
					"code":    code,
					"message": message,
				})
			}
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
		LogoutResponse: func(c *gin.Context, i int) {
			claims := jwt.ExtractClaims(c)
			fmt.Println("----", claims)
			exp := time.Unix(int64(claims["exp"].(float64)), 0).Format(time.RFC3339)
			userLogUseCase.Update(model.UserLog{Status: enum.UserLogStatusInit, ExpiredTime: exp}, model.UserLog{Status: enum.UserLogStatusDone})
			c.Writer.WriteHeader(i)
		},
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

// @tags User
// @Summary user logout
// @Description user logout
// @start_time default
// @Param Authorization header string true "'Bearer ' + token"
// @Success 200 {string} string	""
// @Router /api/v1/auth/logout [post]
func ABC() {

}
