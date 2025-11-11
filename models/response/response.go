package response

type Resp struct {
	ReturnCode    int    `json:"return_code"`
	ReturnMessage string `json:"return_message"`
	Data          any    `json:"data,omitempty,omitzero"`
}
