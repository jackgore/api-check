package builder

type Cookie struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type APITest struct {
	Description string      `json:"description"`
	Hostname    string      `json:"hostname"`
	Endpoint    string      `json:"endpoint"`
	Method      string      `json:"method"`
	Request     APIRequest  `json:"request"`
	Response    APIResponse `json:"response"`
}

type APIRequest struct {
	Body        string            `json:"body"`
	Headers     map[string]string `json:"headers"`
	JSON        interface{}       `json:"json"`
	QueryParams map[string]string `json:"query-params"`
	Cookies     []Cookie          `json:"cookies"`
}

type APIResponse struct {
	Body       string            `json:"body"`
	Headers    map[string]string `json:"headers"`
	JSON       interface{}       `json:"json"`
	StatusCode int               `json:"code"`
}
