package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	appcomment "antblog/internal/application/comment"
	"antblog/internal/interfaces/http/middleware"
	"antblog/pkg/response"
	"antblog/pkg/utils"
	"antblog/pkg/validator"
)

// CommentHandler 前台评论 HTTP 处理器
type CommentHandler struct {
	useCase appcomment.ICommentUseCase
}

// NewCommentHandler 创建前台评论处理器
func NewCommentHandler(uc appcomment.ICommentUseCase) *CommentHandler {
	return &CommentHandler{useCase: uc}
}

// ListComments godoc
// @Summary  获取文章评论树（顶级+子评论，仅已通过，分页）
// @Tags     comment
// @Produce  json
// @Param    article_id query int true  "文章 ID"
// @Param    page       query int false "页码"     default(1)
// @Param    page_size  query int false "每页条数"  default(20)
// @Success  200 {object} response.Response{data=appcomment.CommentTreeResp}
// @Router   /api/v1/comments [get]
func (h *CommentHandler) ListComments(c *gin.Context) {
	var req appcomment.ListCommentReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "参数错误："+err.Error())
		return
	}
	req.Page = utils.NormalizePage(req.Page)
	req.PageSize = utils.NormalizePageSize(req.PageSize)

	result, err := h.useCase.ListComments(c.Request.Context(), &req)
	if err != nil {
		response.FailWithError(c, err)
		return
	}
	response.OK(c, result)
}

// CreateComment godoc
// @Summary  发表评论（支持登录用户和游客）
// @Tags     comment
// @Accept   json
// @Produce  json
// @Param    body body appcomment.CreateCommentReq true "评论内容"
// @Success  201  {object} response.Response{data=appcomment.CommentResp}
// @Router   /api/v1/comments [post]
func (h *CommentHandler) CreateComment(c *gin.Context) {
	var req appcomment.CreateCommentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求体解析失败："+err.Error())
		return
	}
	if err := validator.Validate(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// 获取可选用户 ID（未登录时为 nil，走游客流程）
	userID := middleware.GetOptionalUserID(c)
	ip := c.ClientIP()
	ua := c.GetHeader("User-Agent")

	comment, err := h.useCase.CreateComment(c.Request.Context(), &req, userID, ip, ua)
	if err != nil {
		response.FailWithError(c, err)
		return
	}
	response.Created(c, comment)
}

// LikeComment godoc
// @Summary  评论点赞（无需登录）
// @Tags     comment
// @Produce  json
// @Param    id path int true "评论 ID"
// @Success  200 {object} response.Response
// @Router   /api/v1/comments/{id}/like [post]
func (h *CommentHandler) LikeComment(c *gin.Context) {
	id, err := parseCommentID(c)
	if err != nil {
		response.BadRequest(c, "无效的评论 ID")
		return
	}
	if err := h.useCase.LikeComment(c.Request.Context(), id); err != nil {
		response.FailWithError(c, err)
		return
	}
	response.OKWithMsg(c, "点赞成功", nil)
}

// ─── 辅助 ────────────────────────────────────────────────────────────────────

func parseCommentID(c *gin.Context) (uint64, error) {
	return strconv.ParseUint(c.Param("id"), 10, 64)
}
