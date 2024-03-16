package traffic_logger

import (
	"github.com/gin-gonic/gin"
	"github.com/valyala/bytebufferpool"
	"net/http"
	"time"
)

// recordableGinResponseWriter 支持 response 记录的 response writer
type recordableGinResponseWriter struct {
	gin.ResponseWriter
	buffer *bytebufferpool.ByteBuffer
	status int
}

func (r *recordableGinResponseWriter) WriteHeader(statusCode int) {
	r.status = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func (r *recordableGinResponseWriter) Write(p []byte) (int, error) {
	if r.buffer != nil {
		r.buffer.Write(p)
	}
	return r.ResponseWriter.Write(p)
}

func (l *TrafficLogger) Gin() gin.HandlerFunc {
	return func(c *gin.Context) {
		// we don't track traffic without api name
		apiName := l.extractor.APIName(c.Request)
		if len(apiName) == 0 {
			c.Next()
			return
		}

		// request time
		reqStartTime := time.Now()

		// get request body
		var reqBuffer *bytebufferpool.ByteBuffer
		if !l.ignore.Req(apiName) {
			reqBuffer = bytebufferpool.Get()
			defer bytebufferpool.Put(reqBuffer)
			if err := l.getRequestBody(c.Request, reqBuffer); err != nil {
				_ = c.AbortWithError(http.StatusInternalServerError, err)
				return
			}
		}

		// get resp
		var respBuffer *bytebufferpool.ByteBuffer
		if !l.ignore.Resp(apiName) {
			respBuffer = bytebufferpool.Get()
			defer bytebufferpool.Put(respBuffer)
		}
		nw := &recordableGinResponseWriter{ResponseWriter: c.Writer, buffer: respBuffer}

		// wrap
		c.Writer = nw
		c.Next()

		l.performLogging(reqStartTime, apiName, c.Request, nw.status, reqBuffer, respBuffer)
	}
}
