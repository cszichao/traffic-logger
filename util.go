package traffic_logger

import (
	"io"
	"net/http"

	"github.com/valyala/bytebufferpool"
)

// nullReadCloser 假的 request body
type nullReadCloser struct{}

func (r *nullReadCloser) Read([]byte) (int, error) {
	return 0, io.EOF
}

func (r *nullReadCloser) Close() error {
	return nil
}

// recordableResponseWriter 支持 response 记录的 response writer
type recordableResponseWriter struct {
	http.ResponseWriter
	buffer *bytebufferpool.ByteBuffer
	status int
}

func (r *recordableResponseWriter) WriteHeader(statusCode int) {
	r.status = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func (r *recordableResponseWriter) Write(p []byte) (int, error) {
	if r.buffer != nil {
		r.buffer.Write(p)
	}
	return r.ResponseWriter.Write(p)
}
