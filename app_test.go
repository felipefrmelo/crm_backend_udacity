package crm_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/felipefrmelo/crm_backend_udacity"
	"github.com/google/uuid"
)

type FakeCustomerRepository struct {
	customers []crm.Customer
}

func (c *FakeCustomerRepository) GetCustomers() []crm.Customer {
	return c.customers
}

func (c *FakeCustomerRepository) GetCustomerByID(id string) (*crm.Customer, error) {

	for i := range c.customers {
		if c.customers[i].ID == id {
			return &c.customers[i], nil
		}
	}

	return nil, &crm.ResponseError{crm.NotFoundErrorMsg}
}

func (c *FakeCustomerRepository) AddCustomer(createCustomer crm.CreateCustomer) (*crm.Customer, error) {
	customer := crm.Customer{
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

func (c *FakeCustomerRepository) UpdateCustomer(id string, updateCustomer crm.UpdateCustomer) (*crm.Customer, error) {

	customer, err := c.GetCustomerByID(id)

	if err != nil {
		return nil, &crm.ResponseError{crm.NotFoundErrorMsg}
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

func (c *FakeCustomerRepository) DeleteCustomer(id string) error {

	for i := range c.customers {
		if c.customers[i].ID == id {
			c.customers = append(c.customers[:i], c.customers[i+1:]...)
			return nil
		}
	}
	return &crm.ResponseError{crm.NotFoundErrorMsg}
}

func NewRepo() *FakeCustomerRepository {

	customers := []crm.Customer{{
		ID:        "867836be-ba5f-4877-b7a7-27fa26a1f6da",
		Name:      "Test",
		Role:      "enginner",
		Email:     "test@test.com",
		Phone:     "18 2818312301",
		Contacted: true,
	}}

	return &FakeCustomerRepository{customers}

}

func TestGetCustomers(t *testing.T) {

	t.Run("test get customer return 200", func(t *testing.T) {

		req := getCustomersRequest()

		repo := NewRepo()

		app := crm.NewServer(repo)

		resp, _ := app.Test(req)

		assertStatusCode(t, resp.StatusCode, http.StatusOK)
		assertGetCustomerBody(t, resp, repo)
	})

}

func getCustomersRequest() *http.Request {
	req := httptest.NewRequest(http.MethodGet, "/customers", nil)
	return req
}

func TestGetCustomerByID(t *testing.T) {

	t.Run("test get customer by id return 200", func(t *testing.T) {

		customerID := "867836be-ba5f-4877-b7a7-27fa26a1f6da"
		req := getCustomerRequest(customerID)

		repo := NewRepo()

		app := crm.NewServer(repo)

		resp, _ := app.Test(req)

		assertStatusCode(t, resp.StatusCode, http.StatusOK)
		assertGetCustomerByID(t, resp, customerID)
	})

	t.Run("test return 404 when id not found", func(t *testing.T) {

		customerID := "042174cd-7c03-4adf-8a08-e5f1180b4c55"
		req := httptest.NewRequest(http.MethodGet, "/customers/"+customerID, nil)

		repo := NewRepo()

		app := crm.NewServer(repo)

		resp, _ := app.Test(req)

		assertStatusCode(t, resp.StatusCode, http.StatusNotFound)

		got := parseResponseBody[crm.ResponseError](t, resp)

		want := crm.NotFoundErrorMsg
		if got.Error() != want {
			t.Errorf("want error  %v got = %v", got.Error(), want)
		}
	})

}

func TestAddCustomer(t *testing.T) {

	t.Run("test add customer return 201", func(t *testing.T) {

		customer := crm.CreateCustomer{
			Name:      "Test",
			Role:      "enginner",
			Email:     "test@test.com",
			Phone:     "18 2818312301",
			Contacted: true,
		}

		req := createRequest(customer)

		repo := NewRepo()
		app := crm.NewServer(repo)
		resp, _ := app.Test(req)
		assertStatusCode(t, resp.StatusCode, http.StatusCreated)
		got := parseResponseBody[crm.Customer](t, resp)
		if got.Name != customer.Name {
			t.Errorf("want customer  %v got = %v", customer.Name, got)
		}

		req = getCustomersRequest()

		resp, _ = app.Test(req)
		assertStatusCode(t, resp.StatusCode, http.StatusOK)
		customers := parseResponseBody[[]crm.Customer](t, resp)

		if len(customers) != 2 {
			t.Errorf("want 2 customers got %v", len(customers))
		}

	})
}

func TestUpdateCustomer(t *testing.T) {

	t.Run("test update customer return 200", func(t *testing.T) {

		repo := NewRepo()
		app := crm.NewServer(repo)

		req := getCustomersRequest()

		resp, _ := app.Test(req)
		customers := parseResponseBody[[]crm.Customer](t, resp)

		customerID := customers[0].ID

		role := "new role"
		payload := crm.UpdateCustomer{
			Role: &role,
		}

		req = updateRequest(payload, customerID)

		resp, _ = app.Test(req)

		assertStatusCode(t, resp.StatusCode, http.StatusOK)

		customer, _ := repo.GetCustomerByID(customerID)

		if customer.Role != *payload.Role {
			t.Errorf("want customer  %v got = %v", payload.Role, customer.Role)
		}

		want := "Test"
		if customer.Name != want {
			t.Errorf("want customer  %v got = %v", want, customer.Name)
		}

	})
}

func TestDeleteCustomer(t *testing.T) {

	t.Run("test delete customer return 200", func(t *testing.T) {
		repo := NewRepo()
		app := crm.NewServer(repo)

		req := getCustomersRequest()

		resp, _ := app.Test(req)

		customers := parseResponseBody[[]crm.Customer](t, resp)

		customerID := customers[0].ID

		req = httptest.NewRequest(http.MethodDelete, "/customers/"+customerID, nil)

		resp, _ = app.Test(req)

		assertStatusCode(t, resp.StatusCode, http.StatusNoContent)

		req = getCustomersRequest()

		resp, _ = app.Test(req)

		customers = parseResponseBody[[]crm.Customer](t, resp)

		if len(customers) != 0 {
			t.Errorf("want 0 customers got %v", len(customers))
		}

	})
}

func updateRequest(customer crm.UpdateCustomer, customerID string) *http.Request {
	payload, _ := json.Marshal(customer)
	req := httptest.NewRequest(http.MethodPut, "/customers/"+customerID, bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	return req
}

func assertGetCustomerByID(t *testing.T, resp *http.Response, customerID string) {
	got := parseResponseBody[crm.Customer](t, resp)
	if got.ID != customerID {
		t.Errorf("want customer  %v got = %v", customerID, got)
	}
}

func assertGetCustomerBody(t *testing.T, resp *http.Response, repo *FakeCustomerRepository) {
	t.Helper()
	got := parseResponseBody[[]crm.Customer](t, resp)

	if len(got) != len(repo.GetCustomers()) {
		t.Errorf("Want  customers got %v", got)
	}

	want := repo.GetCustomers()[0]
	if got[0] != want {
		t.Errorf("Want a customer %v got %v", want, got[0])
	}
}

func parseResponseBody[T any](t *testing.T, resp *http.Response) T {
	t.Helper()
	var result T
	err := json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		t.Errorf("error in parsing the body: %v", err)
	}
	return result
}

func assertStatusCode(t *testing.T, got int, want int) {
	t.Helper()
	if got != want {
		t.Errorf("did not get correct status, got %d, want %d", got, want)
	}
}

func getCustomerRequest(customerID string) *http.Request {
	req := httptest.NewRequest(http.MethodGet, "/customers/"+customerID, nil)
	return req
}

func createRequest(customer crm.CreateCustomer) *http.Request {
	payload, _ := json.Marshal(customer)
	req := httptest.NewRequest(http.MethodPost, "/customers", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	return req
}
