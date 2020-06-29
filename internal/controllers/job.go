package controllers

import (
	"net/http"
	"net/url"

	"github.com/google/uuid"

	"github.com/egeback/playdownloader/internal/models"
	"github.com/gin-gonic/gin"
)

//Scheduler global object
var scheduler = models.NewScheduler()

// AddJob function to add jobs to the API
// @Summary Add job
// @Description Add job to API for download
// @Tags jobs
// @Accept json
// @Produce json
// @Param url query string true "url to download" Format(str)
// @Success 200 {object} models.Job
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /jobs [post]
func (c *Controller) AddJob(ctx *gin.Context) {
	u := ctx.DefaultQuery("url", "")
	_, err := url.ParseRequestURI(u)

	if u == "" || err != nil {
		c.createErrorResponse(ctx, http.StatusBadRequest, 100, "no valid url provided")
		return
	}

	download := models.CreateDownload(u)

	id := uuid.New()
	uuid := id.String()

	scheduler.AddJob(models.Job{UUID: uuid, Download: &download})

	ctx.JSON(http.StatusAccepted, gin.H{
		"job_id": uuid,
	})
	return
}

//Jobs ...
// Jobs function to list all jobs active in the API
// @Summary List jobs
// @Description List all jobs active in the API
// @Tags jobs
// @Accept json
// @Produce json
// @Success 200 {array} models.Job
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /jobs [get]
func (c *Controller) Jobs(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, scheduler.GetJobs())
}

//GetJob function to get a specific job by UUID
// Jobs function to get a specific job by UUID
// @Summary Get job
// @Description Get a specific job by UUID
// @Tags jobs
// @Accept json
// @Produce json
// @Param uuid path string true "job uuid" Format(str)
// @Success 200 {object} models.Job
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /jobs/:uuid [get]
func (c *Controller) GetJob(ctx *gin.Context) {
	uuid := ctx.Param("uuid")

	job := scheduler.GetJobByUUID(uuid)
	if job == nil {
		c.createErrorResponse(ctx, http.StatusNotFound, 101, "job id does not exist")
		return
	}

	ctx.JSON(http.StatusAccepted, job)
	return
}

//StopJob ...
// Jobs function to stop a specific job by UUID
// @Summary Stop jobs
// @Description Stop a specific job by UUID
// @Tags jobs
// @Accept json
// @Produce json
// @Param uuid path string true "job uuid" Format(str)
// @Success 200 {object} controller.Response
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /jobs/:uuid [get]
func (c *Controller) StopJob(ctx *gin.Context) {
	uuid := ctx.Param("uuid")

	job := scheduler.GetJobByUUID(uuid)
	if job == nil {
		c.createErrorResponse(ctx, http.StatusNotFound, 101, "job id does not exist")
		return
	}

	err := job.Stop()
	if err != nil {
		c.createErrorResponse(ctx, http.StatusNotFound, 102, "Could not stop job")
		return
	}

	ctx.JSON(http.StatusAccepted, Response{Data: "ok"})
	return
}
