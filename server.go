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
	"github.com/gracew/widget/store"
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
	for _, model := range []interface{}{
		(*model.API)(nil),
		(*model.Deploy)(nil),
		(*model.Auth)(nil),
		(*model.TestToken)(nil),
		(*model.CustomLogic)(nil),
	} {
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

	// delegate to parse
	vars := mux.Vars(r)
	res, err := parse.CreateObject(vars["apiID"], vars["env"], originalReq)
	if err != nil {
		panic(err)
	}

	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		panic(err)
	}
}

func readHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	if r.Method == http.MethodOptions {
		return
	}

	// get the userId
	parseToken := r.Header["X-Parse-Session-Token"][0]
	userID, err := parse.GetUserId(parseToken)
	if err != nil {
		panic(err)
	}

	// delegate to parse
	vars := mux.Vars(r)
	res, err := parse.GetObject(vars["apiID"], vars["env"], vars["id"])
	if err != nil {
		panic(err)
	}

	// fetch the authorization policy
	// TODO(gracew): parallelize some of these requests
	auth, err := store.Auth(vars["apiID"])
	if err != nil {
		panic(err)
	}

	if auth.ReadPolicy.Type == model.AuthPolicyTypeCreatedBy {
		if userID != res.CreatedBy {
			json.NewEncoder(w).Encode(&unauthorized{Message: "unauthorized"})
			return
		}
	}
	// TODO(gracew): support other authz policies

	json.NewEncoder(w).Encode(&res)
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	if r.Method == http.MethodOptions {
		return
	}

	// get the userId
	parseToken := r.Header["X-Parse-Session-Token"][0]
	userID, err := parse.GetUserId(parseToken)
	if err != nil {
		panic(err)
	}

	// delegate to parse
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

	// fetch the authorization policy
	// TODO(gracew): parallelize some of these requests
	auth, err := store.Auth(vars["apiID"])
	if err != nil {
		panic(err)
	}

	var filtered []parse.ObjectRes
	if auth.ReadPolicy.Type == model.AuthPolicyTypeCreatedBy {
		for i := 0; i < len(res.Results); i++ {
			if userID == res.Results[i].CreatedBy {
				filtered = append(filtered, res.Results[i])
			}
		}
	}
	// TODO(gracew): support other authz policies

	json.NewEncoder(w).Encode(&parse.ListRes{Results: filtered})
}

type createdBy struct {
	CreatedBy string `json:"createdBy"`
}

type unauthorized struct {
	Message string `json:"message"`
}