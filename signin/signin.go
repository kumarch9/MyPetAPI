package signin

import (
	"context"
	"mypet/connection"
	"mypet/hashing"
	"mypet/tokens"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type LoginCredential struct {
	Email    string `bson:"email"`
	Password string `bson:"password"`
	Name     string `bson:"name"`
}

var (
	timeVar    = time.Time{}
	timeLayout = "2006-01-02" //date time format for time parse

)

func UserLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		con := connection.ConnectionDb()
		userCol := connection.GetDBCollection(con, "pet_db", "users_info")
		userEmail := c.Param("email")
		userPsw := c.Param("password")
		var user_cred = LoginCredential{}
		dateNil := timeVar.String()
		time_nil, _ := time.Parse(timeLayout, dateNil)
		filterByEmail_date := bson.M{"$and": []bson.M{{"email": userEmail},
			{"deletedAt": bson.M{"$eq": primitive.NewDateTimeFromTime(time_nil)}},
		}}
		user_projection := options.FindOne().SetProjection(bson.D{
			{Key: "_id", Value: 0},
			{Key: "name", Value: 1},
			{Key: "email", Value: 1},
			{Key: "password", Value: 1},
		})
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if (userEmail == "" || userPsw == "") || (userEmail == "" && userPsw == "") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user password or email is empty!"})
			return
		}

		if err := userCol.FindOne(ctx, filterByEmail_date, user_projection).Decode(&user_cred); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorize user can not access!"})
			return
		}

		if pswMatch := hashing.VarifyPassword(user_cred.Password, userPsw); pswMatch {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "password is incorrect!"})
			return
		}

		userToken, Oktoken := tokens.GenerateToken(user_cred.Name, user_cred.Email)
		if !Oktoken {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "signin again!"})
			return
		}

		//log.Println("usertoken:", userToken)
		c.SetSameSite(http.SameSiteLaxMode)
		c.SetCookie("user_token", userToken, 60*10, "", "", false, true) //age => in seconds (60s * 2 ) mean 2 minutes
		c.JSON(http.StatusOK, gin.H{"message": "signin successed"})

	}

}

func AdminLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		con := connection.ConnectionDb()
		adminCol := connection.GetDBCollection(con, "pet_db", "admin_info")
		adm_cred := LoginCredential{}
		adm_projection := options.FindOne().SetProjection(bson.D{
			{Key: "_id", Value: 0},
			{Key: "name", Value: 1},
			{Key: "email", Value: 1},
			{Key: "password", Value: 1},
		})
		//adminVar := admins.Admin{AdminName: "admin", Email: "admin_host@host.com", Password: "adminadmin"}
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		email := c.Param("email")
		password := c.Param("password")

		if (email == "" || password == "") || (email == "" && password == "") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "email or password is empty, try again!"})
			return
		}

		if err := adminCol.FindOne(ctx, bson.M{"email": email}, adm_projection).Decode(&adm_cred); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorize admin can not access!"})
			return
		}
		defer con.Disconnect(ctx)

		if adm_cred.Password != password {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "admin password is incorrect, try again!"})
			return
		}

		adminToken, Oktoken := tokens.GenerateToken(adm_cred.Name, adm_cred.Email)
		if !Oktoken {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "signin again!"})
			return
		}
		c.SetSameSite(http.SameSiteLaxMode)
		c.SetCookie("adm_token", adminToken, 3600*1, "", "", false, true) //age => in seconds (60s * 2 ) mean 2 minutes
		c.JSON(http.StatusOK, gin.H{"message": "signin successed"})
	}

}
