package schema

type RequestBody struct {
	Model         string       `json:"model"`
	Messages      []Message    `json:"messages"`
	Stream        bool         `json:"stream"`
	StreamOptions StreamOption `json:"stream_options"`
}

type StreamOption struct {
	IncludeUsage bool `json:"include_usage"`
}
