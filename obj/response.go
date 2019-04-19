package obj

// Response a server response
type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
