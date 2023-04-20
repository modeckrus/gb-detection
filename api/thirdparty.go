package api

import (
	"bytes"
	"context"
	"encoding/json"
	"gb-detection/model"
	"gb-detection/store"

	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ThirdpartyAPI struct {
	db
	st
}

func NewThirdpartyAPI() {

}

func (a *ThirdpartyAPI) TimerecordStream(c *gin.Context) {

}
func (a *ThirdpartyAPI) AddThirdpartyEndPoint(c *gin.Context) {
	var timeRecord model.Thirdparty
	err := c.BindJSON(&timeRecord)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	timerecord, err := a.st.AddUrl(timeRecord.Url)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, timerecord)
}

func (a *ThirdpartyAPI) DeleteThirdpartyEndPoint(c *gin.Context) {
	var timeRecord model.Thirdparty
	err := c.BindJSON(&timeRecord)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	err = a.st.DeleteUrl(timeRecord.Url)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, timerecord)
}

func (a *ThirdpartyAPI) AllThirdpartyEndPoint(c *gin.Context) {
	timerecord, err := a.st.All()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, timerecord)
}

func (a *ThirdpartyAPI) RegisterStream() {
	ch, err := a.db.Stream(context.Background())
	if err != nil {
		panic(err)
	}
	for {
		thirdparty := <-ch
		js, err := json.Marshal(thirdparty)
		if err != nil {
			log.Println(err)
			continue
		}
		go func() {
			urls, err := a.st.All()
			if err != nil {
				log.Println(err)
				return
			}
			for _, url := range urls {
				url := url
				go func() {
					buf := bytes.NewBuffer(js)
					resp, err := http.Post(url, "application/json", buf)
					if err != nil {
						panic(err)
					}
					defer resp.Body.Close()
				}()
			}

		}()

	}
}

func (a *ThirdpartyAPI) ChekThirdpartyEndPoint(c *gin.Context) {
	var input interface{}
	err := c.BindJSON(&input)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	log.Println("Check thirdparty")
	bytes, err := json.Marshal(input)
	if err != nil {
		log.Println(err)
	}
	log.Println(string(bytes))
	c.JSON(200, input)
}
func ThirdpartyGroup(r *gin.Engine, db store.TimeRecordDb, st store.ThirdpartyDb, url string) {
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
