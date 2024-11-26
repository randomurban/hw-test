package internalhttp

import (
	"fmt"
	"net"
	"net/http"
	"time"
)

// Logging: 66.249.65.3 [25/Feb/2020:19:11:24 +0600] GET /hello?q=1 HTTP/1.1 200 30 "Mozilla/5.0"
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
