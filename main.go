package main

import (
	"crypto/tls"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"tile38-webserver/validator"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalln("Error loading .env file. Make sure it is present at the data folder")
	}
}

func validateRequest(validationType string, authHeader string) bool {
	if validationType == "jwt_rsa" {
		token := strings.Split(authHeader, " ")[1]
		result, err := validator.ValidateJWTRSA(token, os.Getenv("PUBKEY_URI"))
		if err != nil {
			log.Println(err)
			return false
		} else if !result {
			log.Println("Unauthorized request")
			return false
		}
	} else if validationType == "jwt_hmac" {
		token := strings.Split(authHeader, " ")[1]
		result, err := validator.ValidateJWTHMAC(token, os.Getenv("VALIDATION_SECRET"))
		if err != nil {
			log.Println(err)
			return false
		} else if !result {
			log.Println("Unauthorized request")
			return false
		}
	} else if validationType == "basic" {
		if authHeader != os.Getenv("VALIDATION_SECRET") {
			log.Println("Unauthorized request")
			return false
		}
	}
	return true
}

func main() {
	validationType := os.Getenv("VALIDATION_TYPE")
	router := mux.NewRouter()

	router.HandleFunc("/{query}", func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header["Authorization"]
		if len(authHeader) < 1 {
			log.Println("Missing Authorization header.")
			w.WriteHeader(401)
			return
		}
		if !validateRequest(validationType, authHeader[0]) {
			w.WriteHeader(401)
			return
		}
		// Ignore some browser requests
		if r.URL.String() == "/favicon.ico" {
			w.WriteHeader(404)
			return
		}

		resp, err := http.Get(os.Getenv("TILE38_URI") + r.URL.String())
		if err != nil {
			w.WriteHeader(500)
			return
		}
		defer resp.Body.Close()
		io.Copy(w, resp.Body)
	})

	writeTimeout, err := strconv.Atoi(os.Getenv("SERVER_WRITE_TIMEOUT"))
	if err != nil {
		return
	}
	readTimeout, err := strconv.Atoi(os.Getenv("SERVER_READ_TIMEOUT"))
	if err != nil {
		return
	}
	cfg := &tls.Config{
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
	}
	srv := &http.Server{
		Handler:      router,
		Addr:         os.Getenv("SERVER_ADDR"),
		WriteTimeout: time.Duration(writeTimeout) * time.Second,
		ReadTimeout:  time.Duration(readTimeout) * time.Second,
		TLSConfig:    cfg,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
	}

	log.Fatal(srv.ListenAndServeTLS(os.Getenv("SERVER_CERT"), os.Getenv("SERVER_KEY")))
}
