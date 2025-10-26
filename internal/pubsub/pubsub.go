package pubsub

type Acktype int

type SimpleQueueType struct {
	Durable   bool
	Transient bool
}

const (
	Ack Acktype = iota
	NackDiscard
	NackRequeue
)
