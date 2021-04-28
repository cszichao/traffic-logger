package traffic_logger

import (
	"bytes"
	"encoding/json"
	"github.com/rs/zerolog"
	"github.com/valyala/bytebufferpool"
	"io/ioutil"
	"net/http"
	"time"
)

type Options struct {
	Logger    *zerolog.Logger
	Extractor FieldsExtractor
	Ignore    Ignore
}

type TrafficLogger struct {
	logger    *zerolog.Logger
	extractor FieldsExtractor
	ignore    Ignore
}

func New(options *Options) *TrafficLogger {
	extractor := options.Extractor
	if extractor == nil {
		extractor = &DefaultExtractor{}
	}
	ignore := options.Ignore
	if ignore == nil {
		ignore = &DefaultIgnore{}
	}
	return &TrafficLogger{logger: options.Logger, extractor: extractor, ignore: ignore}
}

func (l *TrafficLogger) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// we don't track traffic without a api name
		apiName := l.extractor.APIName(r)
		if len(apiName) == 0 {
			next.ServeHTTP(w, r)
			return
		}

		// request time
		reqStartTime := time.Now()

		// get request body
		var reqBuffer *bytebufferpool.ByteBuffer
		if !l.ignore.Req(apiName) {
			reqBuffer = bytebufferpool.Get()
			defer bytebufferpool.Put(reqBuffer)
			n, err := reqBuffer.ReadFrom(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			_ = r.Body.Close()
			if n > 0 {
				r.Body = ioutil.NopCloser(bytes.NewReader(reqBuffer.B))
			} else {
				r.Body = &nullReadCloser{}
			}
		}

		// get resp
		var respBuffer *bytebufferpool.ByteBuffer
		if !l.ignore.Resp(apiName) {
			respBuffer = bytebufferpool.Get()
			defer bytebufferpool.Put(respBuffer)
		}
		nw := &recordableResponseWriter{ResponseWriter: w, buffer: respBuffer}

		// wrap
		next.ServeHTTP(nw, r)

		// request log
		logRequestEvent := zerolog.Dict().
			Str("method", r.Method).
			Str("host", l.extractor.Host(r)).
			Str("path", r.URL.Path).
			Str("query", r.URL.RawQuery)
		if reqBuffer != nil {
			logRequestEvent = logBodyEvent(logRequestEvent, reqBuffer.B)
		}

		// response log
		logResponseEvent := zerolog.Dict().Int("status", nw.status)
		if respBuffer != nil {
			logResponseEvent = logBodyEvent(logResponseEvent, respBuffer.B)
		}

		// full log
		l.logger.Log().
			Int64("timestamp", reqStartTime.Unix()).
			Str("api_name", apiName).
			Str("ip", l.extractor.ClientIP(r)).
			Str("operator", l.extractor.Operator(r)).
			Dur("latency", time.Now().Sub(reqStartTime)).
			Dict("request", logRequestEvent).
			Dict("response", logResponseEvent).
			Send()
	})
}

func logBodyEvent(e *zerolog.Event, b []byte) *zerolog.Event {
	if len(b) == 0 {
		return e
	}
	if json.Valid(b) {
		// remove carriage returns
		for i, _b := range b {
			if _b == '\n' || _b == '\r' {
				b[i] = ' '
			}
		}
		return e.RawJSON("body", b)
	}
	return e.Bytes("body", b)
}
