package api

import (
	"encoding/json"
	"gb-detection/model"
	"gb-detection/store"
	"io/ioutil"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TimeRecordAPI struct {
	db store.TimeRecordDb
}

func NewTimeRecordAPI(db store.TimeRecordDb) *TimeRecordAPI {
	return &TimeRecordAPI{db}
}

func (a *TimeRecordAPI) AddTimeRecordEndPoint(c *gin.Context) {
	var addTime model.AddTimeRecord
	buf, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
	}
	err = json.Unmarshal(buf, &addTime)
	if err != nil {
		log.Fatalf("unmarshal error: %s\n", err.Error())
	}
	timeRecord, err := a.db.Add(addTime)
	if err != nil {
	}
	c.JSON(200, timeRecord)
}

func (a *TimeRecordAPI) UpdateTimeRecordEndPoint(c *gin.Context) {
	var updateTime model.UpdateTimeRecord
	buf, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})

	}
	err = json.Unmarshal(buf, &updateTime)
	if err != nil {
		log.Fatalf("unmarshal error: %s\n", err.Error())
	}
	timeRecord, err := a.db.Update(updateTime)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})

	}
	c.JSON(200, timeRecord)

}

func (a *TimeRecordAPI) DeleteTimeRecordEndPoint(c *gin.Context) {
	idReq := c.Query("id")
	id, err := strconv.Atoi(idReq)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	errorD := a.db.Delete(id)
	if errorD != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(204, gin.H{})
}

func (a *TimeRecordAPI) GetTimeRecordEndPoint(c *gin.Context) {
	idReq := c.Query("id")
	id, err := strconv.Atoi(idReq)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	timeRecords, err := a.db.Get(id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, timeRecords)
}

func (a *TimeRecordAPI) AllTimeRecordsEndPoint(c *gin.Context) {
	timeRecords, err := a.db.All()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, timeRecords)
}

func (a *TimeRecordAPI) TimeRecordsByEmployeeEndPoint(c *gin.Context) {
	idReq := c.Query("id")
	id, err := strconv.Atoi(idReq)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	timeRecords, err := a.db.ByEmployeeId(id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, timeRecords)
}

type TimeRecordsByDateRequest struct {
	EmployeeId *int           `json:"employee_id"`
	Start      model.DateTime `json:"start"`
	End        model.DateTime `json:"end"`
}

func (a *TimeRecordAPI) TimeRecordsByDateEndPoint(c *gin.Context) {
	var req TimeRecordsByDateRequest
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	// Какой должен быть формат времени? Я почти везде поставил time.Time
	// Или нужен метод для форматирования?

	timeRecords, err := a.db.ByDate(req.Start.Time(), req.End.Time(), req.EmployeeId)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, timeRecords)
}

func (a *TimeRecordAPI) TimeRecordLastByEmployeeEndPoint(c *gin.Context) {
	idReq := c.Query("id")
	id, err := strconv.Atoi(idReq)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	timeRecord, err := a.db.LastByEmployeeId(id)
	if err != nil {
		c.JSON(204, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, timeRecord)
}

func TimeRecordGroup(r *gin.Engine, db store.TimeRecordDb, url string) {
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
