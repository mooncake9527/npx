package response

type Responses interface {
	SetCode(int)
	GetCode() int
	SetMsg(string)
	SetData(interface{})
}
