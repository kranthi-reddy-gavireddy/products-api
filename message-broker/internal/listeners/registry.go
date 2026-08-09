package listeners

type Handler func(message []byte) error

var handlers = map[string]Handler{}

func Register(topic string, h Handler) {
	handlers[topic] = h
}

func GetTopics() []string {
	topics := make([]string, 0, len(handlers))

	for topic := range handlers {
		topics = append(topics, topic)
	}
	return topics
}

func GetHandler(topic string) (Handler, bool) {
	h, exists := handlers[topic]
	return h, exists
}
