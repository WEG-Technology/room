package room

type connectionObserver struct{}

func (c connectionObserver) OnDo(request IRequest) {}

func (c connectionObserver) OnRequestCreated(request IRequest) {}

func (c connectionObserver) OnResponseCompleted(request IRequest, response IResponse, err error) {}
