package main

import (
	"backend/modules"
	"github.com/gin-gonic/gin"
)


func main() {

	r := gin.Default()
	r.GET("/user/:userID", modules.GetUserByID)

	r.POST("/user/register", modules.Register)

	r.POST("/user/receipt", modules.SendReceipt)

	r.GET("/receipt/nonscore", modules.GetReceiptNonScore)

	r.POST("/receipt/addscore", modules.AddScore)

	r.POST("/product/create", modules.CreateProduct)

	r.GET("/product/all", modules.GetAllProduct)

	r.POST("/user/exchange", modules.Exchange)



	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
