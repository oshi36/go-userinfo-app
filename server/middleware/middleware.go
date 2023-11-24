package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"htmlgo/server/models"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/joho/godotenv"
	"os"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
    "go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo"

)


// Database Name
const dbName = "userdata"

// Collection name
const collName = "userapp"

// collection object/instance
var collection *mongo.Collection

// create connection with mongo db
func init() {
 
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	  }
	
	connectionString := os.Getenv("MONGO_URI")

	// Set client options
	clientOptions := options.Client().ApplyURI(connectionString)

          // Add the instrumentation to the client
	clientOptions.Monitor = otelmongo.NewMonitor()

	// connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}
 
  
	fmt.Println("Connected to MongoDB!")

	collection = client.Database(dbName).Collection(collName)

	fmt.Println("Collection instance created!")
	initTracer()


}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		var tmpl = template.Must(template.New("form").ParseFiles("templates/index.html"))
		var err = tmpl.Execute(w, nil)

		//payload := getUsers()
		//json.NewEncoder(w).Encode(payload)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	http.Error(w, "", http.StatusBadRequest)
}

// func getUsers() []primitive.M {
// 	cur, err := collection.Find(context.Background(), bson.D{{}})
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	var results []primitive.M
// 	for cur.Next(context.Background()) {
// 		var result bson.M
// 		e := cur.Decode(&result)
// 		if e != nil {
// 			log.Fatal(e)
// 		}
// 		// fmt.Println("cur..>", cur, "result", reflect.TypeOf(result), reflect.TypeOf(result["_id"]))
// 		results = append(results, result)

// 	}

// 	if err := cur.Err(); err != nil {
// 		log.Fatal(err)
// 	}

// 	cur.Close(context.Background())
// 	return results
// }

func NewUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		template.Must(template.New("form").ParseFiles("templates/index.html"))
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data := models.Detail{
			Name:  r.FormValue("name"),
			Email: r.FormValue("email"),
		}
		_ = json.NewDecoder(r.Body).Decode(&data)
		insertOneUser(data)
		json.NewEncoder(w).Encode(data)

		// if err := tmpl.Execute(w, data); err != nil {
		// 	http.Error(w, err.Error(), http.StatusInternalServerError)
		// }
		// return

	}
	http.Error(w, "", http.StatusBadRequest)
}

func insertOneUser(detail models.Detail) {
	insertResult, err := collection.InsertOne(context.Background(), detail)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted a Single Record ", insertResult.InsertedID)

}

func User(w http.ResponseWriter, r *http.Request) {
	if r.Method == "DELETE" || r.Method == "GET" {
		var tmpl = template.Must(template.New("form").ParseFiles("templates/index.html"))
		var err = tmpl.Execute(w, nil)

		if err = r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		params := mux.Vars(r)
		deleteUser(params["id"])
		json.NewEncoder(w).Encode(params["id"])

		http.Error(w, "", http.StatusBadRequest)

	}
}

func deleteUser(detail string) {
	fmt.Println(detail)
	id, _ := primitive.ObjectIDFromHex(detail)
	filter := bson.M{"_id": id}
	d, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Deleted Document", d.DeletedCount)

}

func initTracer() {
	ctx := context.Background()

	client := otlptracehttp.NewClient()
	exporter, err := otlptrace.New(ctx, client)
	if err != nil {
		log.Fatalf("failed to initialize exporter: %e", err)
	}

	res, err := resource.New(ctx)
	if err != nil {
		log.Fatalf("failed to initialize resource: %e", err)
	}

	// Create the trace provider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)

	// Set the global trace provider
	otel.SetTracerProvider(tp)

	// Set the propagator
	propagator := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})
	otel.SetTextMapPropagator(propagator)
}
