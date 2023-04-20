package api

import (
	"gb-detection/recognition"

"github.com/gin-gonic/gin"

)

type StaffAPI struct {
	db store.StaffDb
}

AddStaffEndPoint

UpdateStaffEndPoint
DeleteStaffEndPoint
GetStaffEndPoint


RecognizeStaffEndPoint

FindStaffEndPoint



func StaffGroup(r *gin.Engine, db store.StaffDb, imageDB store.imageDB, recognizer *recognition.Recognizer, url string) {
	api := NewStaffAPI(db, imageDB, recognizer, url)
	staff := r.Group("/api/staff")

	{
		staff.POST("/add", api.AddStaffEndPoint)
		staff.PUT("/update", api.UpdateStaffEndPoint)
		staff.DELETE("/delete", api.DeleteStaffEndPoint)
		staff.GET("/get", api.GetStaffEndPoint)
		staff.GET("/all", api.GetAllStaffEndPoint)
		staff.POST("/recognize", api.RecognizeStaffEndPoint)
		staff.POST("/find", api.FindStaffEndPoint)
	}
}
