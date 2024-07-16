package exception

type BadRequestStruct struct {
	ErrorMsg string
}

func NewBadRequestHandler(msg string) *BadRequestStruct {
	return &BadRequestStruct{
		ErrorMsg: msg,
	}
}

func (e *BadRequestStruct) Error() string {
	return e.ErrorMsg
}
