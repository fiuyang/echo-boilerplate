package exception

type UnauthorizedStruct struct {
	ErrorMsg string
}

func NewUnauthorizedHandler(msg string) *UnauthorizedStruct {
	return &UnauthorizedStruct{
		ErrorMsg: msg,
	}
}

func (e *UnauthorizedStruct) Error() string {
	return e.ErrorMsg
}
