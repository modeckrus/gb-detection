package api

func ThirdpartyGroup(r *gin.Engine url string) {
	api := NewThirdpartyAPI(db, st)
	group := r.Group("/api/thirdparty")
	{
		group.GET("/timerecordStream", api.TimerecordStream)
		group.POST("/add", api.AddThirdpartyEndPoint)
		group.DELETE("/delete", api.DeleteThirdpartyEndPoint)
		group.GET("/all", api.AllThirdpartyEndPoint)
		group.POST("/check", api.ChekThirdpartyEndPoint)
	}
	go api.RegisterStream()
}
