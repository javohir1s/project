package handler

import (
	"context"
	"database/sql"
	"net/http"

	"exam/config"
	"exam/models"
	"exam/pkg/helpers"

	"github.com/gin-gonic/gin"
)

// @Summary Create a new Coming
// @Description Create a new Coming in the market system.
// @Tags Coming
// @Accept json
// @Produce json
// @Param Coming body models.CreateComing true "Coming information"
// @Success 201 {object} models.Coming "Created Coming"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /coming [post]
func (h *Handler) CreateComing(c *gin.Context) {

	var createComing models.CreateComing
	err := c.ShouldBindJSON(&createComing)
	if err != nil {
		handleResponse(c, 400, "ShouldBindJSON err:"+err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.CtxTimeout)
	defer cancel()

	incrementId, err := h.strg.IncrementID().GetLast(ctx, "coming", "increment_id")
	if err != nil {
		handleResponse(c, http.StatusInternalServerError, err)
		return
	}

	createComing.IncrementID= "C-" + incrementId

	resp, err := h.strg.Coming().Create(ctx, &createComing)
	if err != nil {
		handleResponse(c, http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, http.StatusCreated, resp)
}

// @Summary Get an Coming by ID
// @Description Get Coming details by its ID.
// @Tags Coming
// @Accept json
// @Produce json
// @Param id path string true "Coming ID"
// @Success 200 {object} models.Coming "Coming details"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 404 {object} ErrorResponse "Coming not found"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /coming/{id} [get]
func (h *Handler) GetByIDComing(c *gin.Context) {

	var id = c.Param("id")
	if !helpers.IsValidUUID(id) {
		handleResponse(c, http.StatusBadRequest, "id is not uuid")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.CtxTimeout)
	defer cancel()

	resp, err := h.strg.Coming().GetByID(ctx, &models.ComingPrimaryKey{Id: id})
	if err == sql.ErrNoRows {
		handleResponse(c, http.StatusBadRequest, "no rows in result set")
		return
	}

	if err != nil {
		handleResponse(c, http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, http.StatusOK, resp)
}

// @Summary Get a list of Comings
// @Description Get a list of Comings with optional filtering.
// @Tags Coming
// @Accept json
// @Produce json
// @Param limit query int false "Number of items to return (default 10)"
// @Param offset query int false "Number of items to skip (default 0)"
// @Success 200 {array} models.Coming "List of Comings"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /coming [get]
func (h *Handler) GetListComing(c *gin.Context) {

	limit, err := getIntegerOrDefaultValue(c.Query("limit"), 10)
	if err != nil {
		handleResponse(c, http.StatusBadRequest, "invalid query limit")
		return
	}

	offset, err := getIntegerOrDefaultValue(c.Query("offset"), 0)
	if err != nil {
		handleResponse(c, http.StatusBadRequest, "invalid query offset")
		return
	}


	ctx, cancel := context.WithTimeout(context.Background(), config.CtxTimeout)
	defer cancel()

	var (
		resp = &models.GetListComingResponse{}
	)

	if len(resp.Cominges) <= 0 {
		resp, err = h.strg.Coming().GetList(ctx, &models.GetListComingRequest{
			Limit:  limit,
			Offset: offset,
		})
		if err != nil {
			handleResponse(c, http.StatusInternalServerError, err)
			return
		}

		if err != nil {
			handleResponse(c, http.StatusInternalServerError, err)
			return
		}

	}

	handleResponse(c, http.StatusOK, resp)
}

// @Summary Update an Coming
// @Description Update an existing Coming.
// @Tags Coming
// @Accept json
// @Produce json
// @Param id path string true "Coming ID"
// @Param Coming body models.UpdateComing true "Updated Coming information"
// @Success 202 {object} models.Coming "Updated Coming"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 404 {object} ErrorResponse "Coming not found"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /coming/{id} [put]
func (h *Handler) UpdateComing(c *gin.Context) {

	var updateComing models.UpdateComing

	err := c.ShouldBindJSON(&updateComing)
	if err != nil {
		handleResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var id = c.Param("id")
	if !helpers.IsValidUUID(id) {
		handleResponse(c, http.StatusBadRequest, "id is not uuid")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.CtxTimeout)
	defer cancel()

	rowsAffected, err := h.strg.Coming().Update(ctx, &updateComing)
	if err != nil {
		handleResponse(c, http.StatusInternalServerError, err)
		return
	}

	if rowsAffected == 0 {
		handleResponse(c, http.StatusBadRequest, "no rows affected")
		return
	}

	ctx, cancel = context.WithTimeout(context.Background(), config.CtxTimeout)
	defer cancel()

	resp, err := h.strg.Coming().GetByID(ctx, &models.ComingPrimaryKey{Id: updateComing.Id})
	if err != nil {
		handleResponse(c, http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, http.StatusAccepted, resp)
}

// @Summary Delete an Coming
// @Description Delete an existing Coming.
// @Tags Coming
// @Accept json
// @Produce json
// @Param id path string true "Coming ID"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /coming/{id} [delete]
func (h *Handler) DeleteComing(c *gin.Context) {
	var id = c.Param("id")

	if !helpers.IsValidUUID(id) {
		handleResponse(c, http.StatusBadRequest, "id is not uuid")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.CtxTimeout)
	defer cancel()

	err := h.strg.Coming().Delete(ctx, &models.ComingPrimaryKey{Id: id})
	if err != nil {
		handleResponse(c, http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, http.StatusNoContent, nil)
}
