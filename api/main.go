package api

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

var (
	jwtKey          = []byte("UseASecureKeyHere")
	serverStartTime time.Time
)

const version = "1.0"

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func Init() {
	serverStartTime = time.Now()

	// no deberia hacer falta usar Gin pero por ahora me facilita un poco la vida. "net/http" es mas ligero (duh) y cubre todas las necesidades
	router := gin.Default()

	router.LoadHTMLGlob("web/html/*")
	router.Static("/images", "./web/images")

	endpointsGroup := router.Group("/devtools")

	endpointsGroup.POST("/login", Login)
	endpointsGroup.GET("/status")
	endpointsGroup.GET("/home", home)

	authorized := endpointsGroup.Group("/")
	authorized.Use(authenticate())
	{
		authorized.GET("/validate", validate)
	}

	if err := router.Run(":8080"); err != nil {
		log.Fatal("Error while running server:", err)
	}
}
