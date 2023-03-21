package httpcheck

import "fmt"

type FormData struct {
	Key   string
	Value string
	// Option
	FileName string
}

// String to string
//
// --HTTPCheckerBoundery
// Content-Disposition: form-data; name="${Key}"
//
// ${Value}
func (f FormData) String() string {
	options := ""
	if len(f.FileName) > 0 {
		options = fmt.Sprintf("filename=\"%s\"\nContent-Type: text/csv", f.FileName)
	}
	return fmt.Sprintf(`--HTTPCheckBoundary
	Content-Disposition: form-data; name="%s"; %s
	
	%s
	`, f.Key, options, f.Value)
}

func (tt *Tester) WithMultipart(items ...FormData) *Tester {
	tt.WithHeader(
		"Content-Type",
		"multipart/form-data; boundary=HTTPCheckBoundary",
	)
	payload := ""
	for _, item := range items {
		payload += item.String()
	}
	payload += "--HTTPCheckBoundary--"
	print(payload)
	return tt.WithBody([]byte(payload))
}
