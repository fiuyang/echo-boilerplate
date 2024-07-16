package exception

type NotFoundStruct struct {
	ErrorMsg string
}

func NewNotFoundHandler(msg string) *NotFoundStruct {
	return &NotFoundStruct{
		ErrorMsg: msg,
	}
}

func (e *NotFoundStruct) Error() string {
	return e.ErrorMsg
}
