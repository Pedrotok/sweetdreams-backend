package app

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"

	"github.com/rs/cors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"golang.org/x/net/context"

	middleware "SweetDreams/app/middleware"
	"SweetDreams/config"
	"SweetDreams/controller"
	"SweetDreams/db"
)

type App struct {
	Router *mux.Router
	DB     *mongo.Database
}

// ConfigAndRunApp will create and initialize App structure. App factory function.
func ConfigAndRunApp(config config.GeneralConfig) {
	app := new(App)
	app.initialize(config)
	app.Run(config.ServerHost)
}

func (app *App) initialize(config config.GeneralConfig) {
	app.DB = db.InitialConnection(config.Mongo.Name, config.Mongo.Endpoint)
	app.createIndexes()

	app.Router = mux.NewRouter()
	app.UseMiddleware(middleware.JSONContentTypeMiddleware)
	app.setRouters()
}

func (app *App) createIndexes() {
	keys := bsonx.Doc{
		{Key: "email", Value: bsonx.Int32(1)},
	}
	users := app.DB.Collection("User")
	db.SetIndexes(users, keys)
}

// SetupRouters will register routes in router
func (app *App) setRouters() {

	app.post("/product", middleware.AuthMiddleware(app.DB, app.handleRequest(controller.CreateProduct)))
	app.put("/product", middleware.AuthMiddleware(app.DB, app.handleRequest(controller.UpdateProduct)))
	app.get("/product/{id}", app.handleRequest(controller.GetProduct))
	//app.delete("/product", app.handleRequest(controller.DeleteProduct))
	app.get("/product", app.handleRequest(controller.GetAllProducts))
	app.get("/product", app.handleRequest(controller.GetAllProducts), "page", "{page}")

	app.post("/user/authenticate", app.handleRequest(controller.Authenticate))
	app.post("/user/register", app.handleRequest(controller.RegisterUser))

	// TODO remove this endpoint
	app.get("/user", app.handleRequest(controller.GetAllUsers))
	app.get("/user", app.handleRequest(controller.GetAllUsers), "page", "{page}")

	app.post("/token/refresh", app.handleRequest(controller.RefreshToken))
}

// Run will start the http server on host that you pass in. host:<ip:port>
func (app *App) Run(host string) {
	// use signals for shutdown server gracefully.
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, os.Interrupt, os.Kill)
	go func() {
		c := cors.New(cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowCredentials: true,
			AllowedHeaders:   []string{"*"},
		})

		log.Fatal(http.ListenAndServe(host, c.Handler(app.Router)))
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
type RequestHandlerFunction func(db *mongo.Database, w http.ResponseWriter, r *http.Request) error

// handleRequest is a middleware we create for pass in db connection to endpoints.
func (app *App) handleRequest(handler RequestHandlerFunction) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := handler(app.DB, w, r); err != nil {
			switch e := err.(type) {
			case controller.StatusError:
				log.Printf("HTTP %d - %s", e.Status(), e)
				http.Error(w, e.Error(), e.Status())
			default:
				http.Error(w, http.StatusText(http.StatusInternalServerError),
					http.StatusInternalServerError)
			}
		}
	}
}
