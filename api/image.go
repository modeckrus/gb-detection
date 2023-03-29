package api

import "github.com/gin-gonic/gin"

type ImageAPI struct{}

func NewImageApi() ImageAPI {
	return ImageAPI{}
}

func ImageGroup(r *gin.Engine, url string) {
	api := NewImageApi()
	group := r.Group("/api/images")
	{
		group.POST("/upload", api.UploadEndPoint)
		group.GET("/data", api.DataEndPoint)
		group.GET("/file", api.FileEndPoint)
		group.POST("/recognize", api.RecognizeEndPoint)
		group.POST("/save", api.SaveEndPoint)
	}

}
