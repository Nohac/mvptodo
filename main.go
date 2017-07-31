//go:generate granate

package main

import (
	"fmt"
	"net/http"

	"github.com/gocraft/dbr"
	"github.com/graphql-go/handler"
	"github.com/noh4ck/mvptodo/models"
	"github.com/noh4ck/mvptodo/schema"
	"github.com/noh4ck/mvptodo/storage"
	"golang.org/x/net/context"
)

type SessionHandler struct {
	SchemaHandler *handler.Handler
	Database      *dbr.Connection
}

func (h SessionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Got request")
	db := h.Database.NewSession(nil)

	session := models.SessionType{
		Request:  r,
		Response: w,
		DB:       db,
	}

	auth := models.NewAuth(&session)

	ctx := context.Background()
	ctx = context.WithValue(ctx, "session", &session)
	ctx = context.WithValue(ctx, "auth", auth)

	auth.Authenticate()

	w.Header().Set("Content-Type", "application/json")
	h.SchemaHandler.ContextHandler(ctx, w, r)

}

func main() {
	fmt.Println("Starting mvptodo graphql server")
	database := storage.New("devuser", "devpassword", "mvptodo")

	// session := database.NewSession(nil)
	// seeds.Run(session)

	root := models.Root{}

	schema.Init(schema.ProviderConfig{
		Query:    root,
		Mutation: root,
		Relay:    root,
	})

	schemaHandler := handler.New(&handler.Config{
		Schema: schema.Schema(),
		Pretty: true,
	})

	http.Handle("/graphql", SessionHandler{
		SchemaHandler: schemaHandler,
		Database:      database,
	})
	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(graphiql)
	}))
	http.ListenAndServe(":8080", nil)
}

var graphiql = []byte(`
<!DOCTYPE html>
<html>
   <head>
      <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/graphiql/0.7.8/graphiql.css" />
      <script src="https://cdnjs.cloudflare.com/ajax/libs/fetch/1.0.0/fetch.min.js"></script>
      <script src="https://cdnjs.cloudflare.com/ajax/libs/react/15.3.2/react.min.js"></script>
      <script src="https://cdnjs.cloudflare.com/ajax/libs/react/15.3.2/react-dom.min.js"></script>
      <script src="https://cdnjs.cloudflare.com/ajax/libs/graphiql/0.7.8/graphiql.js"></script>
   </head>
   <body style="width: 100%; height: 100%; margin: 0; overflow: hidden;">
      <div id="graphiql" style="height: 100vh;">Loading...</div>
      <script>
         function graphQLFetcher(graphQLParams) {
            graphQLParams.variables = graphQLParams.variables ? JSON.parse(graphQLParams.variables) : null;
            return fetch("/graphql", {
               method: "post",
               body: JSON.stringify(graphQLParams),
               credentials: "include",
            }).then(function (response) {
               return response.text();
            }).then(function (responseBody) {
               try {
                  return JSON.parse(responseBody);
               } catch (error) {
                  return responseBody;
               }
            });
         }
         ReactDOM.render(
            React.createElement(GraphiQL, {fetcher: graphQLFetcher}),
            document.getElementById("graphiql")
         );
      </script>
   </body>
</html>
`)
