package main

import (
	"context"
	"log"
	"net/http"

	"github.com/Okwujohngibson25/movie-go/internal/config"
	"github.com/Okwujohngibson25/movie-go/internal/database"
	"github.com/Okwujohngibson25/movie-go/internal/middleware"
	"github.com/Okwujohngibson25/movie-go/internal/movie"
	"github.com/Okwujohngibson25/movie-go/internal/postgres"
	"github.com/Okwujohngibson25/movie-go/internal/screen"
	"github.com/Okwujohngibson25/movie-go/internal/showing"
	"github.com/Okwujohngibson25/movie-go/internal/user"
	"github.com/Okwujohngibson25/movie-go/internal/validation"
	"github.com/go-chi/chi/v5"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	validate := validation.NewValidator()

	ctx := context.Background()
	pool, err := database.ConnectPostgres(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}

	defer pool.Close()

	userRepo := postgres.NewUserSqlcRepository(pool)
	userService := user.NewUserService(userRepo, cfg)
	userHandler := user.NewUserHandler(userService, validate)
	userService.SeedAdmin(cfg)

	movieRepo := postgres.NewMovieSqlcRepository(pool)
	movieService := movie.NewMovieService(movieRepo)
	movieHandler := movie.NewMovieHandler(movieService, validate)

	screenRepo := postgres.NewScreenSqlcRepository(pool)
	screenService := screen.NewScreenService(screenRepo)
	screenHandler := screen.NewScreenHandler(screenService, validate)

	showingRepo := postgres.NewShowingSqlcRepository(pool)
	showingService := showing.NewShowingService(showingRepo, movieRepo)
	showingHandler := showing.NewShowingHandler(showingService)

	r := chi.NewRouter()
	r.Route("/api/v1", func(r chi.Router) {
		// public routes
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", userHandler.CreateUser)
			r.Post("/login", userHandler.Login)
		})

		r.Route("/admin", func(r chi.Router) {
			// authenticated routes
			r.Group(func(r chi.Router) {
				r.Use(middleware.RequireAuth(cfg.App.JWTSecret))
				r.Use(middleware.RequireAdmin)
				r.Get("/movies", movieHandler.GetAllMovies)
				r.Get("/movies/{movieID}", movieHandler.GetMovieByID)
				r.Put("/movies/{movieID}", movieHandler.UpdateMovieByID)
				r.Post("/movies", movieHandler.CreateMovie)
				r.Delete("/movies/{movieID}", movieHandler.DeleteMovieByID)

				r.Get("/screens", screenHandler.GetAllScreens)
				r.Get("/screens/{screensID}", screenHandler.GetScreensByID)
				r.Put("/screens/{screensID}", screenHandler.UpdateScreensByID)
				r.Post("/screens", screenHandler.CreateScreen)

				r.Get("/screens/{screensID}/seats", screenHandler.GetSeatsByScreenID)

				r.Get("/showings", showingHandler.GetAllShowings)
				r.Post("/showings", showingHandler.CreateShowing)
			})
		})
	})

	if err := http.ListenAndServe(cfg.App.Port, r); err != nil {
		log.Fatalf("error starting server: %v", err)
	}
}
