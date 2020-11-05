package validator

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
)

var (
	valid_alg []string
	keyCache  = PubKeyCache{"", time.Time{}}
)

func init() {
	// Not checking for errors since main.go already checked
	_ = godotenv.Load(".env")
	valid_alg = strings.Split(os.Getenv("JWT_VALID_ALG"), " ")
}

type PubKeyCache struct {
	pemFile string
	time    time.Time
}

func getKey(url string) string {
	cacheDuration, _ := strconv.Atoi(os.Getenv("CACHE_DURATION"))
	if time.Since(keyCache.time) >= time.Duration(cacheDuration)*time.Minute {
		pem, err := downloadFile(url)
		if err != nil {
			log.Printf("Cannot download public key from %s!", url)
			log.Fatal(err)
		}
		keyCache.pemFile = pem
		keyCache.time = time.Now()
	}
	return keyCache.pemFile
}

func validAlg(x string) bool {
	for _, n := range valid_alg {
		if x == n {
			return true
		}
	}
	return false
}

func downloadFile(url string) (f string, err error) {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bad status: %s", resp.Status)
	}

	// Writer the body to file
	buf := new(strings.Builder)
	_, err = io.Copy(buf, resp.Body)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func ValidateJWTRSA(jwtToken string, pubKeyURL string) (bool, error) {
	pubKeyFile := getKey(pubKeyURL)
	pubKey, keyErr := jwt.ParseRSAPublicKeyFromPEM([]byte(pubKeyFile))
	if keyErr != nil {
		return false, keyErr
	}
	token, jwtError := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		alg := token.Header["alg"]
		if !validAlg(alg.(string)) {
			return nil, fmt.Errorf("Invalid alg: %s", alg)
		}
		return pubKey, nil
	})
	if !token.Valid {
		if ve, ok := jwtError.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				fmt.Println("Malformed token")
			} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				// Token is either expired or not active yet
				fmt.Println("Invalid or expired token")
				return false, jwtError
			} else {
				fmt.Println("Couldn't handle this token:")
				return false, jwtError
			}
		}
	}
	return true, nil
}

func ValidateJWTHMAC(jwtToken string, secret string) (bool, error) {
	token, jwtError := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		alg := token.Header["alg"]
		if !validAlg(alg.(string)) {
			return nil, fmt.Errorf("Invalid alg: %s", alg)
		}
		return []byte(secret), nil
	})
	if jwtError != nil {
		log.Fatal("Error", jwtError)
		return false, jwtError
	} else if !token.Valid {
		log.Fatal("Token is invalid")
		return false, nil
	} else if ve, ok := jwtError.(*jwt.ValidationError); ok {
		log.Println(ve)
		log.Println(ok)
	}
	return true, nil
}
