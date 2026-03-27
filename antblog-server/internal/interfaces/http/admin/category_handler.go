// Package admin 后台管理 HTTP 处理层（需要 Admin 角色）。
package admin

import (
	"github.com/gin-gonic/gin"

	appcategory "antblog/internal/application/category"
	"antblog/pkg/response"
	"antblog/pkg/validator"
)

// CategoryHandler 后台分类管理处理器
type CategoryHandler struct {
	useCase appcategory.ICategoryUseCase
}

// NewAdminCategoryHandler 创建后台分类处理器
func NewAdminCategoryHandler(uc appcategory.ICategoryUseCase) *CategoryHandler {
	return &CategoryHandler{useCase: uc}
}

// ListCategories godoc
// @Summary  后台获取所有分类
// @Tags     admin-category
// @Security BearerAuth
// @Produce  json
// @Success  200 {object} response.Response{data=[]appcategory.CategoryResp}
// @Router   /api/admin/categories [get]
func (h *CategoryHandler) ListCategories(c *gin.Context) {
	list, err := h.useCase.ListCategories(c.Request.Context())
	if err != nil {
		response.FailWithError(c, err)
		return
	}
	response.OK(c, list)
}

// GetCategory godoc
// @Summary  后台按 ID 获取分类
// @Tags     admin-category
// @Security BearerAuth
// @Produce  json
// @Param    id path int true "分类 ID"
// @Success  200 {object} response.Response{data=appcategory.CategoryResp}
// @Router   /api/admin/categories/{id} [get]
func (h *CategoryHandler) GetCategory(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		response.BadRequest(c, "无效的分类 ID")
		return
	}

	cat, err := h.useCase.GetCategoryByID(c.Request.Context(), id)
	if err != nil {
		response.FailWithError(c, err)
		return
	}
	response.OK(c, cat)
}

// CreateCategory godoc
// @Summary  创建分类
// @Tags     admin-category
// @Security BearerAuth
// @Accept   json
// @Produce  json
// @Param    body body appcategory.CreateCategoryReq true "分类信息"
// @Success  201  {object} response.Response{data=appcategory.CategoryResp}
// @Router   /api/admin/categories [post]
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var req appcategory.CreateCategoryReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求体解析失败："+err.Error())
		return
	}
	if err := validator.Validate(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	cat, err := h.useCase.CreateCategory(c.Request.Context(), &req)
	if err != nil {
		response.FailWithError(c, err)
		return
	}
	response.Created(c, cat)
}

// UpdateCategory godoc
// @Summary  更新分类
// @Tags     admin-category
// @Security BearerAuth
// @Accept   json
// @Produce  json
// @Param    id   path int true "分类 ID"
// @Param    body body appcategory.UpdateCategoryReq true "分类信息"
// @Success  200  {object} response.Response{data=appcategory.CategoryResp}
// @Router   /api/admin/categories/{id} [put]
func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		response.BadRequest(c, "无效的分类 ID")
		return
	}

	var req appcategory.UpdateCategoryReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求体解析失败："+err.Error())
		return
	}
	if err := validator.Validate(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	cat, err := h.useCase.UpdateCategory(c.Request.Context(), id, &req)
	if err != nil {
		response.FailWithError(c, err)
		return
	}
	response.OK(c, cat)
}

// DeleteCategory godoc
// @Summary  删除分类
// @Tags     admin-category
// @Security BearerAuth
// @Produce  json
// @Param    id path int true "分类 ID"
// @Success  200 {object} response.Response
// @Router   /api/admin/categories/{id} [delete]
func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		response.BadRequest(c, "无效的分类 ID")
		return
	}

	if err := h.useCase.DeleteCategory(c.Request.Context(), id); err != nil {
		response.FailWithError(c, err)
		return
	}
	response.OKWithMsg(c, "删除成功", nil)
}
