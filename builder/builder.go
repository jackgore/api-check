package builder

type Cookie struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type APITest struct {
	Description string      `json:"description"`
	Method      string      `json:"method"`
	Hostname    string      `json:"hostname"`
	Endpoint    string      `json:"endpoint"`
	Request     APIRequest  `json:"request"`
	Response    APIResponse `json:"response"`
}

type APIRequest struct {
	Body        string            `json:"body"`
	Headers     map[string]string `json:"headers"`
	QueryParams map[string]string `json:"query-params"`
	JSON        interface{}       `json:"json,omitempty"`
	Cookies     []Cookie          `json:"cookies,omitempty"`
}

type APIResponse struct {
	Body       string            `json:"body"`
	Headers    map[string]string `json:"headers"`
	JSON       interface{}       `json:"json,omitempty"`
	StatusCode int               `json:"code"`
}
