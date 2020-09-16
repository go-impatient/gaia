package model

// Data is the parent model for returning data in this api,
// includes meta for pagination
type Data struct {
	Data       interface{}     `json:"data"`
	Pagination *PaginationData `json:"pagination,omitempty"`
}

// PaginationData ...
type PaginationData struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
	Total int `json:"total"`
}

// Response ...
type Response struct {
	Status     *int             `json:"status"`
	StatusDesc *string          `json:"status_desc,"`
	Errors     *[]string        `json:"errors,omitempty"`
	Data       map[string]*Data `json:"data,omitempty"`
}

// ErrorResponseType 定义了标准的 API 接口错误时返回数据模型
type ErrorResponseType struct {
	Error     string `json:"error"`
	ErrorCode int    `json:"error_code"`
	Message   string `json:"message"`
}

type ResponseBody struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

// SuccessResponseType 定义了标准的 API 接口成功时返回数据模型
type SuccessResponseType struct {
	TotalSize  int         `json:"totalSize,omitempty"`
	PageNumber string      `json:"pageNumber,omitempty"`
	Result     interface{} `json:"result"`
}

// ResponseType ...
type ResponseType struct {
	ErrorResponseType
	SuccessResponseType
}

// BoolRes ...
type BoolRes struct {
	SuccessResponseType
	Result bool `json:"result"`
}

// StringRes ...
type StringRes struct {
	SuccessResponseType
	Result string `json:"result"`
}
