package api

import (
	"context"
	//"gb-detection/image"
	"gb-detection/model"
	"gb-detection/recognition"
	"gb-detection/store"

	"io"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type ImageAPI struct {
	db         store.ImageDb
	recognizer *recognition.Recognizer
	url        string
}

func NewImageApi(db store.ImageDb, recognizer *recognition.Recognizer, url string) *ImageAPI {
	return &ImageAPI{db, recognizer, url}
}

func (api *ImageAPI) UploadEndPoint(c *gin.Context) {
	buf := c.Request.Body
	jpegBuf, err := api.recognizer.JpegFormatter(buf)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	image, err := api.db.Add(jpegBuf)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	descriptor, err := api.recognizer.RecognizeSingleFile(image.Path)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	image.Descriptor = descriptor
	err = api.db.Update(image)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, image.ToImagePrintable(api.url))
}

func (api *ImageAPI) DataEndPoint(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	image, err := api.db.Get(id)
	if err != nil {
		c.JSON(404, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, image.ToImagePrintable(api.url))
}

func (api *ImageAPI) FileEndPoint(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	size, file, err := api.db.File(id)
	if err != nil {
		c.JSON(404, gin.H{"error": err.Error()})
		return
	}
	c.DataFromReader(200, size, "image/jpeg", file, nil)
}

func (api *ImageAPI) recognize(ctx context.Context, buf io.ReadCloser) (model.Image, error) {
	descriptor, err := api.recognizer.RecognizeBuf(buf)
	if err != nil {
		return model.Image{}, err
	}
	images, err := api.recognizer.Classify(api.db, *descriptor)
	if err != nil {
		return model.Image{}, err
	}
	return images, nil
}

func (api *ImageAPI) RecognizeEndPoint(c *gin.Context) {
	ctx, done := context.WithTimeout(context.Background(), 2*time.Second)
	defer done()
	buf := c.Request.Body
	resultChan := make(chan model.Image)
	errorChan := make(chan error)
	go func() {
		image, err := api.recognize(ctx, buf)
		if err != nil {
			errorChan <- err
		}
		resultChan <- image
	}()
	select {
	case result := <-resultChan:
		c.JSON(200, result.ToImagePrintable(api.url))
	case err := <-errorChan:
		c.JSON(500, gin.H{"error": err.Error()})
	case <-ctx.Done():
		c.JSON(408, gin.H{"error": "timeout"})
	}

}

func (api *ImageAPI) SaveEndPoint(c *gin.Context) {
	buf := c.Request.Body
	jpegBuf, err := api.recognizer.JpegFormatter(buf)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, image.ToImagePrintable(api.url))
}

func ImageGroup(r *gin.Engine, db store.ImageDb, recognizer *recognition.Recognizer, url string) {
	api := NewImageApi(db, recognizer, url)
	group := r.Group("/api/images")
	{
		group.POST("/upload", api.UploadEndPoint)
		group.GET("/data", api.DataEndPoint)
		group.GET("/file", api.FileEndPoint)
		group.POST("/recognize", api.RecognizeEndPoint)
		group.POST("/save", api.SaveEndPoint)
	}

}
