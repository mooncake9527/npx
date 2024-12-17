package response

import (
	"net/http"
	"github.com/gin-gonic/gin"
)


// Error 失败数据处理
func Error(c *gin.Context, code int, msg string) {
	res := &Response{}
	res.SetMsg(msg)
	res.SetCode(code)
	c.Set("result", res)
	c.Set("status", code)
	c.AbortWithStatusJSON(http.StatusOK, res)
}

func CreateResponse(c *gin.Context, code int, msg string) Responses {
	res := &Response{}
	res.SetMsg(msg)
	res.SetCode(code)
	return res
}

// OK 通常成功数据处理
func OK(c *gin.Context, data interface{}, msg string) {
	res := &Response{}
	res.SetData(data)
	if msg != "" {
		res.SetMsg(msg)
	}
	res.SetCode(1)
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
	// data["requestId"] = utils.GetReqId(c)
	c.Set("result", data)
	c.AbortWithStatusJSON(http.StatusOK, data)
}
