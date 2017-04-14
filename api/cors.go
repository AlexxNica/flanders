package api

import "net/http"

//CORS the CORS middleware
func CORS(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "POST,PUT,GET,HEAD,DELETE,OPTIONS")
		w.Header().Set("Access-Control-Max-Age", "1728000")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
