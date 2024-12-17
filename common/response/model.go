package response

type Response struct {
	Code int         `protobuf:"varint,2,opt,name=code,proto3" json:"code" example:"1"`
	Msg  string      `protobuf:"bytes,3,opt,name=message,proto3" json:"message" example:"ok"`
	Data interface{} `json:"data" example:"null"`
}

type Page struct {
	Count int64 `json:"count"`
	// PageIndex int   `json:"pageIndex"`
	// PageSize  int   `json:"pageSize"`
}

type page struct {
	Page
	List interface{} `json:"list"`
}

func (e *Response) SetData(data interface{}) {
	e.Data = data
}

func (e *Response) SetMsg(s string) {
	e.Msg = s
}

func (e *Response) SetCode(code int) {
	e.Code = code
}

func (e *Response) GetCode() int {
	return e.Code
}
