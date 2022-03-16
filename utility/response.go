package utility

type JsonResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Error   error       `json:"error"`
}
