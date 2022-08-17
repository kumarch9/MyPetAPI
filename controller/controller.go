package controller

import (
	cn "mypet/connection"
	"mypet/model"

	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var petCollection *mongo.Collection = cn.GetDBCollection(cn.ConnectionDb(), "golangDB", "mypetCol")
var validate = validator.New()

func PostPet() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		var pet model.Pet
		defer cancel()

		if err := c.BindJSON(&pet); err != nil {
			c.JSON(http.StatusNotImplemented, gin.H{"message": err.Error()})
			return
		}

		//using validator to check becoming data in same pet model fromate
		if validationErr := validate.Struct(&pet); validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "post data is not valid format!"})
			return
		}

		//using the pointer for memory save in newPet
		newPet := &model.Pet{
			Id:           primitive.NewObjectID(),
			GenInfo:      pet.GenInfo,
			ActivityInfo: pet.ActivityInfo,
			FeedInfo:     pet.FeedInfo,
			VetInfo:      pet.VetInfo,
		}

		result, err := petCollection.InsertOne(ctx, newPet)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		var id string = newPet.Id.Hex()
		//c.JSON(http.StatusCreated, gin.H{"data": data})
		c.JSON(http.StatusCreated, bson.D{bson.E{Key: id, Value: result}})

	}
}

func GetPetById() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		pId := c.Param("Id")
		var pet model.Pet
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(pId)

		err := petCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&pet)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusFound, bson.M{"data": pet})
	}
}

func GetPetByName() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		name := c.Param("name")
		var pet model.Pet
		defer cancel()
		err := petCollection.FindOne(ctx, bson.M{"geninfo.name": name}).Decode(&pet) //find in the value in object key
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusFound, bson.M{"data": pet})
	}
}

func EditPetById() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		pId := c.Param("Id")
		var pet model.Pet
		defer cancel()
		objId, _ := primitive.ObjectIDFromHex(pId)
		if err := c.BindJSON(&pet); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		//using validator to check becoming data in same pet model fromate
		if validationErr := validate.Struct(&pet); validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "post data is not valid format!"})
			return
		}
		updateData := bson.M{"geninfo": pet.GenInfo, "activityinfo": pet.ActivityInfo, "feedinfo": pet.FeedInfo, "vetinfo": pet.VetInfo}
		result, err := petCollection.UpdateOne(ctx, bson.M{"_id": objId}, bson.M{"$set": updateData})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		var updatePet model.Pet
		if result.ModifiedCount == 1 {
			if err := petCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&updatePet); err != nil {
				c.JSON(http.StatusNotModified, gin.H{"err": err.Error()})
				return
			}
		}
		c.JSON(http.StatusOK, gin.H{"new document": updatePet})
	}
}

func DeletePetById() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		pId := c.Param("Id")
		defer cancel()
		objId, _ := primitive.ObjectIDFromHex(pId)
		result, err := petCollection.DeleteOne(ctx, bson.M{"_id": objId})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		if result.DeletedCount < 1 {
			c.JSON(http.StatusBadRequest, gin.H{"Message": "deleted document not found of ID_" + pId})
			return
		}
		c.JSON(http.StatusOK, gin.H{"Message": "deleted document successed of ID_" + pId})
	}
}

func GetPets() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var pets []model.Pet
		defer cancel()
		result, errFind := petCollection.Find(ctx, bson.M{})
		if errFind != nil {
			c.JSON(http.StatusNoContent, gin.H{"error": errFind})
			return
		}
		defer result.Close(ctx)
		var onePet model.Pet
		for result.Next(ctx) {
			if err := result.Decode(&onePet); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err})
				return
			}
			pets = append(pets, onePet)
		}
		c.JSON(http.StatusFound, gin.H{"documents": pets})
	}
}
