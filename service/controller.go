package service

// // onboard
// func Submit(c *gin.Context) {}

// // finance
// func CheckAllDbsCountAndID(c *gin.Context) {
// 	faileddbs := fin.Fin.CheckAllDbsCountAndID()
// 	// return failed dbs
// 	if len(faileddbs) > 0 {
// 		msg := strings.Join(faileddbs, "\n")
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"requestID": requestid.Get(c),
// 			"code":      http.StatusInternalServerError,
// 			"message":   msg,
// 		})
// 		return
// 	}

// 	// sucess
// 	c.JSON(http.StatusOK, gin.H{
// 		"code":    http.StatusOK,
// 		"message": "OK",
// 	})
// }

// func CheckAllWorkloadAndAmount(c *gin.Context) {
// 	faileddbs := fin.Fin.CheckAllWorkloadAndAmount()
// 	// return failed dbs
// 	if len(faileddbs) > 0 {
// 		msg := strings.Join(faileddbs, "\n")
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"requestID": requestid.Get(c),
// 			"code":      http.StatusInternalServerError,
// 			"message":   msg,
// 		})
// 		return
// 	}
// 	// sucess
// 	c.JSON(http.StatusOK, gin.H{
// 		"code":    http.StatusOK,
// 		"message": "OK",
// 	})
// }

// func UpdateAllWorkToFin(c *gin.Context) {
// 	faileds := fin.Fin.UpdateAllWorkToFin()
// 	// return failed
// 	if len(faileds) > 0 {
// 		msg := strings.Join(faileds, "\n")
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"requestID": requestid.Get(c),
// 			"code":      http.StatusInternalServerError,
// 			"message":   msg,
// 		})
// 		return
// 	}
// 	// sucess
// 	c.JSON(http.StatusOK, gin.H{
// 		"code":    http.StatusOK,
// 		"message": "OK",
// 	})
// }

// func UpdateAllFinToProgress(c *gin.Context) {
// 	var req UpdateFinParams
// 	err := c.ShouldBindJSON(&req)
// 	if err != nil {
// 		log.Errorln(err)
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"requestID": requestid.Get(c),
// 			"code":      http.StatusInternalServerError,
// 			"message":   err.Error(),
// 		})
// 		return
// 	}

// 	faileds := fin.Fin.UpdateAllFinToProgress(
// 		req.PaymentDateStr,
// 		req.ActualToken, req.ActualPrice,
// 		req.TargetToken, req.TargetPrice)
// 	// return failed
// 	if len(faileds) > 0 {
// 		msg := strings.Join(faileds, "\n")
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"requestID": requestid.Get(c),
// 			"code":      http.StatusInternalServerError,
// 			"message":   msg,
// 		})
// 		return
// 	}
// 	// sucess
// 	c.JSON(http.StatusOK, gin.H{
// 		"code":    http.StatusOK,
// 		"message": "OK",
// 	})

// }

// func UpdateFinToProgress(c *gin.Context) {
// 	var req UpdateFinParams
// 	err := c.ShouldBindJSON(&req)
// 	if err != nil {
// 		log.Errorln(err)
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"requestID": requestid.Get(c),
// 			"code":      http.StatusInternalServerError,
// 			"message":   err.Error(),
// 		})
// 		return
// 	}
// 	if req.FinNid == "" {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"requestID": requestid.Get(c),
// 			"code":      http.StatusInternalServerError,
// 			"message":   "fin nid is nil",
// 		})
// 		return
// 	}
// 	faileds := fin.Fin.UpdateFinToProgress(
// 		req.FinNid, req.PaymentDateStr,
// 		req.ActualToken, req.ActualPrice,
// 		req.TargetToken, req.TargetPrice)
// 	// return failed
// 	if len(faileds) > 0 {
// 		msg := strings.Join(faileds, "\n")
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"requestID": requestid.Get(c),
// 			"code":      http.StatusInternalServerError,
// 			"message":   msg,
// 		})
// 		return
// 	}
// 	// sucess
// 	c.JSON(http.StatusOK, gin.H{
// 		"code":    http.StatusOK,
// 		"message": "OK",
// 	})
// }

// func PayAll(c *gin.Context) {
// 	faileds := fin.Fin.PayAll()
// 	// return failed
// 	if len(faileds) > 0 {
// 		msg := strings.Join(faileds, "\n")
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"requestID": requestid.Get(c),
// 			"code":      http.StatusInternalServerError,
// 			"message":   msg,
// 		})
// 		return
// 	}
// 	// sucess
// 	c.JSON(http.StatusOK, gin.H{
// 		"code":    http.StatusOK,
// 		"message": "OK",
// 	})
// }

// func Pay(c *gin.Context) {
// 	var fnid string
// 	err := c.ShouldBindJSON(&fnid)
// 	if err != nil {
// 		log.Errorln(err)
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"requestID": requestid.Get(c),
// 			"code":      http.StatusInternalServerError,
// 			"message":   err.Error(),
// 		})
// 		return
// 	}

// 	faileds := fin.Fin.Pay(fnid)
// 	// return failed
// 	if len(faileds) > 0 {
// 		msg := strings.Join(faileds, "\n")
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"requestID": requestid.Get(c),
// 			"code":      http.StatusInternalServerError,
// 			"message":   msg,
// 		})
// 		return
// 	}
// 	// sucess
// 	c.JSON(http.StatusOK, gin.H{
// 		"code":    http.StatusOK,
// 		"message": "OK",
// 	})
// }
