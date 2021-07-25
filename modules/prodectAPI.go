package modules

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"strings"
	"time"
)


type Product struct {
	Name  string `bson:"name"`
	Image string `bson:"image"`
	Score int `bson:"score"`
}

func CreateProduct(c *gin.Context) {
	name := c.PostForm("name")
	image, _ := c.FormFile("image")
	score := c.PostForm("score")

	scoreInt,_ := StrToInt(score)


	if name!="" && scoreInt > 0{
		imageFileName := name+"-"+time.Now().Format("01-02-2006-15:04:05.000000000")+".jpg"
		imageFileName = strings.ReplaceAll(imageFileName,":",".")

		// Upload the file to specific dst.
		c.SaveUploadedFile(image, "./image_product/"+imageFileName)

		product := Product{
			Name:  name,
			Image: imageFileName,
			Score: scoreInt,
		}

		collection := db.Collection("products")

		ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)
		res, err := collection.InsertOne(ctx, product)
		if err != nil {
			c.JSON(500, err)
		} else {
			id := res.InsertedID
			c.JSON(200, id)
		}
	}else{
		c.JSON(417, "Enter name and Score more than 0")
	}
}

func GetAllProduct(c *gin.Context) {
	collection := db.Collection("products")
	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(500, err)
	}
	var products []bson.M
	if err = cursor.All(ctx, &products); err != nil {
		c.JSON(500, err)
	}

	c.JSON(200, products)
}

