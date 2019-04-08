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

// JSONType describes the structure of the expected JSON to receive.
type JSONType = map[string]interface{}

// JSONArray describes the structure of an array of expected JSON.
type JSONArrayOf struct {
    Value interface{}
}

// APIResponse describes that api-check will expect from the server following
// running an APITest.
type APIResponse struct {
    // Describes the exact body of the HTTP response that should be received
    // from the server.
	Body       string            `json:"body"`

    // Describes the exact JSON expected from the server.
	JSON       interface{}       `json:"json,omitempty"`

    // TypeOf describes what type should be expected from the server. Instead
    // of specifying exactly what should be received it describes the structure
    // of what should be received.
    TypeOf     *interface{}          `json:"ofType,omitempty"`

    // Array describes the structure of the expected response in the 
    ArrayOf    *interface{}       `json:"arrayOf,omitempty"`

    // Describes the headers that are expected to be received from the server.
	Headers    map[string]string `json:"headers"`

    // Describes the status code expected from the server.
	StatusCode int               `json:"code"`
}
