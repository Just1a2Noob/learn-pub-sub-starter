package pubsub

type Acktype int

type SimpleQueueType int

const (
	SimpleQueueDurable SimpleQueueType = iota
	SimpleQueueTransient
)

const (
	QueueClassicType   = "classic"
	QueueClassicQuorum = "quorum"
)

const (
	Ack Acktype = iota
	NackDiscard
	NackRequeue
)
