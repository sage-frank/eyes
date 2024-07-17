package v2

type Responses interface {
	SetCode(int32)
	SetTraceID(string)
	SetMsg(string)
	SetInfo(string)
	SetData(interface{})
	SetSuccess(bool)
	Clone() Responses // 初始化/重置
}
