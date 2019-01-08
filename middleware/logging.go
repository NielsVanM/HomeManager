package middleware

import (
	"net/http"
	"strings"
	"time"

	"nvmtech.nl/homemanager/tools/log"
)

// LogHTTP logs the http request to the console
func LogHTTP(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: Fix long time format in log.txt
		// str := "[" + r.Method + "] " + r.RequestURI + "\tAddress: " + r.RemoteAddr + "\tUser-Agent: " + r.UserAgent() + "\tTimeStamp: " + time.Now().String()
		str := "[" + r.Method + "] " + r.RequestURI + "\tAddress: " + r.RemoteAddr + "\tTimeStamp: " + strings.Split(time.Now().UTC().String(), ".")[0]

		go log.Info("HTTPLog", str)

		next.ServeHTTP(w, r)
	})
}
