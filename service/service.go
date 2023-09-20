package service

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/permadao/dnotion/db"
	"github.com/permadao/dnotion/fin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	DB  *db.NotionDB
	Fin *fin.Finance
)

func StartServe(ndb *db.NotionDB, fin *fin.Finance) {
	log.Info("Starting server...")

	// router
	router := gin.Default()
	router.Use(gin.Recovery())
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// routes
	InitRouter(router)

	port := fmt.Sprintf(":%s", viper.GetString("service.port"))
	err := router.Run(port)
	if err != nil {
		log.Fatal("run service error: ", err.Error())
	}

	DB = ndb
	Fin = fin
}
