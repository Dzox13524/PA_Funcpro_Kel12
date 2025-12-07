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
	db.AutoMigrate(&domain.User{}, &domain.Product{})

	userRepoGetByID := repository.NewGetUserByIDRepository(db)
	userRepoGetByEmail := repository.NewGetUserByEmailRepository(db)
	userRepoCreate := repository.NewCreateUserRepository(db)
	userRepoUpdate := repository.NewUpdateUserRepository(db)

	prodRepoCreate := repository.NewCreateProductRepository(db)
	prodRepoGetAll := repository.NewGetAllProductsRepository(db)
	prodRepoGetByID := repository.NewGetProductByIDRepository(db)
	prodRepoUpdate := repository.NewUpdateProductRepository(db)
	prodRepoDelete := repository.NewDeleteProductRepository(db)

	createUserService := service.NewCreateUser(userRepoCreate, userRepoGetByEmail)
	getUserByIDService := service.NewGetUserByID(userRepoGetByID)
	loginService := service.NewLoginService(userRepoGetByEmail)
	updateUserService := service.NewUpdateUser(userRepoUpdate)

	svcCreateProduct := service.NewCreateProductService(prodRepoCreate)
	svcGetAllProducts := service.NewGetAllProductsService(prodRepoGetAll)
	svcGetProductByID := service.NewGetProductByIDService(prodRepoGetByID)
	svcUpdateProduct := service.NewUpdateProductService(prodRepoUpdate, prodRepoGetByID)
	svcDeleteProduct := service.NewDeleteProductService(prodRepoDelete, prodRepoGetByID)
	svcUploadImage := service.NewUploadProductImageService(prodRepoGetByID, prodRepoUpdate)

	log.SetFlags(0)
	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/v1/auth/register", handle.HandleCreateUser(createUserService))
	mux.HandleFunc("POST /api/v1/auth/login", handle.HandleLogin(loginService))
	mux.HandleFunc("GET /api/v1/users/{id}", handle.HandleGetUserByID(getUserByIDService))

	mux.HandleFunc("GET /api/v1/users/me", middleware.AuthMiddleware(handle.HandleGetMe(getUserByIDService)))
	mux.HandleFunc("PATCH /api/v1/users/me", middleware.AuthMiddleware(handle.HandleUpdateMe(updateUserService)))

	mux.HandleFunc("POST /api/v1/market/products", middleware.AuthMiddleware(handle.HandleCreateProduct(svcCreateProduct)))
	mux.HandleFunc("GET /api/v1/market/products", handle.HandleGetAllProducts(svcGetAllProducts))
	mux.HandleFunc("GET /api/v1/market/products/{id}", handle.HandleGetProductByID(svcGetProductByID))
	mux.HandleFunc("PUT /api/v1/market/products/{id}", middleware.AuthMiddleware(handle.HandleUpdateProduct(svcUpdateProduct)))
	mux.HandleFunc("DELETE /api/v1/market/products/{id}", middleware.AuthMiddleware(handle.HandleDeleteProduct(svcDeleteProduct)))
	mux.HandleFunc("POST /api/v1/market/products/{id}/upload", middleware.AuthMiddleware(handle.HandleUploadProductImage(svcUploadImage)))


	var finalHandler http.Handler = mux
	finalHandler = middleware.Logging(finalHandler)
	finalHandler = middleware.CORSMiddleware(finalHandler) //untuk akses option (web)

	log.Println("Server running at http://localhost:8080")
	http.ListenAndServe(":8080", finalHandler)
}