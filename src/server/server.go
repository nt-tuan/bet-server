package server

import (
	"bet-server/src/auth"
	"bet-server/src/bet"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type Server struct {
	betService  *bet.BetService
	authService *auth.AuthService
}

func NewServer() Server {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	var uri = os.Getenv("MONGO_URI")
	var database = os.Getenv("MONGO_DATABASE")
	var betService = bet.NewBetService(uri, database)
	var authService = auth.NewAuthServive("https://login.microsoftonline.com/b035d9fd-a388-48f2-9abf-cea8457262ce/v2.0", os.Getenv("JWT_SECRET"))

	return Server{
		betService:  &betService,
		authService: &authService,
	}
}

// AuthRequired is a simple middleware to check the session
func (s *Server) AuthRequired(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")
	username, err := s.authService.ValidateAccessToken(token)
	if err != nil {
		c.AbortWithStatus(401)
		return
	}
	c.Set("user", username)
}

func (s *Server) Run() {
	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowCredentials = true
	config.AllowHeaders = []string{"Authorization"}

	r.Use(cors.New(config))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})
	r.POST("/login", s.login)

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
		private.GET("/highlight-deal", s.getHighlightDeal)
	}

	r.Run()
}
