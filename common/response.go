package common

// SuccessBody 成功响应结构体
type ResponseBody struct {
	Status   string `json:"status"`
	ErrorMsg string `json:"error"`
	Data     interface{}
}

// Success 成功
func Success(data interface{}) ResponseBody {
	return ResponseBody{
		Status: "200",
		Data:   data,
	}
}

// Error 失败
func Error(errcode, msg string) ResponseBody {
	return ResponseBody{
		Status:   errcode,
		ErrorMsg: msg,
	}
}
