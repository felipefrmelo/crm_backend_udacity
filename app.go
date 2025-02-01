package crm

import (
	"net/http"
)

type HandlerHttpEngine func(HttpEngine) error

type Server interface {
	Get(path string, handler HandlerHttpEngine)
	Post(path string, handler HandlerHttpEngine)
	Put(path string, handler HandlerHttpEngine)
	Delete(path string, handler HandlerHttpEngine)
	Listen(addr string) error
	Test(req *http.Request, msTimeout ...int) (*http.Response, error)
}

func NewServer(repo CustomerRepository, server ...string) Server {

	var serverName string
	if len(server) == 0 {
		serverName = "fiber"
	} else {
		serverName = server[0]
	}

	app := FactoryServer(serverName)

	crmApp := new(AppCrm)
	crmApp.repo = repo

	app.Get("/customers", crmApp.GetCustomers)
	app.Get("/customers/:id", crmApp.GetCustomerByID)
	app.Post("/customers", crmApp.AddCustomer)
	app.Put("/customers/:id", crmApp.UpdateCustomer)
	app.Delete("/customers/:id", crmApp.DeleteCustomer)

	return app
}
