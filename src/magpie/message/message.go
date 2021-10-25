package message

// import (
// 	"fmt"
// )

type MessageType int

const (
	EVAL_LINE MessageType = iota
	CALL
	METHOD_CALL
	RETURN
)

type Message struct {
	Type MessageType
	Body interface{}
}

/* All classes that listen to messages must implement this interface. */
type MessageListener interface {
	/**
	 * Called to receive a message sent by a message producer.
	 * @param message the message that was sent.
	 */
	MessageReceived(message Message)
}

type MessageHandler struct {
	message   Message
	listeners map[MessageListener]bool
}

func NewMessageHandler() *MessageHandler {
	return &MessageHandler{listeners: make(map[MessageListener]bool)}
}

func (m *MessageHandler) AddListener(listener MessageListener) {
	m.listeners[listener] = true
}

func (m *MessageHandler) RemoveListener(listener MessageListener) {
	delete(m.listeners, listener)
}

func (m *MessageHandler) SendMessage(message Message) {
	m.message = message
	m.notifyListeners()
}

func (m *MessageHandler) notifyListeners() {
	for l, _ := range m.listeners {
		l.MessageReceived(m.message)
	}
}
