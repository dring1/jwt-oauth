package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dring1/jwt-oauth/lib/contextkeys"
	"github.com/dring1/jwt-oauth/logger"
)

const (
	ApacheFormatPattern = "[%s] %s %d %s %s %d %s %f\n"
)

type ApacheLogRecord struct {
	http.ResponseWriter

	ip                    string
	time                  time.Time
	method, uri, protocol string
	status                int
	responseBytes         int64
	elapsedTime           time.Duration
	agent                 string
	id                    string
}

func (r *ApacheLogRecord) Log(logger logger.Service) {
	timeFormatted := r.time.Format("02/Jan/2006 03:04:05")
	requestLine := fmt.Sprintf("%s %s %s", r.method, r.uri, r.protocol)
	//fmt.Fprintf(out, ApacheFormatPattern, timeFormatted, r.id, r.status, r.ip, requestLine, r.responseBytes, r.agent, r.elapsedTime.Seconds())
	logger.Infof(ApacheFormatPattern, timeFormatted, r.id, r.status, r.ip, requestLine, r.responseBytes, r.agent, r.elapsedTime.Seconds())
}

func (r *ApacheLogRecord) Write(p []byte) (int, error) {
	written, err := r.ResponseWriter.Write(p)
	r.responseBytes += int64(written)
	return written, err
}

func (r *ApacheLogRecord) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

type ApacheLoggingHandler struct {
	handler http.Handler
	logger  logger.Service
}

func NewApacheLoggingHandler(logger logger.Service) Middleware {
	return func(handler http.Handler) http.Handler {
		return &ApacheLoggingHandler{
			handler: handler,
			logger:  logger,
		}
	}
}

func (h *ApacheLoggingHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	clientIP := r.RemoteAddr
	if colon := strings.LastIndex(clientIP, ":"); colon != -1 {
		clientIP = clientIP[:colon]
	}

	val := r.Context().Value(contextkeys.ReqId)
	reqId, ok := val.(string)
	if !ok {
		reqId = "INVALID ID"
	}

	record := &ApacheLogRecord{
		ResponseWriter: rw,
		ip:             clientIP,
		time:           time.Time{},
		method:         r.Method,
		uri:            r.RequestURI,
		protocol:       r.Proto,
		status:         http.StatusOK,
		elapsedTime:    time.Duration(0),
		agent:          r.UserAgent(),
		id:             reqId,
	}

	startTime := time.Now()
	h.handler.ServeHTTP(record, r)
	finishTime := time.Now()

	record.time = finishTime.UTC()
	record.elapsedTime = finishTime.Sub(startTime)

	record.Log(h.logger)
}
