package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/rs/xid"
	"github.com/truongnqse05461/ewallet/internal/cx"
	"github.com/truongnqse05461/ewallet/internal/log"
	"github.com/truongnqse05461/ewallet/internal/metrics"
)

type Middleware struct {
	logger log.Logger
	db     *sqlx.DB
	metric *metrics.Metric
}

func NewMiddleware(
	logger log.Logger,
	db *sqlx.DB,
	metric *metrics.Metric,
) *Middleware {
	return &Middleware{
		logger: logger,
		db:     db,
		metric: metric,
	}
}

func (i *Middleware) Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request = c.Request.WithContext(cx.SetLogger(c.Request.Context(), i.logger))
		c.Next()
	}
}
func (i *Middleware) Trace() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := xid.New().String()
		c.Request = c.Request.WithContext(cx.SetTrace(c.Request.Context(), traceID))

		log := cx.GetLogger(c.Request.Context()).WithField("trace_id", traceID)
		c.Request = c.Request.WithContext(cx.SetLogger(c.Request.Context(), log))

		c.Writer.Header().Set("X-Trace-Id", traceID)
		c.Next()
	}
}

func (i *Middleware) Metric() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		i.metric.IncRequestCount(c.Request.Method, c.FullPath(), c.Writer.Status())
		i.metric.ObserveRequestDuration(c.Request.Method, c.FullPath(), c.Writer.Status(), time.Since(start))
	}
}

func (i *Middleware) Access() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		buf := bytes.NewBuffer(nil)
		c.Writer = &ResponseMultiWriter{ResponseWriter: c.Writer, buf: buf}

		bodyField := ""
		if c.Request.Body != http.NoBody {
			body, _ := c.GetRawData()
			c.Request.Body = io.NopCloser(bytes.NewReader(body))
			bodyField = string(body)
		}

		c.Writer.Header().Set("Content-Type", "application/json")
		c.Next()

		fields := log.Fields{
			// request
			"url":            c.Request.URL.String(),
			"headers":        i.getHeaders(c),
			"method":         c.Request.Method,
			"user_agent":     c.Request.UserAgent(),
			"request_body":   bodyField,
			"request_length": len(bodyField),
			// response
			"response_time":   time.Since(start).Milliseconds(),
			"status":          c.Writer.Status(),
			"response_body":   buf.String(),
			"response_length": buf.Len(),
		}
		cx.GetLogger(c.Request.Context()).WithFields(fields).Info("access log")
	}
}

func (i *Middleware) ErrorHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		log := cx.GetLogger(c.Request.Context())

		if err := c.Errors.Last(); err != nil {
			log.WithErr(err).Error("unknown error occurred")
			c.AbortWithStatusJSON(http.StatusInternalServerError, err)
			// if e, ok := err.Err.(errs.AppError); ok {
			// 	log.WithErr(err).WithField("stack_trace", e.StackTrace()).Error("app error occurred")
			// 	c.AbortWithStatusJSON(e.Status, e)
			// } else {
			// 	log.WithErr(err).Error("unknown error occurred")
			// 	c.AbortWithStatusJSON(http.StatusInternalServerError, "internal server error")
			// }
		}
	}
}

func (i *Middleware) Tx() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := cx.GetLogger(c.Request.Context())

		tx, err := i.db.Beginx()
		if err != nil {
			_ = c.Error(err)
			c.Abort()
			return
		}

		defer func() {
			if err := recover(); c.Errors.Last() != nil || err != nil {
				if err := tx.Rollback(); err != nil {
					log.WithErr(err).Error("tx rollback failed")
				} else {
					log.Debug("tx rollbacked")
				}

				if err != nil {
					panic(err)
				}
			}
		}()
		log.Debug("tx started")

		c.Request = c.Request.WithContext(cx.SetTx(c.Request.Context(), tx))
		c.Next()

		if c.Errors.Last() == nil {
			if err := tx.Commit(); err != nil {
				log.WithErr(err).Error("tx commit failed")
				panic(err)
			} else {
				log.Debug("tx committed")
			}
		}
	}
}

func (i *Middleware) getHeaders(c *gin.Context) string {
	headers, err := json.Marshal(c.Request.Header)
	if err != nil {
		return ""
	}
	return string(headers)
}

type ResponseMultiWriter struct {
	gin.ResponseWriter
	buf *bytes.Buffer
}

func (w *ResponseMultiWriter) Write(b []byte) (int, error) {
	w.buf.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w *ResponseMultiWriter) WriteString(s string) (int, error) {
	w.buf.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}
