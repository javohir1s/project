package handler

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"exam/config"
	"exam/models"
	"exam/pkg/helpers"

	"github.com/gin-gonic/gin"
)

// @Summary create a Remainder
// @Description Create Remainder
// @Tags Remainder
// @Accept json
// @Produce json
// @Param object body models.CreateRemainder true "Remainder ID"
// @Success 200 {object} models.Remainder "Remainder details"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 404 {object} ErrorResponse "Remainder not found"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /remainder [post]
func (h *Handler) CreateRemainder(c *gin.Context) {

	var createRemainder models.Remainder
	err := c.ShouldBindJSON(&createRemainder)
	if err != nil {
		handleResponse(c, 400, "ShouldBindJSON err:"+err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.CtxTimeout)
	defer cancel()

	product, err := h.strg.Product().GetByID(ctx, &models.ProductPrimaryKey{Id: createRemainder.ProductID})
	if err != nil {
		handleResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	createRemainder.Name = product.Name
	fmt.Println(product)

	resp, err := h.strg.Remainder().Create(ctx, &createRemainder)
	if err != nil {
		handleResponse(c, http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, http.StatusCreated, resp)
}

// @Summary Get a Remainder by ID
// @Description Get Remainder details by its ID.
// @Tags Remainder
// @Accept json
// @Produce json
// @Param id path string true "Remainder ID"
// @Success 200 {object} models.Remainder "Remainder details"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 404 {object} ErrorResponse "Remainder not found"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /remainder/{id} [get]
func (h *Handler) GetByIDRemainder(c *gin.Context) {

	var id = c.Param("id")

	if !helpers.IsValidUUID(id) {
		handleResponse(c, http.StatusBadRequest, "id is not uuid")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.CtxTimeout)
	defer cancel()

	resp, err := h.strg.Remainder().GetByID(ctx, &models.RemainderPrimaryKey{Id: id})
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

// @Summary Get List Remainder
// @Description Get List Remainder details by its ok.
// @Tags Remainder
// @Accept json
// @Produce json
// @Param offset path int false "Offset"
// @Param limit path int false "Limit"
// @Param search path int false "search"
// @Success 200 {object} models.Remainder "Remainder details"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 404 {object} ErrorResponse "Remainder not found"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /remainder [get]
func (h *Handler) GetListRemainder(c *gin.Context) {

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

	search := c.Query("search")
	if err != nil {
		handleResponse(c, http.StatusBadRequest, "invalid query search")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.CtxTimeout)
	defer cancel()

	var (
		resp = &models.GetListRemainderResponse{}
	)

	if len(resp.Remainders) <= 0 {
		resp, err = h.strg.Remainder().GetList(ctx, &models.GetListRemainderRequest{
			Limit:  limit,
			Offset: offset,
			Search: search,
		})
		if err != nil {
			handleResponse(c, http.StatusInternalServerError, err)
			return
		}

	}

	handleResponse(c, http.StatusOK, resp)
}

// @Summary Update Remainder
// @Description Get List Remainder details by its ok.
// @Tags Remainder
// @Accept json
// @Produce json
// @Param object body models.UpdateRemainder true "models.UpdateRemainder"
// @Param id path string true "id"
// @Success 200 {object} models.Remainder "Remainder details"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 404 {object} ErrorResponse "Remainder not found"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /remainder/{id} [put]
func (h *Handler) UpdateRemainder(c *gin.Context) {

	var updateRemainder models.Remainder

	err := c.ShouldBindJSON(&updateRemainder)
	if err != nil {
		handleResponse(c, 400, "ShouldBindJSON err:"+err.Error())
		return
	}

	var id = c.Param("id")
	if !helpers.IsValidUUID(id) {
		handleResponse(c, http.StatusBadRequest, "id is not uuid")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.CtxTimeout)
	defer cancel()

	updateRemainder.Id = id

	rowsAffected, err := h.strg.Remainder().Update(ctx, &updateRemainder)
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

	resp, err := h.strg.Remainder().GetByID(ctx, &models.RemainderPrimaryKey{Id: updateRemainder.Id})
	if err != nil {
		handleResponse(c, http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, http.StatusAccepted, resp)
}

// @Summary Get List Remainder
// @Description Get List Remainder details by its ok
// @Tags Remainder
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} models.Remainder "Remainder details"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 404 {object} ErrorResponse "Remainder not found"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /remainder/{id} [delete]
func (h *Handler) DeleteRemainder(c *gin.Context) {
	var id = c.Param("id")

	if !helpers.IsValidUUID(id) {
		handleResponse(c, http.StatusBadRequest, "id is not uuid")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.CtxTimeout)
	defer cancel()

	err := h.strg.Remainder().Delete(ctx, &models.RemainderPrimaryKey{Id: id})
	if err != nil {
		handleResponse(c, http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, http.StatusOK, "deleted")

}
