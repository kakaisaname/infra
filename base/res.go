package base

type ResCode int

const (
	ResCodeOk                 ResCode = 1000
	ResCodeValidationError    ResCode = 2000
	ResCodeRequestParamsError ResCode = 2100
	ResCodeInnerServerError   ResCode = 5000
	ResCodeBizError           ResCode = 6000
)

type Code struct {
	Val int
	Msg string
}

//定义一个struct 返回的数据格式
type Res struct {
	Code    ResCode     `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
