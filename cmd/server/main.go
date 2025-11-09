package main

import (
	"log"
	"net/http"

	handle "github.com/Dzox13524/PA_Funcpro_Kel12/internal/handler"
	"github.com/Dzox13524/PA_Funcpro_Kel12/internal/middleware"
	"github.com/Dzox13524/PA_Funcpro_Kel12/internal/platform/database"
	"github.com/Dzox13524/PA_Funcpro_Kel12/internal/repository"
	"github.com/Dzox13524/PA_Funcpro_Kel12/internal/response"
	"github.com/Dzox13524/PA_Funcpro_Kel12/internal/service"
)

func main() {
	// Konfigurasi database
	db := database.NewConnection()
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)	
	userHandler := handle.NewUserHandler(userService)
	// konfig comsole
	log.SetFlags(0)
	
	
	mux := http.NewServeMux()
	mux.HandleFunc("POST /users", userHandler.HandleCreateUser)
    mux.HandleFunc("GET /users", userHandler.HandleGetUserByID)
	mux.HandleFunc("GET /ping", func (w http.ResponseWriter, r *http.Request){
		type User struct {
        ID   string 
        Name string
    }
    user := User{ID: "123", Name: "John Doe"}
		response.WriteJSON(w, http.StatusOK, "succes", user)
	}) 

	var finalHandler http.Handler = mux
	finalHandler = middleware.Logging(finalHandler)


	http.ListenAndServe(":8080", finalHandler);
}	