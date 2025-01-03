package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"
	xcalistores3 "xcalistore-s3"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)

type server struct {
	listener net.Listener
	config   options
}

var s = server{
	listener: nil,
	config: options{
		8080,
		[]passwordCredentials{{
			Username: "peter",
			Password: "pass",
		}},
	},
}

func main() {
	var err error
	s.listener, err = net.Listen("tcp", fmt.Sprintf(":%d", s.config.port))
	if err != nil {
		panic(fmt.Sprintf("Error while starting to listen at an ephemeral port: %v", err))
	}

	ctx := context.Background()

	store, storeErr := xcalistores3.NewStore(ctx, "test-xcali-backend")
	if storeErr != nil {
		panic(fmt.Sprintf("failed to created S3 store: %v", storeErr))
	}

	startServer(store)
}

func RequestLogger(g *gin.Context) {
	start := time.Now()

	l := getLogger().With().Str("req_xid", xid.New().String()).Logger()

	r := g.Request
	g.Request = r.WithContext(l.WithContext(r.Context()))

	lrw := newLoggingResponseWriter(g.Writer)

	defer func() {
		panicVal := recover()
		if panicVal != nil {
			lrw.statusCode = http.StatusInternalServerError // ensure that the status code is updated
			panic(panicVal)                                 // continue panicking
		}
		l.
			Info().
			Str("method", g.Request.Method).
			Str("url", g.Request.URL.RequestURI()).
			Str("user_agent", g.Request.UserAgent()).
			Int("status_code", lrw.statusCode).
			Dur("elapsed_ms", time.Since(start)).
			Msg("incoming request")
	}()

	g.Next()
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}
