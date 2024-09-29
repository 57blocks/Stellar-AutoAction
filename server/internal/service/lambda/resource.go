package lambda

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/57blocks/auto-action/server/internal/dto"
	"github.com/57blocks/auto-action/server/internal/pkg/errorx"
	"github.com/gorilla/websocket"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	r := c.Request

	// Parse the multipart form
	if err := r.ParseMultipartForm(10 << 20); err != nil { // 10 MB max memory
		c.Error(errorx.Internal(fmt.Sprintf("failed to parse multipart form: %s", err.Error())))
		c.Abort()
		return
	}

	if r.MultipartForm == nil {
		c.Error(errorx.Internal("multipart form is nil"))
		c.Abort()
		return
	}

	reqFiles := make([]*dto.ReqFile, 0, len(r.MultipartForm.File))
	for _, headers := range r.MultipartForm.File {
		header := headers[0]

		file, err := header.Open()
		if err != nil {
			c.Error(errorx.Internal(fmt.Sprintf("failed to open file: %s, err: %s", header.Filename, err.Error())))
			c.Abort()
			return
		}

		bytes, err := io.ReadAll(file)
		file.Close()
		if err != nil {
			c.Error(errorx.Internal(fmt.Sprintf("failed to read file: %s, err: %s", header.Filename, err.Error())))
			c.Abort()
			return
		}

		splits := strings.Split(header.Filename, ".")
		name := splits[0]

		reqFiles = append(reqFiles, &dto.ReqFile{
			Name:  name,
			Bytes: bytes,
		})
	}

	resp, err := ServiceImpl.Register(c, &dto.ReqRegister{
		Expression: r.Form.Get("expression"),
		Payload:    r.Form.Get("payload"),
		Files:      reqFiles,
	})
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, resp)
}

func Invoke(c *gin.Context) {
	req := new(dto.ReqInvoke)

	if err := c.BindUri(req); err != nil {
		c.Error(errorx.BadRequest(err.Error()))
		c.Abort()
		return
	}

	if err := c.ShouldBindJSON(req); err != nil {
		c.Error(errorx.BadRequest(err.Error()))
		c.Abort()
		return
	}

	resp, err := ServiceImpl.Invoke(c, req)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, resp)
}

func List(c *gin.Context) {
	queryParams := new(dto.ReqList)
	if err := c.BindQuery(queryParams); err != nil {
		c.Error(errorx.BadRequest(err.Error()))
		c.Abort()
		return
	}

	resp, err := ServiceImpl.List(c, queryParams.Full)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, resp)
}

func Info(c *gin.Context) {
	req := new(dto.ReqURILambda)

	if err := c.BindUri(req); err != nil {
		c.Error(errorx.BadRequest(err.Error()))
		c.Abort()
		return
	}

	resp, err := ServiceImpl.Info(c, req)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, resp)
}

func Logs(c *gin.Context) {
	req := new(dto.ReqURILambda)

	if err := c.BindUri(req); err != nil {
		c.Error(errorx.BadRequest(err.Error()))
		c.Abort()
		return
	}

	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	if err := ServiceImpl.Logs(c, req, &upgrader); err != nil {
		c.Error(err)
		c.Abort()
		return
	}
}

func Remove(c *gin.Context) {
	req := new(dto.ReqURILambda)

	if err := c.BindUri(req); err != nil {
		c.Error(errorx.BadRequest(err.Error()))
		c.Abort()
		return
	}

	resp, err := ServiceImpl.Remove(c, req)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, resp)
}
