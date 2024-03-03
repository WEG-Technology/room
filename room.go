package room

type IRoom interface {
	Send(request IRequest) IResponse
	GetSegment() ISegment
}

type IRoomObserver interface {
	OnError(request IRequest, response IResponse)
	OnAuthError(request IRequest, response IResponse)
	OnComplete(request IRequest, response IResponse)
	OnAuthComplete(request IRequest, response IResponse)
}

type Room struct {
	observer     IRoomObserver
	connection   IConnection
	authStrategy IAuthStrategy
	segment      ISegment
}

func (r *Room) Send(request IRequest) IResponse {
	r.segment = StartSegmentNow()
	defer r.segment.End()

	if r.authStrategy != nil {
		response := r.connection.Send(r.authStrategy.GetAuthRequest())

		if !response.Ok() {
			r.observer.OnAuthError(r.authStrategy.GetAuthRequest(), response)
			return response
		}

		err := r.authStrategy.Authenticate(request, response)

		if err != nil {
			r.observer.OnAuthError(r.authStrategy.GetAuthRequest(), response)
			return response
		}

		r.observer.OnAuthComplete(r.authStrategy.GetAuthRequest(), response)
	}

	response := r.connection.Send(request)

	r.observer.OnComplete(request, response)

	return response
}

func (r *Room) GetSegment() ISegment {
	return r.segment
}

type OptionRoom func(request *Room)

func WithObserver(observer IRoomObserver) OptionRoom {
	return func(room *Room) {
		room.observer = observer
	}
}

func WithAuth(auth IAuthStrategy) OptionRoom {
	return func(room *Room) {
		room.authStrategy = auth
	}
}

func NewRoom(connection IConnection, opts ...OptionRoom) *Room {
	r := new(Room)
	r.connection = connection

	for _, opt := range opts {
		opt(r)
	}

	if r.observer == nil {
		r.observer = dummyObserver{}
	}

	return r
}
