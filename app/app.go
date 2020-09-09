package app

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"

	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"

	"SweetDreams/config"
	"SweetDreams/controller"
	"SweetDreams/db"
)

type App struct {
	Router *mux.Router
	DB     *mongo.Database
}

// ConfigAndRunApp will create and initialize App structure. App factory function.
func ConfigAndRunApp(config *config.Config) {
	app := new(App)
	app.Initialize(config)
	app.Run(config.ServerHost)
}

// Initialize initialize the app with
func (app *App) Initialize(config *config.Config) {
	app.DB = db.InitialConnection(config.Db.Name, config.Db.Endpoint)
	// app.createIndexes()

	app.Router = mux.NewRouter()
	app.UseMiddleware(controller.JSONContentTypeMiddleware)
	app.setRouters()
}

// SetupRouters will register routes in router
func (app *App) setRouters() {
	app.post("/product", app.handleRequest(controller.CreateProduct))
	app.put("/product", app.handleRequest(controller.UpdateProduct))
	app.get("/product/{id}", app.handleRequest(controller.GetProduct))
	//app.delete("/product", app.handleRequest(product.Delete))
	app.get("/product", app.handleRequest(controller.GetAllProducts), "page", "{page}")

	app.post("/user/authenticate", app.handleRequest(controller.Authenticate))
	app.post("/user/register", app.handleRequest(controller.RegisterUser))
}

// Run will start the http server on host that you pass in. host:<ip:port>
func (app *App) Run(host string) {
	// use signals for shutdown server gracefully.
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, os.Interrupt, os.Kill)
	go func() {
		log.Fatal(http.ListenAndServe(host, app.Router))
	}()
	log.Printf("Server is listening on http://%s\n", host)
	sig := <-sigs
	log.Println("Signal: ", sig)

	log.Println("Stoping MongoDB Connection...")
	app.DB.Client().Disconnect(context.Background())
}

// UseMiddleware will add global middleware in router
func (app *App) UseMiddleware(middleware mux.MiddlewareFunc) {
	app.Router.Use(middleware)
}

// region RESTWrappers
func (app *App) get(path string, endpoint http.HandlerFunc, queries ...string) {
	app.Router.HandleFunc(path, endpoint).Methods("GET").Queries(queries...)
}

func (app *App) post(path string, endpoint http.HandlerFunc, queries ...string) {
	app.Router.HandleFunc(path, endpoint).Methods("POST").Queries(queries...)
}

func (app *App) put(path string, endpoint http.HandlerFunc, queries ...string) {
	app.Router.HandleFunc(path, endpoint).Methods("PUT").Queries(queries...)
}

func (app *App) patch(path string, endpoint http.HandlerFunc, queries ...string) {
	app.Router.HandleFunc(path, endpoint).Methods("PATCH").Queries(queries...)
}

func (app *App) delete(path string, endpoint http.HandlerFunc, queries ...string) {
	app.Router.HandleFunc(path, endpoint).Methods("DELETE").Queries(queries...)
}

// endregion RESTWrappers

// RequestHandlerFunction is a custome type that help us to pass db arg to all endpoints
type RequestHandlerFunction func(db *mongo.Database, w http.ResponseWriter, r *http.Request)

// handleRequest is a middleware we create for pass in db connection to endpoints.
func (app *App) handleRequest(handler RequestHandlerFunction) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(app.DB, w, r)
	}
}
