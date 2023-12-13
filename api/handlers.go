package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Login(c *gin.Context) {
	var credentials Credentials
	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//var storedCredentials Credentials
	result := credentials.Username == "Jamon" && credentials.Password == "Serrano" // this is a placeholder don't forget to verify the credentials :)
	if !result {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "credentials incorrect"})
		return
	}

	// if storedCredentials.Password != credentials.Password {
	// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "credentials incorrect"})
	// 	return
	// }

	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &Claims{
		Username: credentials.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func validate(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"content": "Congratulations! Your JWT is valid."})
}

func home(c *gin.Context) {
	// Render the HTML page for the home route
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "Home Page",
	})
}
