package main

import (
	"backend/database"
	"backend/modules"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"os"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST,HEAD,PATCH, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func getPort() string {
	var port = os.Getenv("PORT")
	if port == "" {
		port = "8080"
		fmt.Println("No Port In Heroku" + port)
	}
	return ":" + port
}

func main() {

	r := gin.Default()

	r.Use(cors.Default())

	database.ConnectDB()

	r.GET("/user/:userID", modules.GetUserByID)

	r.POST("/user/register", modules.Register)

	r.POST("/user/receipt", modules.SendReceipt)

	r.GET("/receipt/nonscore", modules.GetReceiptNonScore)

	r.POST("/receipt/addscore", modules.AddScore)

	r.POST("/product/create", modules.CreateProduct)

	r.GET("/product/all", modules.GetAllProduct)

	r.POST("/user/exchange", modules.Exchange)

	r.GET("/image/product/:filename", modules.GetProductImage)

	r.GET("/image/receipt/:filename", modules.GetReceiptImage)

	r.Run(getPort()) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
