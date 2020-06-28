package controllers

import (
	"net/http"
	"net/url"

	"github.com/google/uuid"

	"github.com/egeback/mediadownloader/pkg/models"
	"github.com/gin-gonic/gin"
)

var scheduler = models.NewScheduler()

// AddJob ...
func (c *Controller) AddJob(ctx *gin.Context) {
	u := ctx.DefaultQuery("url", "")
	_, err := url.ParseRequestURI(u)

	if u == "" || err != nil {
		c.createErrorResponse(ctx, http.StatusBadRequest, 100, "no valid url provided")
		return
	}

	download := models.AddDownload(u)
	//download := actions.AddDownload("https://www.svtplay.se/video/21868842/palmegruppen-tar-langlunch")
	//download := actions.AddDownload("https://www.svtplay.se/video/26987573/you-were-never-really-here")
	//download.Start()
	id := uuid.New()
	uuid := id.String()
	//models.AddJob(models.Job{UUID: uuid, Download: &download})
	scheduler.AddJob(models.Job{UUID: uuid, Download: &download})

	ctx.JSON(http.StatusAccepted, gin.H{
		"job_id": uuid,
	})
	return
}

//Jobs ...
func (c *Controller) Jobs(ctx *gin.Context) {
	//retVal := make([]models.Job, 0, len(models.AllJobs()))
	//for _, job := range models.AllJobs() {
	//	retVal = append(retVal, job)
	//}
	ctx.JSON(http.StatusOK, scheduler.GetJobs())
}

//GetJob ...
func (c *Controller) GetJob(ctx *gin.Context) {
	uuid := ctx.Param("uuid")

	_, exists := models.AllJobs()[uuid]
	if !exists {
		//ctx.JSON(http.StatusNotFound, gin.H{})
		c.createErrorResponse(ctx, http.StatusNotFound, 101, "job id does not exist")
		return
	}

	ctx.JSON(http.StatusAccepted, models.AllJobs()[uuid])
	return
}
