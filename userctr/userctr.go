package userctr

import (
	"context"
	"fmt"
	"mypet/connection"
	"mypet/hashing"
	"mypet/logics"
	"mypet/tokens"
	"mypet/users"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	validate                 = validator.New()
	timeVar                  = time.Time{}
	timeLayout               = "2006-01-02" //date time format for time parse
	petDBName         string = "pet_db"
	userColName       string = "users_info"
	time_nil, _              = time.Parse(timeLayout, timeVar.String())
	bson_obj_nil_time        = primitive.NewDateTimeFromTime(time_nil)
)

type UserFoundID struct {
	Id primitive.ObjectID `bson:"_id"`
}

func AddUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		con := connection.ConnectionDb()
		UserCollection := connection.GetDBCollection(con, petDBName, userColName)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		var userVar = users.User{}
		//var varMobile, varEmail, pswToHash string

		if err := c.BindJSON(&userVar); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if validationErr := validate.Struct(&userVar); validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "data is not well format!"})
			return
		}

		if userVar.Name == "admin" || userVar.Name == "administrator" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "reserved name!"})
			return
		}
		// varMobile = userVar.MobileNum
		// varEmail = userVar.Email
		// pswToHash = userVar.Password

		if foundMobile := logics.KeyValueIsExists(UserCollection, "mobilenum", userVar.MobileNum); foundMobile {
			c.JSON(http.StatusConflict, gin.H{"error": "user's mobile  is already existed!"})
			return
		}

		if foundMail := logics.KeyValueIsExists(UserCollection, "email", userVar.Email); foundMail {
			c.JSON(http.StatusConflict, gin.H{"error": "user's email  is already existed!"})
			return
		}

		c.Request.Body.Close()

		hashPsw, errHash := hashing.CreateHash(userVar.Password)
		if errHash != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": errHash.Error()})
			return
		}

		newUser := &users.User{
			Id:           primitive.NewObjectID(),
			Name:         userVar.Name,
			Email:        userVar.Email,
			MobileNum:    userVar.MobileNum,
			Organisation: userVar.Organisation,
			Designation:  userVar.Designation,
			Address:      userVar.Address,
			AadhaarNum:   userVar.AadhaarNum,
			Password:     hashPsw,
			CreatedAt:    time.Now().Local().UTC(),
			UpdatedAt:    timeVar.Local().UTC(),
			DeletedAt:    timeVar.Local().UTC(),
		}

		result, errInsert := UserCollection.InsertOne(ctx, newUser)

		if errInsert != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"err": "values could not inserted"})
			return
		}
		defer con.Disconnect(ctx)
		c.JSON(http.StatusCreated, gin.H{"succeeded": result.InsertedID})

	}
}

func GetUserById() gin.HandlerFunc {
	return func(c *gin.Context) {
		con := connection.ConnectionDb()
		UserCollection := connection.GetDBCollection(con, petDBName, userColName)

		readCook, _, expiredcook, admin_name, _, errInCook := tokens.ValidateToken(c)
		if errInCook != nil || (expiredcook && readCook != "") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": errInCook.Error()})
			return
		}

		if admin_name != "admin" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized  user can not access!"})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		pId := c.Param("Id")
		var userVar = users.User{}
		c.Request.Body.Close()

		objId, _ := primitive.ObjectIDFromHex(pId)

		filterById_date := bson.M{"$and": []bson.M{{"_id": objId}, {"deletedAt": bson_obj_nil_time}}}

		if err := UserCollection.FindOne(ctx, filterById_date).Decode(&userVar); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		defer con.Disconnect(ctx)
		c.JSON(http.StatusFound, gin.H{"data": userVar})
	}
}

func EditUserById() gin.HandlerFunc {
	return func(c *gin.Context) {
		con := connection.ConnectionDb()
		UserCollection := connection.GetDBCollection(con, petDBName, userColName)

		readCook, _, expiredcook, _, updater_email, errInCook := tokens.ValidateToken(c)
		if errInCook != nil || (expiredcook && readCook != "") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": errInCook.Error()})
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		pId := c.Param("Id")

		var userVar users.User

		objId, _ := primitive.ObjectIDFromHex(pId)

		if err := c.BindJSON(&userVar); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if validationErr := validate.Struct(&userVar); validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "post data is not valid format!"})
			return
		}

		if userVar.Name == "admin" || userVar.Name == "administrator" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "reserved name!"})
			return
		}

		c.Request.Body.Close()
		userFoundId := UserFoundID{}
		errFind := UserCollection.FindOne(ctx, bson.M{"email": updater_email}, options.FindOne().SetProjection(bson.M{"_id": 1})).Decode(&userFoundId)
		if errFind != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": errFind.Error()})
			return
		}

		if userFoundId.Id != objId {
			c.JSON(http.StatusConflict, gin.H{"error": "updater can use only own data!"})
			return
		}
		filterById_date := bson.M{"$and": []bson.M{{"_id": objId}, {"deletedAt": bson_obj_nil_time}}}
		updatedData := bson.M{
			"name":         userVar.Name,
			"email":        userVar.Email,
			"mobilenum":    userVar.MobileNum,
			"organisation": userVar.Organisation,
			"designation":  userVar.Designation,
			"address":      userVar.Address,
			"aadhaarnum":   userVar.AadhaarNum,
			"password":     userVar.Password,
			"updatedAt":    time.Now().Local().UTC()}

		result, err := UserCollection.UpdateOne(ctx, filterById_date, bson.M{"$set": updatedData})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var updateUser users.User
		if result.ModifiedCount == 1 {
			if err := UserCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&updateUser); err != nil {
				c.JSON(http.StatusNotModified, gin.H{"error": err.Error()})
				return
			}
		}
		defer con.Disconnect(ctx)
		c.JSON(http.StatusOK, gin.H{"updated document": updateUser})
	}
}

func DeleteUserById() gin.HandlerFunc {
	return func(c *gin.Context) {
		con := connection.ConnectionDb()
		UserCollection := connection.GetDBCollection(con, petDBName, userColName)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		readCook, _, expiredcook, deleter_name, deleter_email, errInCook := tokens.ValidateToken(c)
		if errInCook != nil || (expiredcook && readCook != "") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": errInCook.Error()})
			return
		}
		pId := c.Param("Id")
		dateNil := timeVar.String()
		time_nil, _ := time.Parse(timeLayout, dateNil)

		c.Request.Body.Close()

		objId, _ := primitive.ObjectIDFromHex(pId)
		userFoundId := UserFoundID{}

		if deleter_name != "admin" {
			errFind := UserCollection.FindOne(ctx, bson.M{"email": deleter_email}, options.FindOne().SetProjection(bson.M{"_id": 1})).Decode(&userFoundId)
			if errFind != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": errFind.Error()})
				return
			}

			fmt.Println()
			if userFoundId.Id != objId {
				c.JSON(http.StatusConflict, gin.H{"error": "updater can use only own data!"})
				return
			}
		}

		filterById_date := bson.M{"$and": []bson.M{{"_id": objId}, {"deletedAt": bson.M{"$eq": primitive.NewDateTimeFromTime(time_nil)}}}}
		result, err := UserCollection.UpdateOne(ctx, filterById_date, bson.M{"$set": bson.M{"deletedAt": time.Now().Local().UTC()}})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var updateUser users.User
		if result.ModifiedCount == 1 {
			if err := UserCollection.FindOne(ctx, filterById_date).Decode(&updateUser); err != nil {
				c.JSON(http.StatusNotModified, gin.H{"err": err.Error()})
				return
			}
		}
		defer con.Disconnect(ctx)
		c.JSON(http.StatusOK, gin.H{"message": "document is deleted of ID_" + pId})
	}
}

func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		con := connection.ConnectionDb()
		UserCollection := connection.GetDBCollection(con, petDBName, userColName)
		readCook, _, expiredcook, user_name, _, errInCook := tokens.ValidateToken(c)
		if errInCook != nil || (expiredcook && readCook != "") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": errInCook.Error()})
			return
		}
		if user_name != "admin" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "accessible only by admin!"})
		}
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var usersVar []users.User
		defer cancel()
		dateNil := timeVar.String()
		time_nil, _ := time.Parse(timeLayout, dateNil)

		filterDate := bson.M{"deletedAt": bson.M{"$eq": primitive.NewDateTimeFromTime(time_nil)}}
		sortByDate := options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}}) // -1 for desc order and 1 asc
		result, errFind := UserCollection.Find(ctx, filterDate, sortByDate)
		if errFind != nil {
			c.JSON(http.StatusNoContent, gin.H{"error": errFind.Error()})
			return
		}

		defer result.Close(ctx)
		var singleUser users.User
		for result.Next(ctx) {
			if err := result.Decode(&singleUser); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			usersVar = append(usersVar, singleUser)
		}
		defer con.Disconnect(ctx)
		c.JSON(http.StatusFound, gin.H{"documents": usersVar})
	}
}
