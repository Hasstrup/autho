package middlewares

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/authenticate/services"
	utils "github.com/authenticate/utilities"
	"github.com/mongodb/mongo-go-driver/mongo"
)

type BusinessMiddleware struct {
	Client *mongo.Client
}

func (b BusinessMiddleware) EnforceApiKey(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("authorization")
		if apiKey == "" {
			utils.RespondWithJSON(w, 400, yieldError("Please send in the apikey in the request header with the 'x-authorization' field"))
			return
		}
		defer func() {
			if err := recover(); err != nil {
				log.Println(err)
				utils.RespondWithJSON(w, 400, yieldError("Something went wrong"))
			}
		}()
		// hash and write a cast
		app, err := services.YieldAppFromApiKey(apiKey, b.Client)
		if err != nil {
			utils.RespondWithJSON(w, 400, yieldError(err.Error()))
			return
		}
		appendToRequestBody(r, "main_application", app)
		next.ServeHTTP(w, r)
	})
}

func yieldError(message string) map[string]string {
	return map[string]string{"error": message}
}

func appendToRequestBody(r *http.Request, key string, dest interface{}) {
	var data map[string]interface{}
	_ = json.NewDecoder(r.Body).Decode(&data)
	data[key] = dest
	b, _ := json.Marshal(data)
	r.Body = ioutil.NopCloser(bytes.NewBuffer(b))
}
