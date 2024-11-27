package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/exp/rand"
	"poc-goland.com/poc/src/entities"
	"poc-goland.com/poc/src/middleware"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func sumNumbers(context *gin.Context) {
	firstNbr, err1 := strconv.Atoi(context.Param("numberFrst"))
	sndNbr, err2 := strconv.Atoi(context.Param("numberSnd"))

	if err1 != nil || err2 != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid numbers",
		})
		return
	}

	sum := firstNbr + sndNbr

	context.JSON(http.StatusOK, gin.H{
		"Result of sum": sum,
	})
}

func helloServer(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{
		"Message": "HI YOU!",
	})
}

func randomNumber(context *gin.Context) {
	maximunNbr, err := strconv.Atoi(context.Param("maxRand"))

	if err != nil || maximunNbr <= 0 {
		context.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid number",
		})
		return
	}
	randNbr := rand.Intn(maximunNbr) + 1

	context.JSON(http.StatusOK, gin.H{
		"Random Number": randNbr,
	})
}

func PublicRoutes(router *gin.Engine, databaseContext entities.Database) {
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/sum/:numberFrst/:numberSnd", sumNumbers)
	router.GET("/random/:maxRand", randomNumber)
	router.POST("/login", databaseContext.LoginAuthentication)
	router.POST("/register", databaseContext.CreateNewUser)
}

func PrivateRoutes(router *gin.Engine, databaseContext entities.Database) {
	router.GET("/hi", middleware.VerifyJWTCookie, helloServer)
	router.PUT("/user/modify-password", middleware.VerifyJWTCookie, databaseContext.ModifyPassword)
	router.DELETE("/user/delete-account", middleware.VerifyJWTCookie, databaseContext.DeleteAccount)
}
