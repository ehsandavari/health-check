package enums

type HttpMethod string

const (
	HttpMethodGET     HttpMethod = "GET"
	HttpMethodPOST    HttpMethod = "POST"
	HttpMethodPUT     HttpMethod = "PUT"
	HttpMethodDELETE  HttpMethod = "DELETE"
	HttpMethodPATCH   HttpMethod = "PATCH"
	HttpMethodHEAD    HttpMethod = "HEAD"
	HttpMethodOPTIONS HttpMethod = "OPTIONS"
)

func (r HttpMethod) String() string {
	return string(r)
}

func (r HttpMethod) IsValid() bool {
	switch r {
	case HttpMethodGET,
		HttpMethodPOST,
		HttpMethodPUT,
		HttpMethodDELETE,
		HttpMethodPATCH,
		HttpMethodHEAD,
		HttpMethodOPTIONS:
		return true
	default:
		return false
	}
}
