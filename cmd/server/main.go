package main

import (
	"log"
	"net/http"
	"github.com/Dzox13524/PA_Funcpro_Kel12/internal/domain"
	handle "github.com/Dzox13524/PA_Funcpro_Kel12/internal/handler"
	"github.com/Dzox13524/PA_Funcpro_Kel12/internal/middleware"
	"github.com/Dzox13524/PA_Funcpro_Kel12/internal/platform/database"
	"github.com/Dzox13524/PA_Funcpro_Kel12/internal/repository"
	"github.com/Dzox13524/PA_Funcpro_Kel12/internal/service"
)

func main() {
	db := database.NewConnection()

	db.AutoMigrate(&domain.PestReport{})

	userRepoGetByID := repository.NewGetUserByIDRepository(db)
	userRepoGetByEmail := repository.NewGetUserByEmailRepository(db)
	userRepoCreate := repository.NewCreateUserRepository(db)
	userRepoUpdate := repository.NewUpdateUserRepository(db)

	createUserService := service.NewCreateUser(userRepoCreate, userRepoGetByEmail)
	getUserByIDService := service.NewGetUserByID(userRepoGetByID)
	loginService := service.NewLoginService(userRepoGetByEmail)
	updateUserService := service.NewUpdateUser(userRepoUpdate)

	pestRepo := repository.NewPestRepository(db)
	pestService := service.NewPestService(pestRepo)

	log.SetFlags(0)
	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/v1/auth/register", handle.HandleCreateUser(createUserService))
	mux.HandleFunc("POST /api/v1/auth/login", handle.HandleLogin(loginService))
	mux.HandleFunc("GET /api/v1/users/{id}", handle.HandleGetUserByID(getUserByIDService))

	mux.HandleFunc("GET /api/v1/users/me", middleware.AuthMiddleware(handle.HandleGetMe(getUserByIDService)))
	mux.HandleFunc("PATCH /api/v1/users/me", middleware.AuthMiddleware(handle.HandleUpdateMe(updateUserService)))

	mux.HandleFunc("POST /api/v1/alerts", middleware.AuthMiddleware(handle.HandleCreateAlert(pestService)))
	mux.HandleFunc("GET /api/v1/alerts/map", handle.HandleGetMapData(pestService))
	mux.HandleFunc("GET /api/v1/alerts/{id}", handle.HandleGetAlertDetail(pestService))
	mux.HandleFunc("POST /api/v1/alerts/{id}/verify", middleware.AuthMiddleware(handle.HandleVerifyAlert(pestService)))

	var finalHandler http.Handler = mux
	finalHandler = middleware.Logging(finalHandler)
	finalHandler = middleware.CORSMiddleware(finalHandler)

	log.Println("Server running on port :8080")
	http.ListenAndServe(":8080", finalHandler)
}