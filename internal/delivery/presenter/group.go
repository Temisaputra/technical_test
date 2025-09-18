package presenter

type GroupRequest struct {
	Name        string `json:"name" validate:"required,not_blank"`
	Description string `json:"description"`
	Active      bool   `json:"active"`
}

type GroupResponse struct {
	ID        int    `json:"id"`
	GroupName string `json:"group_name"`
	Value     string `json:"value"`
	Active    bool   `json:"active"`
}
