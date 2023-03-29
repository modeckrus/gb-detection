package api

import (
	"github.com/gin-gonic/gin"
)

func TimeRecordGroup(r *gin.Engine, url string) {
	api := NewTimeRecordAPI(db)
	group := r.Group("/api/timerecord")
	{
		group.POST("/add", api.AddTimeRecordEndPoint)
		group.PUT("/update", api.UpdateTimeRecordEndPoint)
		group.DELETE("/delete", api.DeleteTimeRecordEndPoint)
		group.GET("/get", api.GetTimeRecordEndPoint)
		group.GET("/all", api.AllTimeRecordsEndPoint)
		group.GET("/byemployee", api.TimeRecordsByEmployeeEndPoint)
		group.POST("/bydate", api.TimeRecordsByDateEndPoint)
		group.GET("/lastbyemployee", api.TimeRecordLastByEmployeeEndPoint)
	}

}
