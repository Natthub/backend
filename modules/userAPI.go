package modules

import (
	db "backend/database"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	_ "go.mongodb.org/mongo-driver/mongo"
	_ "go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"strings"
	"time"
)

type User struct {
	UserID string `bson:"userID"`
	Image  string `bson:"image"`
	Name   string `bson:"name"`
	Coupon struct {
		Used   int `bson:"used"`
		Unused int `bson:"unused"`
	} `bson:"coupon"`
	Product []interface{} `bson:"product"`
}
type UserWithID struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	UserID string             `bson:"userID"`
	Image  string             `bson:"image"`
	Name   string             `bson:"name"`
	Coupon struct {
		Used   int `bson:"used"`
		Unused int `bson:"unused"`
	} `json:"coupon"`
	Product []interface{} `bson:"product"`
}

type Receipt struct {
	ReceiptID string `bson:"receiptID"`
	UserID    string `bson:"userID"`
	Image     string `bson:"image"`
	Score     int    `bson:"score"`
}

type ProductWithID struct {
	ID    primitive.ObjectID `bson:"_id,omitempty"`
	Name  string             `bson:"name"`
	Image string             `bson:"image"`
	Score int                `bson:"score"`
}

func IsDup(err error) bool {
	var e mongo.WriteException
	if errors.As(err, &e) {
		for _, we := range e.WriteErrors {
			if we.Code == 11000 {
				return true
			}
		}
	}
	return false
}

func GetUserByID(c *gin.Context) {

	collection := db.Database.Collection("user")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	userID := c.Param("userID")
	cursor, err := collection.Find(ctx, bson.M{"userID": userID})

	if err != nil {
		c.JSON(500, err)
	}
	var users []bson.M
	if err = cursor.All(context.TODO(), &users); err != nil {
		c.JSON(500, err)
	}

	c.JSON(200, users)
}

func Register(c *gin.Context) {

	userID := c.PostForm("userID")
	name := c.PostForm("name")
	image := c.PostForm("image")
	user := User{
		UserID: userID,
		Name:   name,
		Image:  image,
		Coupon: struct {
			Used   int `bson:"used"`
			Unused int `bson:"unused"`
		}(struct {
			Used   int
			Unused int
		}{Used: 0, Unused: 0}),
		Product: []interface{}{},
	}

	collection := db.Database.Collection("user")

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	res, err := collection.InsertOne(ctx, user)
	if err != nil {
		if IsDup(err) {
			c.JSON(http.StatusOK, gin.H{"status": 0, "message": "duplicate user id"})
		} else {
			c.JSON(500, err)
		}
	} else {
		id := res.InsertedID
		c.JSON(http.StatusOK, gin.H{"status": 1, "id": id})
	}

}

func SendReceipt(c *gin.Context) {
	userID := c.PostForm("userID")
	receiptID := c.PostForm("receiptID")
	image, _ := c.FormFile("image")

	if image != nil {
		imageFileName := receiptID + "-" + time.Now().Format("01-02-2006-15:04:05.000000000") + ".jpg"
		imageFileName = strings.ReplaceAll(imageFileName, ":", ".")

		// Upload the file to specific dst.
		c.SaveUploadedFile(image, "./image_receipt/"+imageFileName)

		receipt := Receipt{
			ReceiptID: receiptID,
			UserID:    userID,
			Score:     0,
			Image:     imageFileName,
		}

		collection := db.Database.Collection("receipt")

		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		res, err := collection.InsertOne(ctx, receipt)
		if err != nil {
			fmt.Println(err)
			if IsDup(err) {
				c.JSON(http.StatusOK, gin.H{"status": 0, "message": "duplicate receipt id"})
			} else {
				c.JSON(500, err)
			}

		} else {
			id := res.InsertedID
			c.JSON(http.StatusOK, gin.H{"status": 1, "id": id})
		}
	} else {
		c.JSON(417, "please select file")
	}
}

func Exchange(c *gin.Context) {
	userID := c.PostForm("userID")
	productID := c.PostForm("productID")
	productObjectID, _ := primitive.ObjectIDFromHex(productID)

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	collection := db.Database.Collection("products")
	cursor := collection.FindOne(ctx, bson.M{"_id": productObjectID})
	var product = ProductWithID{}
	cursor.Decode(&product)

	collection = db.Database.Collection("user")
	cursor = collection.FindOne(ctx, bson.M{"userID": userID})
	var user = UserWithID{}
	cursor.Decode(&user)
	if user.Coupon.Unused >= product.Score && userID != "" && productID != "" {
		collection = db.Database.Collection("user")
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		filter := bson.M{"userID": userID}
		update := bson.M{"$inc": bson.M{"coupon.used": product.Score}}

		res, updateErr := collection.UpdateMany(ctx, filter, update)
		if updateErr != nil {
			fmt.Println("4")
			log.Fatal(updateErr)
		}
		_ = res.ModifiedCount

		update = bson.M{"$inc": bson.M{"coupon.unused": -product.Score}}
		res, updateErr = collection.UpdateMany(ctx, filter, update)
		if updateErr != nil {
			fmt.Println("4")
			log.Fatal(updateErr)
		}
		_ = res.ModifiedCount

		res, pushErr := collection.UpdateOne(
			ctx,
			bson.M{"userID": userID},
			bson.M{"$push": bson.M{"product": product}},
		)
		if pushErr != nil {
			c.JSON(500, pushErr)
		} else {
			c.JSON(http.StatusOK, gin.H{"status": 1, "message": "success"})
		}

	} else {
		c.JSON(417, gin.H{"message": "bad data"})
	}

}
