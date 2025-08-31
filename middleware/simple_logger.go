package middleware

import (
	"bytes"
	"io"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

const LOG_FILE = "custom.log"

var logger *log.Logger

func InitLogging() {
	f, err := os.OpenFile(LOG_FILE, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	logger = log.New(f, "", log.Ldate|log.Ltime|log.Lshortfile)
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w bodyLogWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

func CustomLoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		logger.Printf("Started %s %s", c.Request.Method, c.Request.URL.Path)
		bodyAsByteArray, _ := io.ReadAll(c.Request.Body)
		logger.Printf("\tRequestBody:%s", string(bodyAsByteArray))
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyAsByteArray))
		c.Next()
		logger.Printf("Completed %s %s in %v", c.Request.Method,
			c.Request.URL.Path, time.Since(start))
		logger.Printf("\tResponseBody:%s", blw.body.String())
	}
}
