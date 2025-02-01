package crm

import "github.com/google/uuid"

type InMemoryCustomerRepository struct {
	customers []Customer
}

func (c *InMemoryCustomerRepository) GetCustomers() []Customer {
	return c.customers
}

func (c *InMemoryCustomerRepository) GetCustomerByID(id string) (*Customer, error) {

	for i := range c.customers {
		if c.customers[i].ID == id {
			return &c.customers[i], nil
		}
	}

	return nil, &ResponseError{NotFoundErrorMsg}
}

func (c *InMemoryCustomerRepository) AddCustomer(createCustomer CreateCustomer) (*Customer, error) {
	customer := Customer{
		ID:        uuid.New().String(),
		Name:      createCustomer.Name,
		Role:      createCustomer.Role,
		Email:     createCustomer.Email,
		Phone:     createCustomer.Phone,
		Contacted: createCustomer.Contacted,
	}
	c.customers = append(c.customers, customer)

	return &customer, nil
}

func (c *InMemoryCustomerRepository) UpdateCustomer(id string, updateCustomer UpdateCustomer) (*Customer, error) {

	customer, err := c.GetCustomerByID(id)

	if err != nil {
		return nil, &ResponseError{NotFoundErrorMsg}
	}

	if updateCustomer.Name != nil {
		customer.Name = *updateCustomer.Name
	}
	if updateCustomer.Role != nil {
		customer.Role = *updateCustomer.Role
	}
	if updateCustomer.Email != nil {
		customer.Email = *updateCustomer.Email
	}
	if updateCustomer.Phone != nil {
		customer.Phone = *updateCustomer.Phone
	}
	if updateCustomer.Contacted != nil {
		customer.Contacted = *updateCustomer.Contacted
	}
	return customer, nil
}

func (c *InMemoryCustomerRepository) DeleteCustomer(id string) error {

	for i, customer := range c.customers {
		if customer.ID == id {
			c.customers = append(c.customers[:i], c.customers[i+1:]...)
			return nil
		}
	}

	return &ResponseError{NotFoundErrorMsg}
}

func NewRepo() *InMemoryCustomerRepository {

	customers := []Customer{{
		ID:        "867836be-ba5f-4877-b7a7-27fa26a1f6da",
		Name:      "Test",
		Role:      "enginner",
		Email:     "test@test.com",
		Phone:     "18 2818312301",
		Contacted: true,
	}}

	return &InMemoryCustomerRepository{customers}

}
