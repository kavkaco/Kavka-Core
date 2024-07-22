package stream

import "errors"

var (
	ErrConsumer = errors.New("infra consumer error")
	ErrProducer = errors.New("infra producer error")
)
