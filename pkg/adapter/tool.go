package adapter

type Tool struct {
	Function ToolFunction `json:"function" validate:"required"`
}

type ToolFunction struct {
	Name        string         `json:"name" validate:"required"`
	Description *string        `json:"description,omitempty"`
	Parameters  map[string]any `json:"parameters" validate:"required"`
	Strict      *bool          `json:"strict,omitempty"`
}
