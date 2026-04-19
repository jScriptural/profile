package middleware

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"profile/internal/models"
	"regexp"
	"strings"
)

func CORS(next http.Handler, v string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", v)


		b,_ := json.Marshal(r.Header);
		s := string(b)
		log.Printf("REQUEST HEADERS:\n  method: %v  %v",r.Method,s);
		next.ServeHTTP(w, r)
	})
}

func BearerAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonBytes, _ := json.Marshal(struct {
			IP        string `json:"ip,omitempty"`
			UserAgent string `json:"user_agent,omitempty"`
			Referer   string `json:"referer,omitempty"`
		}{
			r.RemoteAddr,
			r.UserAgent(),
			r.Referer(),
		},
		)

		jsn := string(jsonBytes)
		log.Printf("Client: %v", jsn)

		authToken := r.Header.Get("Authorization")
		envToken := strings.TrimSpace(os.Getenv("BEARER_TOKEN"))
		log.Printf("envToken: %v", envToken)
		bRegexp := regexp.MustCompile(`\s*?Bearer\s+?([A-Za-z0-9]*?)\s*?$`)
		tks := bRegexp.FindStringSubmatch(authToken)

		log.Printf("tks: %v", tks)
		if tks == nil || envToken == "" || tks[1] != envToken {
			log.Printf("Failed access attempt: %v", jsn)
			b, _ := json.Marshal(models.APIResponse{Status: "error", Message: "Unauthorizsed"})
			msg := string(b)
			http.Error(w, msg, http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
