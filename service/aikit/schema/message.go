package schema

const (
	// Assistant is the role of an assistant, means the message is returned by ChatModel.
	Assistant string = "assistant"
	// User is the role of a user, means the message is a user message.
	User string = "user"
	// System is the role of a system, means the message is a system message.
	System string = "system"
	// Tool is the role of a tool, means the message is a tool call output.
	Tool string = "tool"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
