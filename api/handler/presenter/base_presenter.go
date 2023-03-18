package presenter

import "github.com/google/go-cmp/cmp"

type MetaResponse struct {
	Code           int    `json:"code,omitempty"`
	Message        string `json:"message,omitempty"`
	Total          uint64 `json:"total,omitempty"`
	NextCursor     string `json:"next_cursor,omitempty"`
	PreviousCursor string `json:"previous_cursor,omitempty"`
}

type Response struct {
	Meta   MetaResponse   `json:"meta,omitempty"`
	Text   interface{}    `json:"text,omitempty"`
	Data   interface{}    `json:"data,omitempty"`
	Errors ErrorResponses `json:"errors,omitempty"`
}

// IsEmpty check if the struct is empty or not
func (r Response) IsEmpty() bool {
	return cmp.Equal(r, Response{})
}
