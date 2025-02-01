package crm

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"

	"github.com/gorilla/mux"

	"github.com/gofiber/fiber/v2"
)

type FiberHttpEngineAdapter struct {
	c *fiber.Ctx
}

func (f *FiberHttpEngineAdapter) JSON(data interface{}) error {
	return f.c.JSON(data)
}

func (f *FiberHttpEngineAdapter) Status(code int) HttpEngine {
	f.c.Status(code)
	return f
}

func (f *FiberHttpEngineAdapter) Params(key string, d ...string) string {
	return f.c.Params(key)
}

func (f *FiberHttpEngineAdapter) BodyParser(out interface{}) error {
	return f.c.BodyParser(out)
}

type FiberServerAdapter struct {
	app *fiber.App
}

func (f *FiberServerAdapter) Get(path string, handler HandlerHttpEngine) {
	f.app.Get(path, func(c *fiber.Ctx) error {
		return handler(&FiberHttpEngineAdapter{c})
	})
}

func (f *FiberServerAdapter) Post(path string, handler HandlerHttpEngine) {
	f.app.Post(path, func(c *fiber.Ctx) error {
		return handler(&FiberHttpEngineAdapter{c})
	})
}

func (f *FiberServerAdapter) Put(path string, handler HandlerHttpEngine) {
	f.app.Put(path, func(c *fiber.Ctx) error {
		return handler(&FiberHttpEngineAdapter{c})
	})
}

func (f *FiberServerAdapter) Delete(path string, handler HandlerHttpEngine) {
	f.app.Delete(path, func(c *fiber.Ctx) error {
		return handler(&FiberHttpEngineAdapter{c})
	})
}

func (f *FiberServerAdapter) Listen(addr string) error {
	return f.app.Listen(addr)
}

func (f *FiberServerAdapter) Test(req *http.Request, msTimeout ...int) (*http.Response, error) {
	return f.app.Test(req, msTimeout...)
}

func NewFiberServerAdapter() *FiberServerAdapter {
	return &FiberServerAdapter{
		app: fiber.New(),
	}
}

// StdHttpEngineAdapter

type GorillaHttpEngineAdapter struct {
	w http.ResponseWriter
	r *http.Request
}

func (g *GorillaHttpEngineAdapter) JSON(data interface{}) error {
	g.w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(g.w).Encode(data)
}

func (g *GorillaHttpEngineAdapter) Status(code int) HttpEngine {
	g.w.WriteHeader(code)
	return g
}

func (g *GorillaHttpEngineAdapter) Params(key string, d ...string) string {
	vars := mux.Vars(g.r)
	if val, ok := vars[key]; ok {
		return val
	}
	if len(d) > 0 {
		return d[0]
	}
	return ""
}

func (g *GorillaHttpEngineAdapter) BodyParser(out interface{}) error {
	return json.NewDecoder(g.r.Body).Decode(out)
}

type GorillaServerAdapter struct {
	router *mux.Router
}

func NewGorillaServerAdapter() *GorillaServerAdapter {
	return &GorillaServerAdapter{
		router: mux.NewRouter(),
	}
}

func convertFiberRouteToGorilla(route string) string {
	re := regexp.MustCompile(`:([a-zA-Z0-9_]+)`)
	return re.ReplaceAllString(route, `{$1}`)
}

func (g *GorillaServerAdapter) Get(path string, handler HandlerHttpEngine) {
	g.router.HandleFunc(convertFiberRouteToGorilla(path), func(w http.ResponseWriter, r *http.Request) {
		err := handler(&GorillaHttpEngineAdapter{w: w, r: r})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}).Methods("GET")
}

func (g *GorillaServerAdapter) Post(path string, handler HandlerHttpEngine) {
	g.router.HandleFunc(convertFiberRouteToGorilla(path), func(w http.ResponseWriter, r *http.Request) {
		err := handler(&GorillaHttpEngineAdapter{w: w, r: r})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}).Methods("POST")
}

func (g *GorillaServerAdapter) Put(path string, handler HandlerHttpEngine) {
	g.router.HandleFunc(convertFiberRouteToGorilla(path), func(w http.ResponseWriter, r *http.Request) {
		err := handler(&GorillaHttpEngineAdapter{w: w, r: r})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}).Methods("PUT")
}

func (g *GorillaServerAdapter) Delete(path string, handler HandlerHttpEngine) {
	g.router.HandleFunc(convertFiberRouteToGorilla(path), func(w http.ResponseWriter, r *http.Request) {
		err := handler(&GorillaHttpEngineAdapter{w: w, r: r})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}).Methods("DELETE")
}

func (g *GorillaServerAdapter) Listen(addr string) error {
	return http.ListenAndServe(addr, g.router)
}

func (g *GorillaServerAdapter) Test(req *http.Request, msTimeout ...int) (*http.Response, error) {
	recorder := httptest.NewRecorder()
	g.router.ServeHTTP(recorder, req)
	return recorder.Result(), nil
}

func FactoryServer(adapter string) Server {
	if adapter == "fiber" {
		return NewFiberServerAdapter()
	}
	return NewGorillaServerAdapter()
}
