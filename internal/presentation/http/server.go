package http

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/accio/internal/application/dto"
	"github.com/accio/internal/infrastructure/container"
)

// Server represents the HTTP server
type Server struct {
	router    *chi.Mux
	container *container.Container
	port      int
}

// NewServer creates a new HTTP server
func NewServer(container *container.Container, port int) *Server {
	router := chi.NewRouter()

	// Middleware
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(60 * time.Second))

	// CORS
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	return &Server{
		router:    router,
		container: container,
		port:      port,
	}
}

// Start starts the HTTP server
func (s *Server) Start() error {
	// Register routes
	s.registerRoutes()

	// Create server
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", s.port),
		Handler: s.router,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Starting server on port %d", s.port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown server
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Error shutting down server: %v", err)
	}

	log.Println("Server stopped")
	return nil
}

// registerRoutes registers all routes
func (s *Server) registerRoutes() {
	// API routes
	s.router.Route("/api", func(r chi.Router) {
		// Health check
		r.Get("/health", s.handleHealth())

		// Profiles
		r.Route("/profiles", func(r chi.Router) {
			r.Get("/", s.handleGetProfiles())
			r.Get("/{platform}/{username}", s.handleGetProfile())
			r.Get("/search", s.handleSearchProfiles())
		})

		// Search history
		r.Route("/search-history", func(r chi.Router) {
			r.Get("/popular", s.handleGetPopularSearches())
		})
	})

	// Web UI routes
	s.router.Get("/", s.handleIndex())
	s.router.Get("/search", s.handleSearchPage())
	s.router.Get("/profile/{platform}/{username}", s.handleProfilePage())

	// Static files
	fileServer := http.FileServer(http.Dir("./ascendio/static"))
	s.router.Handle("/static/*", http.StripPrefix("/static", fileServer))
}

// handleHealth handles the health check endpoint
func (s *Server) handleHealth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	}
}

// handleIndex handles the index page
func (s *Server) handleIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Render index template
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Accio - Username Search Tool</title>
    <link rel="stylesheet" href="/static/css/styles.css">
    <script src="https://unpkg.com/htmx.org@1.9.6"></script>
</head>
<body>
    <header>
        <h1>Accio</h1>
        <p>Username Search Tool</p>
    </header>
    <main>
        <section class="search-section">
            <h2>Search for Usernames</h2>
            <form hx-get="/search" hx-target="#results" hx-indicator="#spinner">
                <div class="form-group">
                    <label for="search-type">Search by:</label>
                    <select id="search-type" name="type">
                        <option value="username">Username</option>
                        <option value="name">Real Name</option>
                    </select>
                </div>
                <div class="form-group">
                    <label for="query">Search:</label>
                    <input type="text" id="query" name="query" placeholder="Enter username or real name" required>
                </div>
                <div class="form-group">
                    <label for="platforms">Platforms:</label>
                    <select id="platforms" name="platforms" multiple>
                        <option value="Twitter">Twitter</option>
                        <option value="GitHub">GitHub</option>
                        <option value="Twitch">Twitch</option>
                        <option value="Instagram">Instagram</option>
                    </select>
                    <small>Hold Ctrl/Cmd to select multiple platforms</small>
                </div>
                <div class="form-group">
                    <button type="submit">Search</button>
                </div>
            </form>
            <div id="spinner" class="htmx-indicator">
                <div class="spinner"></div>
            </div>
        </section>
        <section id="results" class="results-section">
            <!-- Results will be loaded here -->
        </section>
    </main>
    <footer>
        <p>&copy; 2023 Accio - Username Search Tool</p>
    </footer>
</body>
</html>
`))
	}
}

// handleSearchPage handles the search page
func (s *Server) handleSearchPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get query parameters
		query := r.URL.Query().Get("query")
		searchType := r.URL.Query().Get("type")
		platforms := r.URL.Query()["platforms"]

		if query == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Query parameter is required"))
			return
		}

		ctx := r.Context()
		var profiles []*dto.ProfileDTO

		// Search for profiles
		if searchType == "name" {
			// Search by real name
			result, err := s.container.ProfileService.SearchProfilesByName(ctx, query)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(fmt.Sprintf("Error searching profiles: %v", err)))
				return
			}
			profiles = result
		} else {
			// Search by username
			for _, platform := range platforms {
				profile, err := s.container.ProfileService.GetProfileByUsername(ctx, query, platform)
				if err != nil {
					continue
				}
				if profile != nil {
					profiles = append(profiles, profile)
				}
			}
		}

		// Record search in history
		if s.container.SearchHistoryService != nil {
			go func() {
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				s.container.SearchHistoryService.RecordSearch(ctx, query, len(profiles))
			}()
		}

		// Render search results
		w.Header().Set("Content-Type", "text/html")
		if len(profiles) == 0 {
			w.Write([]byte(`
<div class="no-results">
    <h3>No profiles found</h3>
    <p>No profiles were found matching your search criteria.</p>
</div>
`))
			return
		}

		// Render profiles
		w.Write([]byte(`
<div class="results">
    <h3>Search Results</h3>
    <p>Found ` + fmt.Sprintf("%d", len(profiles)) + ` profiles matching your search criteria.</p>
    <div class="profiles-grid">
`))

		for _, profile := range profiles {
			w.Write([]byte(`
<div class="profile-card" hx-get="/profile/` + profile.Platform + `/` + profile.Username + `" hx-target="#results">
    <div class="profile-image">
        <img src="` + profile.ImageURL + `" alt="` + profile.RealName + `">
    </div>
    <div class="profile-info">
        <h4>` + profile.RealName + `</h4>
        <p>@` + profile.Username + ` on ` + profile.Platform + `</p>
        <p>` + fmt.Sprintf("%d", profile.FollowerCount) + ` followers</p>
    </div>
</div>
`))
		}

		w.Write([]byte(`
    </div>
</div>
`))
	}
}

// handleProfilePage handles the profile page
func (s *Server) handleProfilePage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get path parameters
		platform := chi.URLParam(r, "platform")
		username := chi.URLParam(r, "username")

		if platform == "" || username == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Platform and username parameters are required"))
			return
		}

		// Get profile
		ctx := r.Context()
		profile, err := s.container.ProfileService.GetProfileByUsername(ctx, username, platform)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Error getting profile: %v", err)))
			return
		}

		if profile == nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Profile not found"))
			return
		}

		// Render profile
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`
<div class="profile-detail">
    <div class="profile-header">
        <div class="profile-image">
            <img src="` + profile.ImageURL + `" alt="` + profile.RealName + `">
        </div>
        <div class="profile-info">
            <h3>` + profile.RealName + `</h3>
            <p>@` + profile.Username + ` on ` + profile.Platform + `</p>
            <p>` + fmt.Sprintf("%d", profile.FollowerCount) + ` followers</p>
            <p><a href="` + profile.ProfileURL + `" target="_blank">View Profile</a></p>
        </div>
    </div>
    <div class="profile-bio">
        <h4>Bio</h4>
        <p>` + profile.Bio + `</p>
    </div>
    <div class="profile-data">
        <h4>Profile Data</h4>
        <table>
            <tr>
                <th>Key</th>
                <th>Value</th>
            </tr>
`))

		for key, value := range profile.PlatformData {
			w.Write([]byte(`
            <tr>
                <td>` + key + `</td>
                <td>` + value + `</td>
            </tr>
`))
		}

		w.Write([]byte(`
        </table>
    </div>
    <div class="profile-actions">
        <button hx-get="/search" hx-target="#results" class="back-button">Back to Search</button>
    </div>
</div>
`))
	}
}

// API Handlers

// handleGetProfiles handles the get profiles endpoint
func (s *Server) handleGetProfiles() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: Implement
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotImplemented)
		w.Write([]byte(`{"error":"Not implemented"}`))
	}
}

// handleGetProfile handles the get profile endpoint
func (s *Server) handleGetProfile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: Implement
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotImplemented)
		w.Write([]byte(`{"error":"Not implemented"}`))
	}
}

// handleSearchProfiles handles the search profiles endpoint
func (s *Server) handleSearchProfiles() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: Implement
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotImplemented)
		w.Write([]byte(`{"error":"Not implemented"}`))
	}
}

// handleGetPopularSearches handles the get popular searches endpoint
func (s *Server) handleGetPopularSearches() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: Implement
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotImplemented)
		w.Write([]byte(`{"error":"Not implemented"}`))
	}
}
