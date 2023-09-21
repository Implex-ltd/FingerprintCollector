package main

import (
	"compress/flate"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Implex-ltd/collector/internal/fingerprint"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var (
	enckey         = "lmao15464notgonnagetthekeyifyesyouareagoodboy"
	allowedHeaders = "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization,X-CSRF-Token"
)

func SubmitFp(w http.ResponseWriter, r *http.Request) {
	var requestData FpPayload
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
		return
	}

	userId, err := fingerprint.Decrypt(requestData.ID, "broisatryharderlmao667")
	if err != nil {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))

		log.Println(err)
		return
	}

	_, err = fingerprint.Decrypt(requestData.N, "1337superpasslmaohowcanyoubegaylikethat")
	if err != nil {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
		log.Println(err)
		return
	}

	fmt.Println(userId)

	existingDoc := visitcollection.FindOne(context.TODO(), bson.D{{"visitor_id", userId}})

	if existingDoc.Err() == mongo.ErrNoDocuments {
		insertResult, err := collection.InsertOne(context.TODO(), requestData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, err = visitcollection.InsertOne(context.TODO(), bson.D{{"visitor_id", userId}})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Printf("Inserted a document with ID: %v\n", insertResult.InsertedID)
	} else if existingDoc.Err() != nil {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
		return
	} else {
		fmt.Println("Document with VisitorID already exists.")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func SendJsFile(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "../../assets/hsw.js")
}

func HandleRequests() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Timeout(30 * time.Second))

	compressor := middleware.NewCompressor(flate.BestSpeed)
	r.Use(compressor.Handler)

	r.Post("/submit", SubmitFp)
	r.Get("/hsw.js", SendJsFile)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("serve /")
		http.ServeFile(w, r, "../../assets/index.html")
	})

	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", allowedHeaders)
			w.Header().Set("Access-Control-Expose-Headers", "Authorization")
		}

		log.Println("../../assets/" + r.URL.Path)
		http.ServeFile(w, r, "../../assets/"+r.URL.Path)
	})

	http.ListenAndServe(":80", r)
}

func main() {
	log.Println("online")
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI("mongodb+srv://nikolahellatrigger:t8fU55cN7iBQ23Mu@cluster0.n0a2ebe.mongodb.net/?retryWrites=true&w=majority").SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Err(); err != nil {
		panic(err)
	}

	db := client.Database("fingerprint")
	collection = db.Collection("fp")
	visitcollection = db.Collection("visitors")

	HandleRequests()
}
