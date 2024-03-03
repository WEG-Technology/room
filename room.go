package room

type IRoom interface {
	Send(request IRequest) IResponse
	PutAuthStrategy(auth IAuthStrategy)
	GetSegment() ISegment
}

type Room struct {
	connection   IConnection
	authStrategy IAuthStrategy
	segment      ISegment
}

func (r *Room) PutAuthStrategy(authStrategy IAuthStrategy) {
	r.authStrategy = authStrategy
}

func (r *Room) Send(request IRequest) IResponse {
	r.segment = StartSegmentNow()
	defer r.segment.End()

	if request.PreRequest() != nil {
		response := r.connection.Send(request.PreRequest())

		if !response.Ok() {
			return response
		}

		if r.authStrategy != nil {
			err := r.authStrategy.Authenticate(request, response)

			if err != nil {
				return response
			}
		}
	}

	response := r.connection.Send(request)

	return response
}

func (r *Room) GetSegment() ISegment {
	return r.segment
}

type OptionRoom func(request *Room)

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

	return r
}
