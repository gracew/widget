package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/gorilla/mux"
	"github.com/gracew/widget/graph"
	"github.com/gracew/widget/graph/generated"
	"github.com/gracew/widget/graph/model"
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
	r.HandleFunc("/apis/{apiName}/{env}", createHandler).Methods("POST", "OPTIONS")
	r.HandleFunc("/apis/{apiName}/{env}/{id}", readHandler).Methods("GET", "OPTIONS")
	r.HandleFunc("/apis/{apiName}/{env}", listHandler).Methods("GET", "OPTIONS")
	// TODO(gracew): remove cors later
	r.Use(mux.CORSMethodMiddleware(r))
	http.Handle("/", r)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func createSchema(db *pg.DB) error {
	for _, model := range []interface{}{(*model.API)(nil), (*model.Deploy)(nil)} {
		err := db.CreateTable(model, &orm.CreateTableOptions{
			IfNotExists: true,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

type parseRes struct {
	CreatedAt string `json:"createdAt"`
	ObjectID  string `json:"objectId"`
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	if r.Method == http.MethodOptions {
		return
	}
	vars := mux.Vars(r)
	parseClass := vars["apiName"] + vars["env"]
	// talk to parse, forward request body
	// TODO(gracew): don't hardcode this
	req, err := http.NewRequest("POST", "http://localhost:1337/parse/classes/"+parseClass, r.Body)
	if err != nil {
		panic(err)
	}
	req.Header.Add("X-Parse-Application-Id", "appId")
	req.Header.Add("Content-type", "application/json")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	/*var parseRes parseRes
	err = json.NewDecoder(res.Body).Decode(&parseRes)
	print(parseRes.ObjectID)
	*/
	// return response
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	w.Write(body)
}

func readHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	if r.Method == http.MethodOptions {
		return
	}
	vars := mux.Vars(r)
	parseClass := vars["apiName"] + vars["env"]
	objectID := vars["id"]
	req, err := http.NewRequest("GET", "http://localhost:1337/parse/classes/"+parseClass+"/"+objectID, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Add("X-Parse-Application-Id", "appId")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	w.Write(body)
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	if r.Method == http.MethodOptions {
		return
	}
	pageSize, ok := r.URL.Query()["pageSize"]
	data := ""
	if ok && len(pageSize[0]) >= 1 {
		data = "limit=" + pageSize[0]
	}
	vars := mux.Vars(r)
	parseClass := vars["apiName"] + vars["env"]
	req, err := http.NewRequest("GET", "http://localhost:1337/parse/classes/"+parseClass, strings.NewReader(data))
	if err != nil {
		panic(err)
	}
	req.Header.Add("X-Parse-Application-Id", "appId")
	req.Header.Add("Content-type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	w.Write(body)
}
