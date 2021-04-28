package traffic_logger

import (
	"net"
	"net/http"
	"strings"
)

type FieldsExtractor interface {
	Host(r *http.Request) string
	ClientIP(r *http.Request) string
	APIName(r *http.Request) string
	Operator(r *http.Request) string
}

type DefaultExtractor struct{}

func (DefaultExtractor) Host(r *http.Request) string {
	if host := r.Header.Get("X-Forwarded-Host"); len(host) > 0 {
		return host
	}
	return r.Host
}

func (DefaultExtractor) ClientIP(r *http.Request) string {
	if clientIP := strings.TrimSpace(r.Header.Get("X-Real-Ip")); clientIP != "" {
		return clientIP
	}
	clientIP := r.Header.Get("X-Forwarded-For")
	if clientIP = strings.TrimSpace(strings.Split(clientIP, ",")[0]); clientIP != "" {
		return clientIP
	}
	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		return ip
	}
	return ""
}

func (DefaultExtractor) APIName(r *http.Request) string {
	return r.Header.Get("X-Api-Name")
}

func (DefaultExtractor) Operator(r *http.Request) string {
	return r.Header.Get("X-Forwarded-User-Name")
}
