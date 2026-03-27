package admin

import (
	"github.com/gin-gonic/gin"

	apptag "antblog/internal/application/tag"
	"antblog/pkg/response"
	"antblog/pkg/validator"
)

// TagHandler 后台标签管理处理器
type TagHandler struct {
	useCase apptag.ITagUseCase
}

// NewAdminTagHandler 创建后台标签处理器
func NewAdminTagHandler(uc apptag.ITagUseCase) *TagHandler {
	return &TagHandler{useCase: uc}
}

// ListTags godoc
// @Summary  后台获取所有标签
// @Tags     admin-tag
// @Security BearerAuth
// @Produce  json
// @Success  200 {object} response.Response{data=[]apptag.TagResp}
// @Router   /api/admin/tags [get]
func (h *TagHandler) ListTags(c *gin.Context) {
	list, err := h.useCase.ListTags(c.Request.Context())
	if err != nil {
		response.FailWithError(c, err)
		return
	}
	response.OK(c, list)
}

// GetTag godoc
// @Summary  后台按 ID 获取标签
// @Tags     admin-tag
// @Security BearerAuth
// @Produce  json
// @Param    id path int true "标签 ID"
// @Success  200 {object} response.Response{data=apptag.TagResp}
// @Router   /api/admin/tags/{id} [get]
func (h *TagHandler) GetTag(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		response.BadRequest(c, "无效的标签 ID")
		return
	}

	tag, err := h.useCase.GetTagByID(c.Request.Context(), id)
	if err != nil {
		response.FailWithError(c, err)
		return
	}
	response.OK(c, tag)
}

// CreateTag godoc
// @Summary  创建标签
// @Tags     admin-tag
// @Security BearerAuth
// @Accept   json
// @Produce  json
// @Param    body body apptag.CreateTagReq true "标签信息"
// @Success  201  {object} response.Response{data=apptag.TagResp}
// @Router   /api/admin/tags [post]
func (h *TagHandler) CreateTag(c *gin.Context) {
	var req apptag.CreateTagReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求体解析失败："+err.Error())
		return
	}
	if err := validator.Validate(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	tag, err := h.useCase.CreateTag(c.Request.Context(), &req)
	if err != nil {
		response.FailWithError(c, err)
		return
	}
	response.Created(c, tag)
}

// BatchCreateTags godoc
// @Summary  批量创建标签（幂等）
// @Tags     admin-tag
// @Security BearerAuth
// @Accept   json
// @Produce  json
// @Param    body body apptag.BatchCreateTagReq true "标签列表"
// @Success  201  {object} response.Response{data=[]apptag.TagResp}
// @Router   /api/admin/tags/batch [post]
func (h *TagHandler) BatchCreateTags(c *gin.Context) {
	var req apptag.BatchCreateTagReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求体解析失败："+err.Error())
		return
	}
	if err := validator.Validate(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	tags, err := h.useCase.BatchCreateTags(c.Request.Context(), &req)
	if err != nil {
		response.FailWithError(c, err)
		return
	}
	response.Created(c, tags)
}

// UpdateTag godoc
// @Summary  更新标签
// @Tags     admin-tag
// @Security BearerAuth
// @Accept   json
// @Produce  json
// @Param    id   path int                true "标签 ID"
// @Param    body body apptag.UpdateTagReq true "标签信息"
// @Success  200  {object} response.Response{data=apptag.TagResp}
// @Router   /api/admin/tags/{id} [put]
func (h *TagHandler) UpdateTag(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		response.BadRequest(c, "无效的标签 ID")
		return
	}

	var req apptag.UpdateTagReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求体解析失败："+err.Error())
		return
	}
	if err := validator.Validate(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	tag, err := h.useCase.UpdateTag(c.Request.Context(), id, &req)
	if err != nil {
		response.FailWithError(c, err)
		return
	}
	response.OK(c, tag)
}

// DeleteTag godoc
// @Summary  删除标签
// @Tags     admin-tag
// @Security BearerAuth
// @Produce  json
// @Param    id path int true "标签 ID"
// @Success  200 {object} response.Response
// @Router   /api/admin/tags/{id} [delete]
func (h *TagHandler) DeleteTag(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		response.BadRequest(c, "无效的标签 ID")
		return
	}

	if err := h.useCase.DeleteTag(c.Request.Context(), id); err != nil {
		response.FailWithError(c, err)
		return
	}
	response.OKWithMsg(c, "删除成功", nil)
}
