package api

import "github.com/gin-gonic/gin"

func StaffGroup(r *gin.Engine, url string) {
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
