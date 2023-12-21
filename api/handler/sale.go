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

// @Summary create a Sale
// @Description Create Sale
// @Tags Sale
// @Accept json
// @Produce json
// @Param object body models.CreateSale true "Sale ID"
// @Success 200 {object} models.Sale "Sale details"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 404 {object} ErrorResponse "Sale not found"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /sale [post]
func (h *Handler) CreateSale(c *gin.Context) {

	var createSale models.CreateSale
	err := c.ShouldBindJSON(&createSale)
	if err != nil {
		handleResponse(c, 400, "ShouldBindJSON err:"+err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.CtxTimeout)
	defer cancel()

	incrementId, err := h.strg.IncrementID().GetLast(ctx, "sale", "increment_id")
	if err != nil {
		handleResponse(c, http.StatusInternalServerError, err)
		return
	}

	createSale.IncrementID = "S-" + incrementId

	resp, err := h.strg.Sale().Create(ctx, &createSale)
	if err != nil {
		handleResponse(c, http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, http.StatusCreated, resp)
}

// @Summary Get a Sale by ID
// @Description Get Sale details by its ID.
// @Tags Sale
// @Accept json
// @Produce json
// @Param id path string true "Sale ID"
// @Success 200 {object} models.Sale "Sale details"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 404 {object} ErrorResponse "Sale not found"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /sale/{id} [get]
func (h *Handler) GetByIDSale(c *gin.Context) {

	var id = c.Param("id")

	if !helpers.IsValidUUID(id) {
		handleResponse(c, http.StatusBadRequest, "id is not uuid")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.CtxTimeout)
	defer cancel()

	resp, err := h.strg.Sale().GetByID(ctx, &models.SalePrimaryKey{Id: id})
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

// @Summary Get List Sale
// @Description Get List Sale details by its ok.
// @Tags Sale
// @Accept json
// @Produce json
// @Param offset query int false "Offset"
// @Param limit query int false "Limit"
// @Success 200 {object} models.Sale "Sale details"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 404 {object} ErrorResponse "Sale not found"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /sale [get]
func (h *Handler) GetListSale(c *gin.Context) {

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
		resp = &models.GetListSaleResponse{}
	)

	if len(resp.Sales) <= 0 {
		resp, err = h.strg.Sale().GetList(ctx, &models.GetListSaleRequest{
			Limit:  limit,
			Offset: offset,
		})
		if err != nil {
			handleResponse(c, http.StatusInternalServerError, err)
			return
		}

	}

	handleResponse(c, http.StatusOK, resp)
}

// @Summary Update Sale
// @Description Get List Sale details by its ok.
// @Tags Sale
// @Accept json
// @Produce json
// @Param object body models.UpdateSale true "models.UpdateSale"
// @Param id path string true "id"
// @Success 200 {object} models.Sale "Sale details"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 404 {object} ErrorResponse "Sale not found"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /sale/{id} [put]
func (h *Handler) UpdateSale(c *gin.Context) {

	var updateSale models.UpdateSale

	err := c.ShouldBindJSON(&updateSale)
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

	updateSale.Id = id

	rowsAffected, err := h.strg.Sale().Update(ctx, &updateSale)
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

	resp, err := h.strg.Sale().GetByID(ctx, &models.SalePrimaryKey{Id: updateSale.Id})
	if err != nil {
		handleResponse(c, http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, http.StatusAccepted, resp)
}

// @Summary Get List Sale
// @Description Get List Sale details by its ok
// @Tags Sale
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} models.Sale "Sale details"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 404 {object} ErrorResponse "Sale not found"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /sale/{id} [delete]
func (h *Handler) DeleteSale(c *gin.Context) {
	var id = c.Param("id")

	if !helpers.IsValidUUID(id) {
		handleResponse(c, http.StatusBadRequest, "id is not uuid")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.CtxTimeout)
	defer cancel()

	err := h.strg.Sale().Delete(ctx, &models.SalePrimaryKey{Id: id})
	if err != nil {
		handleResponse(c, http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, http.StatusOK, "deleted")

}
