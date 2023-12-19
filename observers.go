package room

type dummyObserver struct{}

func (d dummyObserver) OnError(request IRequest, response IResponse)        {}
func (d dummyObserver) OnAuthError(request IRequest, response IResponse)    {}
func (d dummyObserver) OnAuthComplete(request IRequest, response IResponse) {}
func (d dummyObserver) OnComplete(request IRequest, response IResponse)     {}

type connectionObserver struct{}

func (c connectionObserver) OnDo(request IRequest) {}

func (c connectionObserver) OnRequestCreated(request IRequest) {}

func (c connectionObserver) OnResponseCompleted(request IRequest, response IResponse, err error) {}
