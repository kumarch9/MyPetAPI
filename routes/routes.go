package routes

import (
	"mypet/home"
	"mypet/petctr"
	"mypet/signin"
	"mypet/signout"
	"mypet/userctr"

	"github.com/gin-gonic/gin"
)

func PetRoute(ginR *gin.Engine) {

	h_route := ginR.Group("/")
	// en.GET("/home", home.HomePage())
	// en.POST("/signout", signout.SignOutAll())
	h_route.GET("home", home.HomePage())
	h_route.POST("signout", signout.SignOutAll())
	l_route := ginR.Group("/signin")
	l_route.POST("/user/:email/:password", signin.UserLogin())
	l_route.POST("/admin/:email/:password", signin.AdminLogin())

	// prefixed  route for pet related data
	p_route := ginR.Group("/pet")
	p_route.GET("/id/:Id", petctr.GetPetById())
	p_route.GET("/name/:name", petctr.GetPetByName())
	p_route.GET("/all", petctr.GetPets())
	p_route.POST("/post", petctr.PostPet())
	p_route.PUT("/update/:Id", petctr.EditPetById())
	p_route.DELETE("/delete/:Id", petctr.DeletePetById())

	//prefixed route for user related data
	c_route := ginR.Group("/user")
	c_route.POST("/post", userctr.AddUser())
	c_route.PUT("update/:Id", userctr.EditUserById())
	c_route.DELETE("/delete/:Id", userctr.DeleteUserById())

}
