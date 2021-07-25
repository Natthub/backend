package modules

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strconv"
	"strings"
	"time"
)


func GetReceiptNonScore(c *gin.Context) {
	collection := db.Collection("receipt")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	cursor, err := collection.Find(ctx, bson.M{"score": bson.M{"$eq" : 0}} )

	if err != nil {
		c.JSON(500, err)
	}
	var receipt []bson.M
	if err = cursor.All(ctx, &receipt); err != nil {
		c.JSON(500, err)
	}

	c.JSON(200, receipt)
}

func StrToInt(str string) (int, error) {
	nonFractionalPart := strings.Split(str, ".")
	return strconv.Atoi(nonFractionalPart[0])
}

func AddScore(c *gin.Context) {
	collection := db.Collection("receipt")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	userID := c.PostForm("userID")
	receiptID := c.PostForm("receiptID")
	score := c.PostForm("score")

	scoreInt, _ := StrToInt(score)

	if scoreInt > 0{
		res, err := collection.UpdateOne(
			ctx,
			bson.M{"receiptID": receiptID},
			bson.D{
				{"$set", bson.D{{"score", scoreInt}}},
			}, )
		if err != nil {
			c.JSON(500, err)
		}

		collection = db.Collection("user")
		ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
		res, err = collection.UpdateOne(
			ctx,
			bson.M{"userID": userID},
			bson.D{
				{"$inc", bson.D{{"coupon.unused", scoreInt}}},
			}, options.Update().SetUpsert(true))


		if err != nil {
			c.JSON(500, err)
		} else{
			id := res.ModifiedCount
			c.JSON(200, id)
		}
	}else{
		c.JSON(417, "Score must be more than 0")
	}
}