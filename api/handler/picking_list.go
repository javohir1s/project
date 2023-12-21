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

// @Summary create a PickingList
// @Description Create PickingList
// @Tags PickingList
// @Accept json
// @Produce json
// @Param object body models.CreatePickingList true "PickingList ID"
// @Success 200 {object} models.PickingList "PickingList details"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 404 {object} ErrorResponse "PickingList not found"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /picking_list [post]
func (h *Handler) CreatePickingList(c *gin.Context) {

	var createPickingList models.PickingList
	err := c.ShouldBindJSON(&createPickingList)
	if err != nil {
		handleResponse(c, 400, "ShouldBindJSON err:"+err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.CtxTimeout)
	defer cancel()

	coming, err := h.strg.Coming().GetList(ctx, &models.GetListComingRequest{
		Limit: 10000,
	})
	if err != nil {
		handleResponse(c, http.StatusInternalServerError, err)
		return
	}
	var (
		cominggg     models.Coming
		id           sql.NullString
		productID    sql.NullString
		productName  sql.NullString
		quantity     sql.NullInt64
		coming_price sql.NullFloat64
		salePrice    sql.NullFloat64
		branchID     sql.NullString
	)
	for _, v := range coming.Cominges {
		if v.IncrementID == createPickingList.ComingIncrementID {
			createPickingList.ComingID = v.Id
			cominggg.BranchID = v.BranchID
		}
	}

	resp, err := h.strg.PickingList().Create(ctx, &createPickingList)
	if err != nil {
		handleResponse(c, http.StatusInternalServerError, err)
		return
	}

	productResp, err := h.strg.Product().GetByID(ctx, &models.ProductPrimaryKey{Id: createPickingList.Product_ID})
	if err != nil {
		handleResponse(c, http.StatusBadRequest, err)
		return
	}

	listRemainder, err := h.strg.Remainder().GetList(ctx, &models.GetListRemainderRequest{Limit: 1000})
	if len(listRemainder.Remainders) == 0 {
		_, err = h.strg.Remainder().Create(ctx, &models.Remainder{
			ProductID:   productResp.Id,
			Name:        productResp.Name,
			Quantity:    createPickingList.Quantity,
			ComingPrice: createPickingList.Price,
			SalePrice:   productResp.Price,
			BranchID:    cominggg.BranchID,
		})
		handleResponse(c, http.StatusAccepted, err)
		return
	}
	if err != nil {
		handleResponse(c, http.StatusBadRequest, err)
		return
	}

	for _, v := range listRemainder.Remainders {

		if v.BranchID == cominggg.BranchID && v.ProductID == createPickingList.Product_ID {
			id.String = v.Id
			productID.String = v.ProductID
			productName.String = v.Name
			quantity.Int64 = int64(v.Quantity)
			coming_price.Float64 = v.ComingPrice
			salePrice.Float64 = productResp.Price
			branchID.String = v.BranchID
			break
		}
	}

	_, err = h.strg.Remainder().Update(ctx, &models.Remainder{
		Id:          id.String,
		ProductID:   productID.String,
		Name:        productName.String,
		Quantity:    int(quantity.Int64) + createPickingList.Quantity,
		ComingPrice: coming_price.Float64,
		SalePrice:   salePrice.Float64,
		BranchID:    branchID.String,
	})
	if err != nil {
		handleResponse(c, http.StatusBadRequest, err)
		return
	}

	handleResponse(c, http.StatusCreated, resp)
}

// @Summary Get a PickingList by ID
// @Description Get PickingList details by its ID.
// @Tags PickingList
// @Accept json
// @Produce json
// @Param id path string true "PickingList ID"
// @Success 200 {object} models.PickingList "PickingList details"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 404 {object} ErrorResponse "PickingList not found"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /picking_list/{id} [get]
func (h *Handler) GetByIDPickingList(c *gin.Context) {

	var id = c.Param("id")

	if !helpers.IsValidUUID(id) {
		handleResponse(c, http.StatusBadRequest, "id is not uuid")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.CtxTimeout)
	defer cancel()

	resp, err := h.strg.PickingList().GetByID(ctx, &models.PickingListPrimaryKey{Id: id})
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

// @Summary Get List PickingList
// @Description Get List PickingList details by its ok.
// @Tags PickingList
// @Accept json
// @Produce json
// @Param offset path int false "Offset"
// @Param limit path int false "Limit"
// @Success 200 {object} models.PickingList "PickingList details"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 404 {object} ErrorResponse "PickingList not found"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /picking_list [get]
func (h *Handler) GetListPickingList(c *gin.Context) {

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
		resp = &models.GetListPickingListResponse{}
	)

	if len(resp.Pickinges) <= 0 {
		resp, err = h.strg.PickingList().GetList(ctx, &models.GetListPickingListRequest{
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

// @Summary Update PickingList
// @Description Get List PickingList details by its ok.
// @Tags PickingList
// @Accept json
// @Produce json
// @Param object body models.UpdatePickingList true "models.UpdatePickingList"
// @Param id path string true "id"
// @Success 200 {object} models.PickingList "PickingList details"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 404 {object} ErrorResponse "PickingList not found"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /picking_list/{id} [put]
func (h *Handler) UpdatePickingList(c *gin.Context) {

	var updatePickingList models.PickingList

	err := c.ShouldBindJSON(&updatePickingList)
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

	updatePickingList.ID = id

	rowsAffected, err := h.strg.PickingList().Update(ctx, &updatePickingList)
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

	resp, err := h.strg.PickingList().GetByID(ctx, &models.PickingListPrimaryKey{Id: updatePickingList.ID})
	if err != nil {
		handleResponse(c, http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, http.StatusAccepted, resp)
}

// @Summary Get List PickingList
// @Description Get List PickingList details by its ok
// @Tags PickingList
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} models.PickingList "PickingList details"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 404 {object} ErrorResponse "PickingList not found"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /picking_list/{id} [delete]
func (h *Handler) DeletePickingList(c *gin.Context) {
	var id = c.Param("id")

	if !helpers.IsValidUUID(id) {
		handleResponse(c, http.StatusBadRequest, "id is not uuid")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.CtxTimeout)
	defer cancel()

	err := h.strg.PickingList().Delete(ctx, &models.PickingListPrimaryKey{Id: id})
	if err != nil {
		handleResponse(c, http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, http.StatusOK, "deleted")

}
