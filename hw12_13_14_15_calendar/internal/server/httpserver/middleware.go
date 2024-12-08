package httpserver

import (
	"fmt"
	"net"
	"net/http"
	"time"
)

func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		ip, _, _ := net.SplitHostPort(r.RemoteAddr)
		msg := fmt.Sprintf("%s [%s] %s %s %s %d %d %s",
			ip,
			start.Format("02/Jan/2006:15:04:05 -0700"),
			r.Method,
			r.URL.String(),
			r.Proto,
			200,
			time.Since(start).Microseconds(),
			r.UserAgent())
		s.logger.Info(msg)
	})
}
