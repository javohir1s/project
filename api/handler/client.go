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

// @Summary create a Client
// @Description Create Client
// @Tags Client
// @Accept json
// @Produce json
// @Param object body models.CreateClient true "Client ID"
// @Success 200 {object} models.Client "Client details"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 404 {object} ErrorResponse "Client not found"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /client [post]
func (h *Handler) CreateClient(c *gin.Context) {

	var createClient models.CreateClient
	err := c.ShouldBindJSON(&createClient)
	if err != nil {
		handleResponse(c, 400, "ShouldBindJSON err:"+err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.CtxTimeout)
	defer cancel()

	resp, err := h.strg.Client().Create(ctx, &createClient)
	if err != nil {
		handleResponse(c, http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, http.StatusCreated, resp)
}

// @Summary Get a Client by ID
// @Description Get Client details by its ID.
// @Tags Client
// @Accept json
// @Produce json
// @Param id path string true "Client ID"
// @Success 200 {object} models.Client "Client details"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 404 {object} ErrorResponse "Client not found"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /client/{id} [get]
func (h *Handler) GetByIDClient(c *gin.Context) {

	var id = c.Param("id")

	if !helpers.IsValidUUID(id) {
		handleResponse(c, http.StatusBadRequest, "id is not uuid")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.CtxTimeout)
	defer cancel()

	resp, err := h.strg.Client().GetByID(ctx, &models.ClientPrimaryKey{Id: id})
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

// @Summary Get List Client
// @Description Get List Client details by its ok.
// @Tags Client
// @Accept json
// @Produce json
// @Param offset query int false "Offset"
// @Param limit query int false "Limit"
// @Param search query string false "search"
// @Success 200 {object} models.Client "Client details"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 404 {object} ErrorResponse "Client not found"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /client [get]
func (h *Handler) GetListClient(c *gin.Context) {

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
	fmt.Println(search,limit,offset)
	if err != nil {
		handleResponse(c, http.StatusBadRequest, "invalid query search")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.CtxTimeout)
	defer cancel()

	var (
		resp = &models.GetListClientResponse{}
	)

	if len(resp.Clients) <= 0 {
		resp, err = h.strg.Client().GetList(ctx, &models.GetListClientRequest{
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

// @Summary Update Client
// @Description Get List Client details by its ok.
// @Tags Client
// @Accept json
// @Produce json
// @Param object body models.UpdateClient true "models.UpdateClient"
// @Param id path string true "id"
// @Success 200 {object} models.Client "Client details"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 404 {object} ErrorResponse "Client not found"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /client/{id} [put]
func (h *Handler) UpdateClient(c *gin.Context) {

	var updateClient models.UpdateClient

	err := c.ShouldBindJSON(&updateClient)
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

	updateClient.Id = id

	rowsAffected, err := h.strg.Client().Update(ctx, &updateClient)
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

	resp, err := h.strg.Client().GetByID(ctx, &models.ClientPrimaryKey{Id: updateClient.Id})
	if err != nil {
		handleResponse(c, http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, http.StatusAccepted, resp)
}

// @Summary Get List Client
// @Description Get List Client details by its ok
// @Tags Client
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} models.Client "Client details"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 404 {object} ErrorResponse "Client not found"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /client/{id} [delete]
func (h *Handler) DeleteClient(c *gin.Context) {
	var id = c.Param("id")

	if !helpers.IsValidUUID(id) {
		handleResponse(c, http.StatusBadRequest, "id is not uuid")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.CtxTimeout)
	defer cancel()

	err := h.strg.Client().Delete(ctx, &models.ClientPrimaryKey{Id: id})
	if err != nil {
		handleResponse(c, http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, http.StatusOK, "deleted")

}
