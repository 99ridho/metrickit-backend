package models

type GenericResponse struct {
	Header *GenericResponseHeader `json:"header"`
	Data   interface{}            `json:"data,omitempty"`
	Error  error                  `json:"error,omitempty"`
}

type GenericResponseHeader struct {
	Status   string   `json:"status_code"`
	Messages []string `json:"messages"`
}
