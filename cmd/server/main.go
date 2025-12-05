package main

import (
	"log"
	"net/http"

	handle "github.com/Dzox13524/PA_Funcpro_Kel12/internal/handler"
	"github.com/Dzox13524/PA_Funcpro_Kel12/internal/middleware"
	"github.com/Dzox13524/PA_Funcpro_Kel12/internal/platform/database"
	"github.com/Dzox13524/PA_Funcpro_Kel12/internal/repository"
	"github.com/Dzox13524/PA_Funcpro_Kel12/internal/service"
)

func main() {
	db := database.NewConnection()

	userRepoGetByID := repository.NewGetUserByIDRepository(db)
    userRepoCreate := repository.NewCreateUserRepository(db)

	createUserService := service.NewCreateUser(userRepoCreate)
    getUserByIDService := service.NewGetUserByID(userRepoGetByID)
	
	log.SetFlags(0)
	
	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/v1/auth/register", middleware.MakeHandler(handle.HandleCreateUser(createUserService)))
	
	mux.HandleFunc("GET /api/v1/users/{id}", middleware.MakeHandler(handle.HandleGetUserByID(getUserByIDService)))

	var finalHandler http.Handler = mux
	finalHandler = middleware.Logging(finalHandler)

	log.Println("Server running on port :8080")
	http.ListenAndServe(":8080", finalHandler)
}