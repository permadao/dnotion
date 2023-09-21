package service

import "github.com/gin-gonic/gin"

func InitRouter(r *gin.Engine) {

	group := r.Group("/v1/")
	// onboard
	group.POST("/onboard/submit", Submit)

	// finance
	group.GET("/finance/checkid", CheckAllDbsCountAndID)
	group.GET("/finance/checkAmount", CheckAllWorkloadAndAmount)
	group.POST("/finance/updatework", UpdateAllWorkToFin)
	group.POST("/finance/updateallfin", UpdateAllFinToProgress)
	group.POST("/finance/updatefin", UpdateFinToProgress)
	group.POST("/finance/payall", PayAll)
	group.POST("/finance/pay", Pay)
}
