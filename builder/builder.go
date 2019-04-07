package builder

// Cookie represents a cookie that will be sent to the server for an APITest.
// This would typically be used as an authentication method.
type Cookie struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// APITest represents a single test that is to be ran be api-check.
type APITest struct {
	Description string      `json:"description"`
	Method      string      `json:"method"`
	Hostname    string      `json:"hostname"`
	Endpoint    string      `json:"endpoint"`
	Request     APIRequest  `json:"request"`
	Response    APIResponse `json:"response"`
}

// APIRequest describes the HTTP request that will be sent by api-check while
// performing an HTTP request.
type APIRequest struct {
	Body        string            `json:"body"`
	Headers     map[string]string `json:"headers"`
	QueryParams map[string]string `json:"query-params"`
	JSON        interface{}       `json:"json,omitempty"`
	Cookies     []Cookie          `json:"cookies,omitempty"`
}

// APIResponse describes that api-check will expect from the server following
// running an APITest.
type APIResponse struct {
	Body       string            `json:"body"`
	Headers    map[string]string `json:"headers"`
	JSON       interface{}       `json:"json,omitempty"`
	StatusCode int               `json:"code"`
}
