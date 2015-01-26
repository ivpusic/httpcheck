package httpcheck

import (
	"net/http"
)

func cookiesToMap(cookies []*http.Cookie) map[string]string {
	result := map[string]string{}

	for _, cookie := range cookies {
		result[cookie.Name] = cookie.Value
	}

	return result
}
