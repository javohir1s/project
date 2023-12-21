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

// @Summary create a SaleProduct
// @Description Create SaleProduct
// @Tags SaleProduct
// @Accept json
// @Produce json
// @Param object body models.CreateSaleProduct true "SaleProduct ID"
// @Success 200 {object} models.SaleProduct "SaleProduct details"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 404 {object} ErrorResponse "SaleProduct not found"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /saleproduct [post]
func (h *Handler) CreateSaleProduct(c *gin.Context) {

	var createSaleProduct models.CreateSaleProduct
	err := c.ShouldBindJSON(&createSaleProduct)
	if err != nil {
		handleResponse(c, 400, "ShouldBindJSON err:"+err.Error())
		return
	}

	

	ctx, cancel := context.WithTimeout(context.Background(), config.CtxTimeout)
	defer cancel()
	

	resp, err := h.strg.SaleProduct().Create(ctx, &createSaleProduct)
	if err != nil {
		handleResponse(c, http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, http.StatusCreated, resp)
}

// @Summary Get a SaleProduct by ID
// @Description Get SaleProduct details by its ID.
// @Tags SaleProduct
// @Accept json
// @Produce json
// @Param id path string true "SaleProduct ID"
// @Success 200 {object} models.SaleProduct "SaleProduct details"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 404 {object} ErrorResponse "SaleProduct not found"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /saleproduct/{id} [get]
func (h *Handler) GetByIDSaleProduct(c *gin.Context) {

	var id = c.Param("id")

	if !helpers.IsValidUUID(id) {
		handleResponse(c, http.StatusBadRequest, "id is not uuid")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.CtxTimeout)
	defer cancel()

	resp, err := h.strg.SaleProduct().GetByID(ctx, &models.SaleProductPrimaryKey{Id: id})
	if err == sql.ErrNoRows {
		handleResponse(c, http.StatusBadRequest, "no rows in result set >>>")
		return
	}

	if err != nil {
		handleResponse(c, http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, http.StatusOK, resp)
}

// @Summary Get List SaleProduct
// @Description Get List SaleProduct details by its ok.
// @Tags SaleProduct
// @Accept json
// @Produce json
// @Param offset query int false "Offset"
// @Param limit query int false "Limit"
// @Success 200 {object} models.SaleProduct "SaleProduct details"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 404 {object} ErrorResponse "SaleProduct not found"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /saleproduct [get]
func (h *Handler) GetListSaleProduct(c *gin.Context) {

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
		resp = &models.GetListSaleProductResponse{}
	)

	if len(resp.SaleProducts) <= 0 {
		resp, err = h.strg.SaleProduct().GetList(ctx, &models.GetListSaleProductRequest{
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

// @Summary Update SaleProduct
// @Description Get List SaleProduct details by its ok.
// @Tags SaleProduct
// @Accept json
// @Produce json
// @Param object body models.UpdateSaleProduct true "models.UpdateSaleProduct"
// @Param id path string true "id"
// @Success 200 {object} models.SaleProduct "SaleProduct details"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 404 {object} ErrorResponse "SaleProduct not found"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /saleproduct/{id} [put]
func (h *Handler) UpdateSaleProduct(c *gin.Context) {

	var updateSaleProduct models.UpdateSaleProduct

	err := c.ShouldBindJSON(&updateSaleProduct)
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

	updateSaleProduct.Id = id

	rowsAffected, err := h.strg.SaleProduct().Update(ctx, &updateSaleProduct)
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

	resp, err := h.strg.SaleProduct().GetByID(ctx, &models.SaleProductPrimaryKey{Id: updateSaleProduct.Id})
	if err != nil {
		handleResponse(c, http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, http.StatusAccepted, resp)
}

// @Summary Get List SaleProduct
// @Description Get List SaleProduct details by its ok
// @Tags SaleProduct
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} models.SaleProduct "SaleProduct details"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 404 {object} ErrorResponse "SaleProduct not found"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /saleproduct/{id} [delete]
func (h *Handler) DeleteSaleProduct(c *gin.Context) {
	var id = c.Param("id")

	if !helpers.IsValidUUID(id) {
		handleResponse(c, http.StatusBadRequest, "id is not uuid")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.CtxTimeout)
	defer cancel()

	err := h.strg.SaleProduct().Delete(ctx, &models.SaleProductPrimaryKey{Id: id})
	if err != nil {
		handleResponse(c, http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, http.StatusOK, "deleted")

}


