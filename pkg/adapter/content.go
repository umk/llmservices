package adapter

type ContentPartText struct {
	Text string `json:"text" validate:"required"`
}

type ContentPartImage struct {
	ImageUrl string `json:"image_url" validate:"required,url"`
}

type ContentPart struct {
	OfContentPartText     *ContentPartText  `json:"text,omitempty"`
	OfContentPartImageUrl *ContentPartImage `json:"image_url,omitempty"`
}
