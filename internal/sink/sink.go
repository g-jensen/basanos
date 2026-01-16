package sink

type Sink interface {
	Emit(event any) error
}
