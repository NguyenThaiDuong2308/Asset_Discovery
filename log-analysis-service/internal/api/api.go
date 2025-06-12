package api

import (
	"database/sql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type Server struct {
	router *mux.Router
	db     *sql.DB
}

func New(db *sql.DB) *Server {
	r := mux.NewRouter()
	r.Use(corsMiddleware)
	s := &Server{
		router: r,
		db:     db,
	}
	s.registerRoutes()

	return s
}

func (s *Server) Start(port string) error {
	log.Printf("Starting API server on port %s", port)
	return http.ListenAndServe(":"+port, s.router)
}

func (s *Server) registerRoutes() {

	s.router.HandleFunc("/api/assets", s.getAssetsHandler).Methods("GET", "OPTIONS")
	s.router.HandleFunc("/api/assets/{ip}", s.getAssetByIPHandler).Methods("GET", "OPTIONS")

	s.router.HandleFunc("/api/logs", s.getLogsHandler).Methods("GET", "OPTIONS")

	s.router.HandleFunc("/api/stats", s.getStatsHandler).Methods("GET", "OPTIONS")

	s.router.HandleFunc("/api/upload/{type}", s.uploadLogHandler).Methods("POST", "OPTIONS")
	s.router.HandleFunc("/api/services", s.getServicesHandler).Methods("GET", "OPTIONS")
	s.router.HandleFunc("/api/services/{ip}", s.getServicesByAssetHandler).Methods("GET", "OPTIONS")

}
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("CORS middleware triggered for:", r.Method, r.URL.Path)
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Nếu là OPTIONS request (preflight), trả về ngay
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
