package handler

type StatusHanlder struct {
}

type StatusService interface {
}

func NewStatusHandler() *StatusHanlder {
	return &StatusHanlder{}
}
