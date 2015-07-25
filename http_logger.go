package threedo

import (
	"github.com/codegangsta/negroni"
	"github.com/go-kit/kit/log"
	"net/http"
	"os"
	"time"
)

type Logger struct {
	log.Logger
}

func NewLogger() *Logger {
	return &Logger{log.NewLogfmtLogger(os.Stderr)}
}

func (l *Logger) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	start := time.Now()

	next(rw, r)

	res := rw.(negroni.ResponseWriter)
	l.Log(
		"action", "request",
		"method", r.Method,
		"path", r.URL.Path,
		"status", res.Status(),
		"time", time.Since(start),
	)
}
