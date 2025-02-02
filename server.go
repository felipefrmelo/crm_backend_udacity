package crm

import (
	"net/http"
)

type CustomerRepository interface {
	GetCustomers() []Customer
	GetCustomerByID(id string) (*Customer, error)
	AddCustomer(customer CreateCustomer) (*Customer, error)
	UpdateCustomer(id string, customer UpdateCustomer) (*Customer, error)
	DeleteCustomer(id string) error
}

type HttpEngine interface {
	Params(key string, d ...string) string
	Status(code int) HttpEngine
	JSON(data interface{}) error
	BodyParser(out interface{}) error
	Render(name string, bind interface{}) error
}

type AppCrm struct {
	repo CustomerRepository
	HttpEngine
}

func (s *AppCrm) GetCustomers(c HttpEngine) error {
	return c.Status(http.StatusOK).JSON(s.repo.GetCustomers())
}

func (s *AppCrm) Home(c HttpEngine) error {
	return c.Render("index", map[string]string{"Title": "Test, World!"})
}

func (s *AppCrm) GetCustomerByID(c HttpEngine) error {
	id := c.Params("id")
	customer, err := s.repo.GetCustomerByID(id)

	if err != nil {
		return c.Status(http.StatusNotFound).JSON(err)
	}

	return c.Status(http.StatusOK).JSON(customer)
}

func (s *AppCrm) AddCustomer(c HttpEngine) error {
	var customer CreateCustomer

	if err := c.BodyParser(&customer); err != nil {
		return c.Status(http.StatusBadRequest).JSON(&ResponseError{err.Error()})
	}

	s.repo.AddCustomer(customer)

	return c.Status(http.StatusCreated).JSON(customer)
}

func (s *AppCrm) UpdateCustomer(c HttpEngine) error {
	var customer UpdateCustomer

	if err := c.BodyParser(&customer); err != nil {
		return c.Status(http.StatusBadRequest).JSON(&ResponseError{err.Error()})
	}

	id := c.Params("id")

	customerUpdated, err := s.repo.UpdateCustomer(id, customer)

	if err != nil {
		return c.Status(http.StatusNotFound).JSON(err)
	}

	return c.Status(http.StatusOK).JSON(customerUpdated)
}

func (s *AppCrm) DeleteCustomer(c HttpEngine) error {
	id := c.Params("id")

	err := s.repo.DeleteCustomer(id)

	if err != nil {
		return c.Status(http.StatusNotFound).JSON(err)
	}

	return c.Status(http.StatusNoContent).JSON(nil)
}
