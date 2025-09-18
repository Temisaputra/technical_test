package presenter

type ContactRequest struct {
	Name        string `json:"name" validate:"required,not_blank"`
	JobPosition string `json:"job_position"`
	Email       string `json:"email" validate:"omitempty,email"`
	Phone       string `json:"phone"`
	Mobile      string `json:"mobile"`
	Main        bool   `json:"main"`
}

type ContactResponse struct {
	ID          int    `json:"id"`
	SupplierID  uint   `json:"supplier_id"`
	Name        string `json:"name"`
	JobPosition string `json:"job_position"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	Mobile      string `json:"mobile"`
	Main        bool   `json:"main"`
}
