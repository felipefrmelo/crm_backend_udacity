package crm

type CreateCustomer struct {
	Name      string `json:"name"`
	Role      string `json:"role"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Contacted bool   `json:"contacted"`
}

type UpdateCustomer struct {
	Name      *string `json:"name"`
	Role      *string `json:"role"`
	Email     *string `json:"email"`
	Phone     *string `json:"phone"`
	Contacted *bool   `json:"contacted"`
}
