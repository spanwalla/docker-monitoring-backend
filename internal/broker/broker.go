package broker

//go:generate mockgen -source=broker.go -destination=mocks/mock.go

type Publisher interface {
	Publish(body []byte) error
}

type Consumer interface {
	Start() error
	Stop() error
}
