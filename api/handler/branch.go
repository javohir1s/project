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

// @Summary create a branch
// @Description Create branch
// @Tags branch
// @Accept json
// @Produce json
// @Param object body models.CreateBranch true "Branch ID"
// @Success 200 {object} models.Branch "Branch details"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 404 {object} ErrorResponse "Branch not found"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /branch [post]
func (h *Handler) CreateBranch(c *gin.Context) {

	var createBranch models.CreateBranch
	err := c.ShouldBindJSON(&createBranch)
	if err != nil {
		handleResponse(c, 400, "ShouldBindJSON err:"+err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.CtxTimeout)
	defer cancel()

	resp, err := h.strg.Branch().Create(ctx, &createBranch)
	if err != nil {
		handleResponse(c, http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, http.StatusCreated, resp)
}

// @Summary Get a branch by ID
// @Description Get branch details by its ID.
// @Tags branch
// @Accept json
// @Produce json
// @Param id path string true "Branch ID"
// @Success 200 {object} models.Branch "Branch details"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 404 {object} ErrorResponse "Branch not found"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /branch/{id} [get]
func (h *Handler) GetByIDBranch(c *gin.Context) {

	var id = c.Param("id")

	if !helpers.IsValidUUID(id) {
		handleResponse(c, http.StatusBadRequest, "id is not uuid")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.CtxTimeout)
	defer cancel()

	resp, err := h.strg.Branch().GetByID(ctx, &models.BranchPrimaryKey{Id: id})
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

// @Summary Get List branch
// @Description Get List branch details by its ok.
// @Tags branch
// @Accept json
// @Produce json
// @Param offset query int false "Offset"
// @Param limit query int false "Limit"
// @Param search query string false "search"
// @Success 200 {object} models.Branch "Branch details"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 404 {object} ErrorResponse "Branch not found"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /branch [get]
func (h *Handler) GetListBranch(c *gin.Context) {

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
		resp = &models.GetListBranchResponse{}
	)

	if len(resp.Branches) <= 0 {
		resp, err = h.strg.Branch().GetList(ctx, &models.GetListBranchRequest{
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

// @Summary Update branch
// @Description Get List branch details by its ok.
// @Tags branch
// @Accept json
// @Produce json
// @Param object body models.UpdateBranch true "models.UpdateBranch"
// @Param id path string true "id"
// @Success 200 {object} models.Branch "Branch details"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 404 {object} ErrorResponse "Branch not found"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /branch/{id} [put]
func (h *Handler) UpdateBranch(c *gin.Context) {

	var updateBranch models.UpdateBranch

	err := c.ShouldBindJSON(&updateBranch)
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

	updateBranch.Id = id

	rowsAffected, err := h.strg.Branch().Update(ctx, &updateBranch)
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

	resp, err := h.strg.Branch().GetByID(ctx, &models.BranchPrimaryKey{Id: updateBranch.Id})
	if err != nil {
		handleResponse(c, http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, http.StatusAccepted, resp)
}

// @Summary Get List branch
// @Description Get List branch details by its ok
// @Tags branch
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} models.Branch "Branch details"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 404 {object} ErrorResponse "Branch not found"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /branch/{id} [delete]
func (h *Handler) DeleteBranch(c *gin.Context) {
	var id = c.Param("id")

	if !helpers.IsValidUUID(id) {
		handleResponse(c, http.StatusBadRequest, "id is not uuid")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.CtxTimeout)
	defer cancel()

	err := h.strg.Branch().Delete(ctx, &models.BranchPrimaryKey{Id: id})
	if err != nil {
		handleResponse(c, http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, http.StatusOK, "deleted")

}
