package petctr

import (
	"context"
	"log"
	cn "mypet/connection"
	"mypet/logics"
	petmodel "mypet/pets"
	"mypet/tokens"

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
	petColName        string = "pets_info"
	time_nil, _              = time.Parse(timeLayout, timeVar.String())
	bson_obj_nil_time        = primitive.NewDateTimeFromTime(time_nil)
)
var bind_document = bson.M{
	"_id":          1,
	"geninfo":      1,
	"activityinfo": 1,
	"feedinfo":     1,
	"vetinfo":      1,
	"createdAt":    1,
	"updatedAt":    1,
	"deletedAt":    1,
}

func PostPet() gin.HandlerFunc {
	return func(c *gin.Context) {
		con := cn.ConnectionDb()
		petCollection := cn.GetDBCollection(con, petDBName, petColName)

		readCook, _, expiredcook, creator_name, creator_email, errInCook := tokens.ValidateToken(c)
		if errInCook != nil || (expiredcook && readCook != "") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": errInCook.Error()})
			return
		}

		var pet petmodel.Pet
		var licenceNum string
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		//assign creator info data into variables
		newCreatorInfo := petmodel.Creator{
			CreatorName:  creator_name,
			CreatorEmail: creator_email,
		}
		newUpdaterInfo := petmodel.Updater{
			UpdaterName:  "0",
			UpdaterEmail: "0",
		}
		newDeleterInfo := petmodel.Deleter{
			DeleterName:  "0",
			DeleterEmail: "0",
		}

		if err := c.BindJSON(&pet); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		if validationErr := validate.Struct(&pet); validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "data is not well format!"})
			return
		}

		c.Request.Body.Close()
		licenceNum = pet.GenInfo.Licence_Number

		if isExistPet := logics.KeyValueIsExists(petCollection, "geninfo.licence_number", licenceNum); isExistPet {
			c.JSON(http.StatusConflict, gin.H{"error": "pet is already existed!"})
			return
		}

		newPet := &petmodel.Pet{
			Id:           primitive.NewObjectID(),
			GenInfo:      pet.GenInfo,
			ActivityInfo: pet.ActivityInfo,
			FeedInfo:     pet.FeedInfo,
			VetInfo:      pet.VetInfo,
			CreatorInfo:  newCreatorInfo,
			UpdaterInfo:  newUpdaterInfo,
			DeleterInfo:  newDeleterInfo,
			CreatedAt:    time.Now().Local().UTC(),
			UpdatedAt:    timeVar.Local().UTC(),
			DeletedAt:    timeVar.Local().UTC(),
		}

		result, err := petCollection.InsertOne(ctx, newPet)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer con.Disconnect(ctx)
		c.IndentedJSON(http.StatusCreated, gin.H{"successed and object Id_:": result.InsertedID})

	}
}

func GetPetById() gin.HandlerFunc {
	return func(c *gin.Context) {
		con := cn.ConnectionDb()
		petCollection := cn.GetDBCollection(con, petDBName, petColName)

		readCook, _, expiredcook, creator_name, _, errInCook := tokens.ValidateToken(c)
		if errInCook != nil || (expiredcook && readCook != "") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": errInCook.Error()})
			return
		}
		//log.Println("read cookie", readCook)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		pId := c.Param("Id")

		defer cancel()
		c.Request.Body.Close()
		objId, _ := primitive.ObjectIDFromHex(pId)

		//must equal keys in to get document from db by use of petctr.bind_document.
		filterById_date := bson.M{"$and": []bson.M{{"_id": objId}, {"deletedAt": bson.M{"$eq": bson_obj_nil_time}}}}

		defer con.Disconnect(ctx)
		//if admin requests by id of pets and in else only user will get bind data in selected fields from db.
		if creator_name == "admin" {
			pet_admin := petmodel.Pet{}
			err := petCollection.FindOne(ctx, filterById_date).Decode(&pet_admin)
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusFound, gin.H{"data": pet_admin})

		} else {
			pet := petmodel.PetBind_User{}
			err := petCollection.FindOne(ctx, filterById_date, options.FindOne().SetProjection(bind_document)).Decode(&pet)
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusFound, gin.H{"data": pet})
		}

	}
}

func GetPetByName() gin.HandlerFunc {
	return func(c *gin.Context) {
		con := cn.ConnectionDb()
		petCollection := cn.GetDBCollection(con, petDBName, petColName)

		readCook, _, expiredcook, requester_name, _, errInCook := tokens.ValidateToken(c)
		if errInCook != nil || (expiredcook && readCook != "") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": errInCook.Error()})
			return
		}

		//log.Println("read cookie", readCook)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		name := c.Param("name")

		defer cancel()

		filterByName_date := bson.M{"$and": []bson.M{{"geninfo.name": name}, {"deletedAt": bson_obj_nil_time}}}

		//sortByDate := options.Find().SetSort(bson.D{{Key: "createdAt", Value: 1}}) // -1 for desc order and 1 asc
		c.Request.Body.Close()

		//if admin requests by name of pets,so get in array of pets and in else only user will get bind data in
		//selected fields in array from db.
		if requester_name == "admin" {
			var pets_admin []petmodel.Pet
			var singlePet_admin petmodel.Pet

			result, errFind := petCollection.Find(ctx, filterByName_date)
			if errFind != nil {
				c.JSON(http.StatusNoContent, gin.H{"error": errFind.Error()})
				return
			}

			defer result.Close(ctx)

			for result.Next(ctx) {
				if err := result.Decode(&singlePet_admin); err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				pets_admin = append(pets_admin, singlePet_admin)
			}
			c.JSON(http.StatusFound, gin.H{"data": pets_admin})

		} else {
			pets := []petmodel.PetBind_User{}
			singlePet := petmodel.PetBind_User{}

			result, errFind := petCollection.Find(ctx, filterByName_date, options.Find().SetProjection(bind_document))
			if errFind != nil {
				c.JSON(http.StatusNoContent, gin.H{"error": errFind.Error()})
				return
			}
			defer result.Close(ctx)

			for result.Next(ctx) {
				if err := result.Decode(&singlePet); err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				pets = append(pets, singlePet)
			}
			defer con.Disconnect(ctx)
			c.JSON(http.StatusFound, gin.H{"data": pets})
		}

	}
}

func EditPetById() gin.HandlerFunc {
	return func(c *gin.Context) {
		//db connection create instance from connection.GetDBCollection pkg
		con := cn.ConnectionDb()

		//mongo db collection create instance from connection.GetDBCollection pkg
		petCollection := cn.GetDBCollection(con, petDBName, petColName)

		//calling the tokens.ValidateToken(*gin.context) from tokens.ValidateToken
		readCook, _, expiredcook, updater_name, updater_email, errInCook := tokens.ValidateToken(c)
		if errInCook != nil || (expiredcook && readCook != "") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": errInCook.Error()})
			return
		}

		newUpdaterInfo := petmodel.Updater{
			UpdaterName:  updater_name,
			UpdaterEmail: updater_email,
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		pId := c.Param("Id")

		//new model of pet from petmodel.Pet pkg
		var pet petmodel.Pet

		//if any error or exit by user the process request in context, it will be called.
		defer cancel()

		//string Id from string into mongo object id.
		objId, _ := primitive.ObjectIDFromHex(pId)

		if err := c.BindJSON(&pet); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		//using validation method from (validator.Validate).Structpkg
		if validationErr := validate.Struct(&pet); validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "post data is not valid format!"})
			return
		}

		// request body close
		c.Request.Body.Close()

		//select filter for Id and  nill time is equal to db deleted time.
		filterById_date := bson.M{"$and": []bson.M{{"_id": objId}, {"deletedAt": bson_obj_nil_time}}}

		//declare and initialize from pets.Pet on pkg
		updatedData := bson.M{
			"geninfo":      pet.GenInfo,
			"activityinfo": pet.ActivityInfo,
			"feedinfo":     pet.FeedInfo,
			"vetinfo":      pet.VetInfo,
			"updaterinfo":  newUpdaterInfo,
			"updatedAt":    time.Now().Local().UTC(),
		}

		result, err := petCollection.UpdateOne(ctx, filterById_date, bson.M{"$set": updatedData})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "not updated!"})
			return
		}

		var updatePet petmodel.PetBind_User
		if result.ModifiedCount != 1 || result.MatchedCount != 1 {
			c.JSON(http.StatusNotFound, gin.H{"error": "requested data not found or has been deleted!"})
			return
		}

		if err := petCollection.FindOne(ctx, bson.M{"_id": objId}, options.FindOne().SetProjection(bind_document)).Decode(&updatePet); err != nil {
			c.JSON(http.StatusNotModified, gin.H{"error": err.Error()})
			return
		}

		defer con.Disconnect(ctx)
		c.JSON(http.StatusOK, gin.H{"updated document": &updatePet})
	}

}
func DeletePetById() gin.HandlerFunc {
	return func(c *gin.Context) {
		con := cn.ConnectionDb()
		petCollection := cn.GetDBCollection(con, petDBName, petColName)

		readCook, _, expiredcook, deleter_name, deleter_email, errInCook := tokens.ValidateToken(c)
		if errInCook != nil || (expiredcook && readCook != "") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": errInCook.Error()})
			return
		}

		//set deleter info into variable
		newDeleterInfo := petmodel.Deleter{
			DeleterName:  deleter_name,
			DeleterEmail: deleter_email,
		}

		//log.Println("read cookie", readCook)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		pId := c.Param("Id")

		defer cancel()
		c.Request.Body.Close()
		objId, _ := primitive.ObjectIDFromHex(pId)

		filterById_date := bson.M{"$and": []bson.M{{"_id": objId}, {"deletedAt": bson.M{"$eq": bson_obj_nil_time}}}}
		deletedData := bson.M{"deletedAt": time.Now().Local().UTC(), "deleterinfo": newDeleterInfo}

		result, err := petCollection.UpdateOne(ctx, filterById_date, bson.M{"$set": deletedData})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		if result.ModifiedCount != 1 {
			c.JSON(http.StatusNotModified, gin.H{"err": "count not deleted, record might not found or already deleted!"})
			return
		}
		defer con.Disconnect(ctx)
		c.JSON(http.StatusOK, gin.H{"message": "document is deleted of ID_" + pId})
	}
}

func GetPets() gin.HandlerFunc {
	return func(c *gin.Context) {
		con := cn.ConnectionDb()
		petCollection := cn.GetDBCollection(con, petDBName, petColName)

		readCook, _, expiredcook, creator_name, _, errInCook := tokens.ValidateToken(c)
		if errInCook != nil || (expiredcook && readCook != "") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": errInCook.Error()})
			return
		}

		// log.Println("readCook", readCook)
		// log.Println("creator_name", creator_name)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		filterDate := bson.M{"deletedAt": bson_obj_nil_time}

		sortByDate := options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}}) // -1 for desc order and 1 asc

		if creator_name == "admin" {
			var pets_admin = []petmodel.Pet{}
			var pet_admin petmodel.Pet
			result, errFind := petCollection.Find(ctx, bson.D{{}}, sortByDate)
			if errFind != nil {
				c.JSON(http.StatusNoContent, gin.H{"error": errFind.Error()})
				return
			}
			defer result.Close(ctx)

			for result.Next(ctx) {
				if err := result.Decode(&pet_admin); err != nil {
					log.Println("result next error :", err.Error())
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				pets_admin = append(pets_admin, pet_admin)
			}
			c.JSON(http.StatusFound, gin.H{"documents": pets_admin})
		} else {
			var pets = []petmodel.PetBind_User{}
			var pet petmodel.PetBind_User
			result, errFind := petCollection.Find(ctx, filterDate, options.Find().SetProjection(bind_document))
			if errFind != nil {
				c.JSON(http.StatusNoContent, gin.H{"error": errFind.Error()})
				return
			}
			defer result.Close(ctx)

			for result.Next(ctx) {
				if err := result.Decode(&pet); err != nil {
					log.Println("result next error :", err.Error())
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				pets = append(pets, pet)
			}
			defer con.Disconnect(ctx)
			c.JSON(http.StatusFound, gin.H{"documents": pets})
		}

	}
}
