package handler

import (
	"context"
	"errors"
	"exam/config"
	"exam/models"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

// @Summary MakePay
// @Description Get List Coming details by its ok.
// @Tags Pay
// @Accept json
// @Produce json
// @Param sale_id query string ture "sale_id"
// @Param money query float64 true "Pay money"
// @Success 200 {object} models.Coming "Coming details"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 404 {object} ErrorResponse "Coming not found"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /make_pay [post]
func (h *Handler) MakePay(c *gin.Context) {
	var (
		incrementID = c.Query("sale_id")
		money       = cast.ToFloat64(c.Query("money"))
		Id          string
		ClientID    string
		BranchID    string
		IncrementID string
		TotalPrice  float64
	)

	fmt.Println(money)
	fmt.Println(incrementID)

	ctx, cancel := context.WithTimeout(context.Background(), config.CtxTimeout)
	defer cancel()

	saleList, err := h.strg.Sale().GetList(ctx, &models.GetListSaleRequest{Limit: 10000})
	if err != nil {
		handleResponse(c, http.StatusInternalServerError, err)
		return
	}

	saleFound := false
	for _, v := range saleList.Sales {
		if strings.EqualFold(strings.TrimSpace(v.IncrementID), strings.TrimSpace(incrementID)) {
			saleFound = true
			if v.TotalPrice/2 < money {
				Id = v.Id
				ClientID = v.ClientID
				BranchID = v.BranchID
				IncrementID = v.IncrementID
				TotalPrice = v.TotalPrice
			}
		}
	}

	if !saleFound {
		handleResponse(c, http.StatusNotFound, errors.New("sale not found"))
		return
	}

	_, err = h.strg.Sale().Update(ctx, &models.UpdateSale{
		Id:          Id,
		ClientID:    ClientID,
		BranchID:    BranchID,
		IncrementID: IncrementID,
		TotalPrice:  TotalPrice,
		Paid:        money,
		Debd:        TotalPrice - money,
	})

	if err != nil {
		handleResponse(c, http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, http.StatusOK, "successful payment")
}

// Report godoc
// @ID    Report
// @Router /report [GET]
// @Summary    Report
// @Description   Report
// @Tags Report
// @Accept json
// @Produce    json
// @Param from_date query string false "Start date in Year-Month-Day format"
// @Param to_date query string false "End date in Year-Month-Day format"
// @Success 201    {object} Response{data=models.Client} "OverallReportBody"
// @Response 400 {object} Response{data=string}    "Invalid Argument"
// @Failure 500     {object} Response{data=string}    "Server Error"
func (h *Handler) OverallReport(c *gin.Context) {
	fromDateStr := c.Query("from_date")
	toDateStr := c.Query("to_date")

	var fromDate, toDate time.Time
	var err error

	if fromDateStr != "" {
		fromDate, err = time.Parse("2006-01-02", fromDateStr)
		if err != nil {
			handleResponse(c, http.StatusBadRequest, err.Error())
			return
		}
	}

	if toDateStr != "" {
		toDate, err = time.Parse("2006-01-02", toDateStr)
		if err != nil {
			handleResponse(c, http.StatusBadRequest, err.Error())
			return
		}
	}

	clent_resp, err := h.strg.Report().GetListReport(c, fromDate, toDate)
	if err != nil {
		handleResponse(c, http.StatusBadRequest, err)
		return
	}

	sale_resp, err := h.strg.Report().GetListSaleBranch(c, fromDate, toDate)
	if err != nil {
		handleResponse(c, http.StatusBadRequest, err)
		return
	}

	handleResponse(c, http.StatusOK, models.OverallReport{
		Clients:           clent_resp.Clients,
		BranchSaleReports: *sale_resp,
	})
}
