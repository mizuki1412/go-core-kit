package schema

type ResponseBody struct {
	Id      string       `json:"id"`
	Model   string       `json:"model"`
	Object  string       `json:"object"`
	Created int64        `json:"created"` // 时间戳s
	Choices []*ResChoice `json:"choices"`
	Usage   *ResUsage    `json:"usage"`
}

type ResChoice struct {
	Delta        *ResSession `json:"delta"`
	Index        int         `json:"index"`
	FinishReason string      `json:"finish_reason"` // stop
}

type ResSession struct {
	Role             string `json:"role"`
	Content          string `json:"content"`           // 回答
	ReasoningContent string `json:"reasoning_content"` // 推理过程
}

type ResUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}
