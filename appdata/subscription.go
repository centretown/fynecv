package appdata

import "encoding/json"

type Subscription struct {
	consumer Consumer
	run      func(consumer Consumer)
}

func NewSubcription(consumer Consumer, run func(Consumer)) *Subscription {
	sub := &Subscription{
		consumer: consumer,
		run:      run,
	}
	return sub
}

func (sub *Subscription) Consume(newState *Entity[json.RawMessage]) {
	sub.consumer.Copy(newState)
	sub.run(sub.consumer)
}
