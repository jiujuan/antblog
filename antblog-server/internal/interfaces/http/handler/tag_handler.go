package handler

import (
	"github.com/gin-gonic/gin"

	apptag "antblog/internal/application/tag"
	"antblog/pkg/response"
)

// TagHandler 标签前台 HTTP 处理器
type TagHandler struct {
	useCase apptag.ITagUseCase
}

// NewTagHandler 创建标签处理器
func NewTagHandler(uc apptag.ITagUseCase) *TagHandler {
	return &TagHandler{useCase: uc}
}

// ListTags godoc
// @Summary  获取所有标签列表
// @Tags     tag
// @Produce  json
// @Success  200 {object} response.Response{data=[]apptag.TagResp}
// @Router   /api/v1/tags [get]
func (h *TagHandler) ListTags(c *gin.Context) {
	list, err := h.useCase.ListTags(c.Request.Context())
	if err != nil {
		response.FailWithError(c, err)
		return
	}
	response.OK(c, list)
}

// GetTagBySlug godoc
// @Summary  按 Slug 获取标签详情
// @Tags     tag
// @Produce  json
// @Param    slug path string true "标签 Slug"
// @Success  200  {object} response.Response{data=apptag.TagResp}
// @Router   /api/v1/tags/{slug} [get]
func (h *TagHandler) GetTagBySlug(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		response.BadRequest(c, "slug 不能为空")
		return
	}

	tag, err := h.useCase.GetTagBySlug(c.Request.Context(), slug)
	if err != nil {
		response.FailWithError(c, err)
		return
	}
	response.OK(c, tag)
}
