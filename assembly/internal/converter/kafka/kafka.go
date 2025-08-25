package kafka

import "github.com/space-wanderer/microservices/assembly/internal/model"

type AssemblyRecodedDecoder interface {
	Decode(data []byte) model.OrderPaidEvent
}
