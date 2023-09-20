package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Implex-ltd/collector/internal/fingerprint"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var (
	enckey = "test"
)

func SubmitFp(w http.ResponseWriter, r *http.Request) {
	var requestData FpPayload

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err := fingerprint.Decrypt(requestData.Data.N, enckey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = fingerprint.Decrypt(requestData.Data.F, enckey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = fingerprint.Decrypt(requestData.Data.D, enckey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	decJ, err := fingerprint.Decrypt(requestData.Data.J, enckey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var j_dec interface{}

	if err = json.Unmarshal([]byte(decJ), &j_dec); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var fp Fpjs
	if err = json.Unmarshal([]byte(decJ), &fp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println(fp.VisitorID)

	filter := bson.D{{"visitor_id", fp.VisitorID}}
	existingDoc := visitcollection.FindOne(context.TODO(), filter)

	if existingDoc.Err() == mongo.ErrNoDocuments {
		insertResult, err := collection.InsertOne(context.TODO(), requestData)
		if err != nil {
			log.Fatal(err)
		}

		_, err = visitcollection.InsertOne(context.TODO(), bson.D{{"visitor_id", fp.VisitorID}})
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Inserted a document with ID: %v\n", insertResult.InsertedID)
	} else if existingDoc.Err() != nil {
		log.Fatal(existingDoc.Err())
	} else {
		fmt.Println("Document with VisitorID already exists.")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func VerificationStart(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, ".../../assets/challenge/index.html")
}

func SendJsFile(w http.ResponseWriter, r *http.Request) {
	workDir, _ := os.Getwd()
	filesDir := filepath.Join(workDir, ".../../assets/challenge")

	http.ServeFile(w, r, filepath.Join(filesDir, "challenge.js"))
}

func HandleRequests() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Timeout(30 * time.Second))

	r.Get("/verify", VerificationStart)
	r.Post("/submitfp", SubmitFp)
	r.Get("/challenge.js", SendJsFile)

	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		log.Println(".../../assets/challenge/" + r.URL.Path)
		http.ServeFile(w, r, ".../../assets/challenge/"+r.URL.Path)
	})

	http.ListenAndServe(":8080", r)
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
