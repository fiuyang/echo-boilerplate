package exception

type ForbiddenStruct struct {
	ErrorMsg string
}

func NewForbiddenHandler(msg string) *ForbiddenStruct {
	return &ForbiddenStruct{
		ErrorMsg: msg,
	}
}

func (e *ForbiddenStruct) Error() string {
	return e.ErrorMsg
}
