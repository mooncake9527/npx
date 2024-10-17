package response

import (
	"net/http"

	"github.com/mooncake9527/npx/common/utils"

	"github.com/gin-gonic/gin"
)

var Default = &response{}

// Error 失败数据处理
func Error(c *gin.Context, code int, msg string) {
	res := Default.Clone()
	res.SetMsg(msg)
	res.SetTraceID(utils.GetReqId(c))
	res.SetCode(code)
	res.SetSuccess(false)
	c.Set("result", res)
	c.Set("status", code)
	c.AbortWithStatusJSON(http.StatusOK, res)
}

func CreateResponse(c *gin.Context, code int, msg string) Responses {
	res := Default.Clone()
	res.SetMsg(msg)
	res.SetTraceID(utils.GetReqId(c))
	res.SetCode(code)
	res.SetSuccess(false)
	return res
}

// OK 通常成功数据处理
func OK(c *gin.Context, data interface{}, msg string) {
	res := Default.Clone()
	res.SetData(data)
	res.SetSuccess(true)
	if msg != "" {
		res.SetMsg(msg)
	}
	res.SetTraceID(utils.GetReqId(c))
	res.SetCode(http.StatusOK)
	c.Set("result", res)
	c.Set("status", http.StatusOK)
	c.AbortWithStatusJSON(http.StatusOK, res)
}

// PageOK 分页数据处理
func PageOK(c *gin.Context, result interface{}, count int64, pageIndex int, pageSize int, msg string) {
	var res page
	res.List = result
	res.Count = count
	res.PageIndex = pageIndex
	res.PageSize = pageSize
	OK(c, res, msg)
}

// Custum 兼容函数
func Custum(c *gin.Context, data gin.H) {
	data["requestId"] = utils.GetReqId(c)
	c.Set("result", data)
	c.AbortWithStatusJSON(http.StatusOK, data)
}
