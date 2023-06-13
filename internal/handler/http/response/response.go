package response

type ResponseType struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func ErrorResponse(msg string) ResponseType {
	return ResponseType{
		Success: false,
		Message: msg,
	}
}
