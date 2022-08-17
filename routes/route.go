package routes

import (
	ct "mypet/controller"

	"github.com/gin-gonic/gin"
)

func PetRoute(en *gin.Engine) {

	route := en.Group("/pet") // prefixed  route it will call before  other  route
	route.GET("/id/:Id", ct.GetPetById())
	route.GET("/name/:name", ct.GetPetByName())
	route.GET("/all", ct.GetPets())
	route.POST("/post", ct.PostPet())
	route.PUT("/update/:Id", ct.EditPetById())
	route.DELETE("/delete/:Id", ct.DeletePetById())

}
