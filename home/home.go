package home

import (
	"context"
	"mypet/connection"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// const mgs string = "Welcome in my pet site....."

var petDBName string = "pet_db"
var petColName string = "pets_info"
var userColName string = "users_info"
var errInFound []string
var NumberOfFound map[string]string

func HomePage() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		con := connection.ConnectionDb()

		c.JSON(http.StatusOK, gin.H{"msg": "hello in home page"})
		errInFound = nil
		colUser := connection.GetDBCollection(con, petDBName, userColName)
		var colPet = connection.GetDBCollection(con, petDBName, petColName)
		totalUser, errInUser := TotalDocument(colUser)
		if errInUser != nil {
			errInFound = append(errInFound, errInUser.Error())
		}

		//here we are using only three category of pet(dog,fox,cat) if any more need, so add here code
		totalPet, errInPet := TotalDocument(colPet)
		if errInPet != nil {
			errInFound = append(errInFound, errInPet.Error())
		}

		totalDog, errInDog := CountByKeyValue(colPet, "geninfo.pettype", "dog")
		if errInDog != nil {
			errInFound = append(errInFound, errInDog.Error())
		}

		totalCat, errInCat := CountByKeyValue(colPet, "geninfo.pettype", "cat")
		if errInCat != nil {
			errInFound = append(errInFound, errInCat.Error())
		}
		totalFox, errInFox := CountByKeyValue(colPet, "geninfo.pettype", "fox")
		if errInFox != nil {
			errInFound = append(errInFound, errInFox.Error())
		}

		defer con.Disconnect(ctx)
		newErrArray := errInFound
		// rrr 500 instead 200
		c.JSON(http.StatusOK, gin.H{
			"totalUser": totalUser,
			"totalPet":  totalPet,
			"totalDog":  totalDog,
			"totalCat":  totalCat,
			"totalFox":  totalFox,
			"error":     newErrArray[0:],
		})

	}

}

func TotalDocument(dbCollectionName *mongo.Collection) (numPets string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	count, errInDB := dbCollectionName.CountDocuments(ctx, bson.M{})
	//defer con.Disconnect(ctx)
	if errInDB != nil {
		return "", errInDB
	}

	var countStr string = strconv.FormatInt(count, 10)
	return countStr, nil

}

func CountByKeyValue(dbCollectionName *mongo.Collection, keyName string, valueName string) (numPets string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	count, errInDB := dbCollectionName.CountDocuments(ctx, bson.M{keyName: valueName})
	//defer con.Disconnect(ctx)
	if errInDB != nil {
		return "", errInDB
	}

	var countStr string = strconv.FormatInt(count, 10)
	return countStr, nil
}
