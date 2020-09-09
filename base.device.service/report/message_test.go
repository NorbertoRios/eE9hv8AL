package report

import "testing"

func TestNewMessage(t *testing.T) {
	message := NewMessage()
	if message == nil {
		t.Error("Unable to create new Message")
	}
}

func TestSourceId(t *testing.T) {
	message := NewMessage()
	message.SetSourceID(125)
	if message.SourceID() != 125 {
		t.Error("Invalid source id value")
	}
}

func TestMessageFloating(t *testing.T) {
	message := NewMessage()
	message.SetValue("TestValue", float32(1.0))
	value := message.GetFloatValue("TestValue", float32(0.0))
	if value != float32(1.0) {
		t.Error("Invalid float value")
	}
}
