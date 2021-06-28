package telegram_notify_service

type DisabledService struct {
}

func (s DisabledService) NotifySupport(m string) {
}

func NewDisabledService() *GenericService {
	var service GenericService = &DisabledService{}
	return &service
}
