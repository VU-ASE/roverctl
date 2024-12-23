package utils

import (
	"fmt"
	"net/http"
)

func ParseHTTPError(err error, htt *http.Response) error {
	if err != nil && htt != nil {
		// Read http response body
		httpRes := make([]byte, htt.ContentLength)
		_, err = htt.Body.Read(httpRes)
		if err != nil {
			return fmt.Errorf("Failed to read http response body: %v", err)
		} else {
			return fmt.Errorf("Operation failed: %v\n%s", err, PrettyJSON(httpRes))
		}
	} else {
		return fmt.Errorf("Operation failed: %v", err)
	}
}
