package main

import (
	"compress/flate"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Implex-ltd/collector/internal/fingerprint"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/surrealdb/surrealdb.go"
)

var (
	KEY_USERID = "broisatryharderlmao667"
	KEY_FP     = "1337superpasslmaohowcanyoubegaylikethat"
)

func submit(w http.ResponseWriter, r *http.Request) {
	defer func() {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}()

	var data FpPayload
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return
	}

	if data.ID == "" || data.N == "" {
		return
	}

	UserID, err := fingerprint.Decrypt(data.ID, KEY_USERID)
	if err != nil {
		return
	}

	Fp, err := fingerprint.Decrypt(data.N, KEY_FP)
	if err != nil {
		return
	}

	// do...

	createdData, err := DB.Create("fp", Fingerprint{
		ID:          UserID,
		Fingerprint: base64.StdEncoding.EncodeToString([]byte(Fp)),
	})
	if err != nil {
		log.Printf("VisitorID (already seen): %s", UserID)
		return
	}

	dbFp := make([]Fingerprint, 1)
	err = surrealdb.Unmarshal(data, &createdData)
	if err != nil {
		log.Println(err)
		return
	}

	log.Printf("VisitorID (new: %s): %s", dbFp[0].ID, UserID)
}

func main() {
	ConnectDB("144.172.76.66", "root", "rootnikoontoplmao5245", 8000)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Timeout(15 * time.Second))
	r.Use(middleware.NewCompressor(flate.BestSpeed).Handler)

	r.Post("/submit", submit)

	r.Get("/hsw.js", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../../assets/hsw.js")
	})

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../../assets/index.html")
	})

	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization,X-CSRF-Token")
			w.Header().Set("Access-Control-Expose-Headers", "Authorization")
		}

		http.ServeFile(w, r, fmt.Sprintf("../../assets/%s", r.URL.Path))
	})

	http.ListenAndServe(":80", r)
}
