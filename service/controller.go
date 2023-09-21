package service

import (
	"net/http"
	"strings"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/permadao/dnotion/fin"
)

// onboard
func Submit(c *gin.Context) {}

// finance
func CheckAllDbsCountAndID(c *gin.Context) {
	faileddbs := fin.Fin.CheckAllDbsCountAndID()
	// return failed dbs
	if len(faileddbs) > 0 {
		msg := strings.Join(faileddbs, "\n")
		c.JSON(http.StatusInternalServerError, gin.H{
			"requestID": requestid.Get(c),
			"code":      http.StatusInternalServerError,
			"message":   msg,
		})
	}

	// sucess
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "All db count and id OK",
	})
}

func CheckAllWorkloadAndAcutal(c *gin.Context) {}

func UpdateAllWorkToFin(c *gin.Context) {}

func UpdateAllFinToProgress(c *gin.Context) {}

func UpdateFinToProgress(c *gin.Context) {}

func PayAll(c *gin.Context) {}

func Pay(c *gin.Context) {}
