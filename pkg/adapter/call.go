package adapter

type ToolCall struct {
	Id       string           `json:"id" validate:"required"`
	Function ToolCallFunction `json:"function" validate:"required"`
}

type ToolCallFunction struct {
	Name      string `json:"name" validate:"required"`
	Arguments string `json:"arguments" validate:"required"`
}
