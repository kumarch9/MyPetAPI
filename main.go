package main

import (
	"context"
	con "mypet/connection"
	rt "mypet/routes"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	router *gin.Engine
	//ctx        = context.TODO()
	ctxTime, _ = context.WithTimeout(context.Background(), 10*time.Second)
)

func main() {
	router = gin.Default()
	conClient := con.ConnectionDb()
	defer conClient.Disconnect(ctxTime)
	rt.PetRoute(router)
	router.Run(":8055")
}
