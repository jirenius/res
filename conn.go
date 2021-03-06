package res

import nats "github.com/nats-io/nats.go"

// Conn is an interface that represents a connection to a NATS server.
// It is implemented by nats.Conn.
type Conn interface {
	// Publish publishes the data argument to the given subject
	Publish(subject string, payload []byte) error

	// PublishRequest publishes a request expecting a response on the reply
	// subject.
	PublishRequest(subject, reply string, data []byte) error

	// ChanSubscribe subscribes to messages matching the subject pattern.
	ChanSubscribe(subject string, ch chan *nats.Msg) (*nats.Subscription, error)

	// ChanQueueSubscribe subscribes to messages matching the subject pattern.
	// All subscribers with the same queue name will form the queue group and
	// only one member of the group will be selected to receive any given
	// message, which will be placed on the channel.
	ChanQueueSubscribe(subject, queue string, ch chan *nats.Msg) (*nats.Subscription, error)

	// Close will close the connection to the server.
	Close()
}
