package server

import (
	"bet-server/src/bet"
	"crypto/rsa"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
	"github.com/lestrrat-go/jwx/v2/jwk"
)

type Server struct {
	betService *bet.BetService
	secret     string
	jwtIssuer  string
}

func NewServer() Server {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	var uri = os.Getenv("MONGO_URI")
	var database = os.Getenv("MONGO_DATABASE")
	var betService = bet.NewBetService(uri, database)

	return Server{
		betService: &betService,
		secret:     os.Getenv("SECRET"),
		jwtIssuer:  os.Getenv("JWT_ISSUER"),
	}
}

// AuthRequired is a simple middleware to check the session
func (s *Server) AuthRequired(c *gin.Context) {
	var bearerToken = c.Request.Header.Get("Authorization")
	keySet, err := jwk.Fetch(c.Request.Context(), "https://login.microsoftonline.com/common/discovery/v2.0/keys")

	var _, token, found = strings.Cut(bearerToken, "Bearer ")
	println("myToken: " + token)
	if !found {
		c.AbortWithStatus(401)
		return
	}

	parsedtoken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, fmt.Errorf("kid header not found")
		}

		keys, ok := keySet.LookupKeyID(kid)
		if !ok {
			return nil, fmt.Errorf("key %v not found", kid)
		}
		log.Printf("publickey ")
		publickey := &rsa.PublicKey{}
		err = keys.Raw(publickey)
		if err != nil {

			return nil, fmt.Errorf("could not parse pubkey")
		}
		log.Printf("publickey %v", publickey)

		return publickey, nil
	})

	if err != nil {
		println("parese error")
		println(err.Error())
		c.AbortWithStatus(401)
		return
	}

	if claims, ok := parsedtoken.Claims.(jwt.MapClaims); ok && parsedtoken.Valid {
		exp := claims["exp"].(int64)
		now := time.Now().UnixMilli()
		if now > exp {
			c.AbortWithStatus(401)
			return
		}

		if claims["iss"] != s.jwtIssuer {
			c.AbortWithStatus(401)
			return
		}

		c.Set("user", claims["unique_name"])
		c.Next()
	}
	c.Status(401)
	return
}

func (s *Server) Run() {
	r := gin.Default()

	private := r.Group("/")
	private.Use(s.AuthRequired)
	{
		private.GET("/deal", s.getDeals)
		private.POST("/deal", s.createDeal)
		private.PATCH("/deal/:id/open", s.openDeal)
		private.PATCH("/deal/:id/cancel", s.cancelDeal)
		private.GET("/deal/:id", s.getDealController)

		private.GET("/deal/:id/bet", s.getBets)
		private.POST("/deal/:id/bet", s.placeBet)
	}

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})
	r.Run()
}
