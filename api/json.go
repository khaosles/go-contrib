package api

/*
   @File: json.go
   @Author: khaosles
   @Time: 2023/3/7 21:54
   @Desc:
*/

type JsonResult struct {
	// code
	Code int `json:"code" default:"0"`
	// response information
	Msg string `json:"msg" default:""`
	// data
	Data any `json:"data,omitempty" default:"nil"` // 默认无数据时不显示该字段
	// whether success
	Success bool `json:"success" default:"false"`
}

func NewYes(data any) *JsonResult {
	return &JsonResult{
		Code:    StatusOK,
		Msg:     "ok",
		Data:    data,
		Success: true,
	}
}

func NewNo(code int, msg string) *JsonResult {
	return &JsonResult{
		Code:    code,
		Msg:     msg,
		Data:    nil,
		Success: false,
	}
}

func (j *JsonResult) SetCode(code int) *JsonResult {
	j.Code = code
	return j
}

func (j *JsonResult) SetMsg(msg string) *JsonResult {
	j.Msg = msg
	return j
}

func (j *JsonResult) SetSuccess(success bool) *JsonResult {
	j.Success = success
	return j
}

func (j *JsonResult) SetData(data any) *JsonResult {
	j.Data = data
	return j
}
