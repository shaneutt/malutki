package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// -----------------------------------------------------------------------------
// Vars & Consts
// -----------------------------------------------------------------------------

var port = 8080

// -----------------------------------------------------------------------------
// Main
// -----------------------------------------------------------------------------

func init() {
	if v := os.Getenv("LISTEN_PORT"); v != "" {
		customPort, err := strconv.Atoi(v)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%d is not a valid port\n", customPort)
			os.Exit(1)
		}
		port = customPort
	}
}

func main() {
	log.Printf("starting server listening on *:%d\n", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), http.HandlerFunc(handler)); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
		os.Exit(1)
	}
}

// -----------------------------------------------------------------------------
// Router
// -----------------------------------------------------------------------------

func handler(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, "/status/") {
		if code := statusCode(r.URL.Path); code != 0 {
			w.WriteHeader(code)
			return
		}
	} else if isPathRoot(r.URL.Path) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, indexHTML)
		return
	}

	w.WriteHeader(http.StatusNotFound)
	return
}

// -----------------------------------------------------------------------------
// Index
// -----------------------------------------------------------------------------

const indexHTML = `<!DOCTYPE html>
<html>
<header>
  <title>Welcome to the Malutki testing server</title>
</header>
<body>
  <h1>Available APIs:</h1>
  <h2><a href="/status/200">/status/</a></h2>
</body>
</html>
`

// -----------------------------------------------------------------------------
// Status Handler
// -----------------------------------------------------------------------------

var statusRegexp = regexp.MustCompile(`^/status/([0-9]+)$`)

func statusCode(path string) int {
	submatches := statusRegexp.FindAllStringSubmatch(path, -1)
	code, _ := strconv.Atoi(submatches[0][1])

	// return not found for invalid status codes
	if http.StatusText(code) == "" {
		return http.StatusNotFound
	}

	// currently, only 2XX, 4XX and 5XX are supported
	if isInRange(code, http.StatusOK, http.StatusIMUsed) ||
		isInRange(code, http.StatusBadRequest, http.StatusUnavailableForLegalReasons) ||
		isInRange(code, http.StatusInternalServerError, http.StatusNetworkAuthenticationRequired) {
		return code
	}

	// if not supported, return "400 Bad Request"
	return http.StatusBadRequest
}

// -----------------------------------------------------------------------------
// Helper Functions
// -----------------------------------------------------------------------------

func isPathRoot(path string) bool {
	if path == "" || path == "/" {
		return true
	}
	return false
}

func isInRange(num, start, end int) bool {
	return num >= start && num <= end
}
