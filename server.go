package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/gorilla/mux"
	"github.com/gracew/widget/graph"
	"github.com/gracew/widget/graph/generated"
	"github.com/gracew/widget/graph/model"
	"github.com/gracew/widget/parse"
	"github.com/rs/cors"
)

const defaultPort = "8080"

func main() {
	db := pg.Connect(&pg.Options{User: "postgres"})
	defer db.Close()

	err := createSchema(db)
	if err != nil {
		panic(err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))

	http.Handle("/playground", playground.Handler("GraphQL playground", "/query"))
	// TODO(gracew): remove cors later
	http.Handle("/query", cors.Default().Handler(srv))

	// individual API routes
	r := mux.NewRouter()
	r.HandleFunc("/apis/{apiID}/{env}", createHandler).Methods("POST", "OPTIONS")
	r.HandleFunc("/apis/{apiID}/{env}/{id}", readHandler).Methods("GET", "OPTIONS")
	r.HandleFunc("/apis/{apiID}/{env}", listHandler).Methods("GET", "OPTIONS")
	// TODO(gracew): remove cors later
	r.Use(mux.CORSMethodMiddleware(r))
	http.Handle("/", r)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func createSchema(db *pg.DB) error {
	for _, model := range []interface{}{(*model.API)(nil), (*model.Deploy)(nil), (*model.Auth)(nil), (*model.TestToken)(nil)} {
		err := db.CreateTable(model, &orm.CreateTableOptions{
			IfNotExists: true,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	if r.Method == http.MethodOptions {
		return
	}
	vars := mux.Vars(r)

	// get the userId
	parseToken := r.Header["X-Parse-Session-Token"][0]
	userID, err := parse.GetUserId(parseToken)
	if err != nil {
		panic(err)
	}

	// add createdBy to the original req
	var originalReq map[string]interface{}
	err = json.NewDecoder(r.Body).Decode(&originalReq)
	originalReq["createdBy"] = userID

	// return response
	res, err := parse.CreateObject(vars["apiID"], vars["env"], originalReq)
	if err != nil {
		panic(err)
	}
	w.Write(res)
}

func readHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	if r.Method == http.MethodOptions {
		return
	}
	vars := mux.Vars(r)

	res, err := parse.GetObject(vars["apiID"], vars["env"], vars["id"])
	if err != nil {
		panic(err)
	}
	w.Write(res)

	// get the userId
	// TODO(gracew): do this in parallel as the object fetch
	/*parseToken := r.Header["X-Parse-Session-Token"][0]
	userID, err := parse.GetUserId(parseToken)
	if err != nil {
		panic(err)
	}

	// get the authorization policy. if CREATED_BY, then check userId == res.createdBy
	auth, err := store.Auth(vars["apiID"])
	if err != nil {
		panic(err)
	}

	if auth.ReadPolicy.Type == model.AuthPolicyTypeCreatedBy {

	}
	var createdBy createdBy*/
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	if r.Method == http.MethodOptions {
		return
	}
	pageSizes, ok := r.URL.Query()["pageSize"]
	pageSize := "100"
	if ok && len(pageSizes[0]) >= 1 {
		pageSize = pageSizes[0]
	}
	vars := mux.Vars(r)
	res, err := parse.ListObjects(vars["apiID"], vars["env"], pageSize)
	if err != nil {
		panic(err)
	}
	w.Write(res)
}
