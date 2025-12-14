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

	db.AutoMigrate(
		&domain.User{},
		&domain.Product{},
		&domain.Article{},
		&domain.Question{},
		&domain.Answer{},
		&domain.QuestionLike{},
		&domain.Favorite{},
		&domain.PestReport{},
		&domain.MarketTransaction{},
		&domain.PestVerification{},
	)

	userRepoGetByID := repository.NewGetUserByIDRepository(db)
	userRepoGetByEmail := repository.NewGetUserByEmailRepository(db)
	userRepoCreate := repository.NewCreateUserRepository(db)
	userRepoUpdate := repository.NewUpdateUserRepository(db)

	prodRepoCreate := repository.NewCreateProductRepository(db)
	prodRepoGetAll := repository.NewGetAllProductsRepository(db)
	prodRepoGetByID := repository.NewGetProductByIDRepository(db)
	prodRepoUpdate := repository.NewUpdateProductRepository(db)
	prodRepoDelete := repository.NewDeleteProductRepository(db)
	prodRepoSearch := repository.NewSearchProductsRepository(db)
	prodRepoMetaCrops := repository.NewGetMetaCropsRepository(db)
	prodRepoMetaRegions := repository.NewGetMetaRegionsRepository(db)

	marketRepoCreate := repository.NewCreateTransactionRepository(db)
	marketRepoGetByID := repository.NewGetTransactionByIDRepository(db)
	marketRepoGetByUser := repository.NewGetTransactionsByUserRepository(db)
	marketRepoUpdateStatus := repository.NewUpdateTransactionStatusRepository(db)

	articlerepoCreate := repository.NewCreateArticleRepository(db)
	articlerepoGetByID := repository.NewGetArticleByIDRepository(db)
	articlerepoGetAll := repository.NewGetAllArticlesRepository(db)

	qRepoCreate := repository.NewCreateQuestionRepository(db)
	qRepoGet := repository.NewGetAllQuestionsRepository(db)
	qRepoDetail := repository.NewGetQuestionByIDRepository(db)
	aRepoCreate := repository.NewCreateAnswerRepository(db)
	likeRepo := repository.NewToggleQuestionLikeRepository(db)
	favRepo := repository.NewToggleFavoriteRepository(db)

	pestRepo := repository.NewPestRepository(db)

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
	svcSearchProducts := service.NewSearchProductsService(prodRepoSearch)
	svcMetaCrops := service.NewGetMetaCropsService(prodRepoMetaCrops)
	svcMetaRegions := service.NewGetMetaRegionsService(prodRepoMetaRegions)

	svcCreateReservation := service.NewCreateReservationService(marketRepoCreate, prodRepoGetByID, prodRepoUpdate)
	svcCreateOrder := service.NewCreateOrderService(marketRepoCreate, prodRepoGetByID, prodRepoUpdate)
	svcGetMyReservations := service.NewGetUserTransactionsService(marketRepoGetByUser, domain.TypeReservation)
	svcGetMyOrders := service.NewGetUserTransactionsService(marketRepoGetByUser, domain.TypeOrder)
	svcGetTransDetail := service.NewGetTransactionDetailService(marketRepoGetByID)
	svcUpdateTransStatus := service.NewUpdateTransactionStatusService(marketRepoUpdateStatus)

	svccreatearticle := service.NewCreateArticleService(articlerepoCreate)
	svcgetarticlebyid := service.NewGetArticleByIDService(articlerepoGetByID)
	svcgetallarticles := service.NewGetAllArticlesService(articlerepoGetAll)

	svcCreateQ := service.NewCreateQuestion(qRepoCreate)
	svcGetFeed := service.NewGetFeed(qRepoGet)
	svcGetDetail := service.NewGetQuestionDetail(qRepoDetail)
	svcAddAns := service.NewAddAnswer(aRepoCreate)
	svcLike := service.NewToggleLike(likeRepo)
	svcFav := service.NewToggleFav(favRepo)

	pestService := service.NewPestService(pestRepo)

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

	mux.HandleFunc("GET /api/v1/market/search", handle.HandleSearchProducts(svcSearchProducts))
	mux.HandleFunc("GET /api/v1/market/meta/crops", handle.HandleGetMetaCrops(svcMetaCrops))
	mux.HandleFunc("GET /api/v1/market/meta/regions", handle.HandleGetMetaRegions(svcMetaRegions))

	mux.HandleFunc("POST /api/v1/market/reservations", middleware.AuthMiddleware(handle.HandleCreateReservation(svcCreateReservation)))
	mux.HandleFunc("POST /api/v1/market/orders", middleware.AuthMiddleware(handle.HandleCreateOrder(svcCreateOrder)))
	mux.HandleFunc("GET /api/v1/users/me/reservations", middleware.AuthMiddleware(handle.HandleGetMyReservations(svcGetMyReservations)))
	mux.HandleFunc("GET /api/v1/users/me/orders", middleware.AuthMiddleware(handle.HandleGetMyOrders(svcGetMyOrders)))
	mux.HandleFunc("GET /api/v1/market/orders/{id}", middleware.AuthMiddleware(handle.HandleGetTransactionDetail(svcGetTransDetail)))
	mux.HandleFunc("POST /api/v1/market/reservations/{id}/confirm", middleware.AuthMiddleware(handle.HandleConfirmReservation(svcUpdateTransStatus)))
	mux.HandleFunc("POST /api/v1/market/reservations/{id}/cancel", middleware.AuthMiddleware(handle.HandleCancelReservation(svcUpdateTransStatus)))
	mux.HandleFunc("PATCH /api/v1/market/orders/{id}/status", middleware.AuthMiddleware(handle.HandleUpdateOrderStatus(svcUpdateTransStatus)))

	mux.HandleFunc("POST /api/v1/articles", middleware.AuthMiddleware(handle.HandleCreateArticle(svccreatearticle)))
	mux.HandleFunc("GET /api/v1/articles", handle.HandleGetAllArticles(svcgetallarticles))
	mux.HandleFunc("GET /api/v1/articles/{id}", handle.HandleGetArticleByID(svcgetarticlebyid))

	mux.HandleFunc("GET /api/v1/questions", middleware.AuthMiddlewareOptional(handle.HandleGetFeed(svcGetFeed)))
	mux.HandleFunc("POST /api/v1/questions", middleware.AuthMiddleware(handle.HandleCreateQuestion(svcCreateQ)))
	mux.HandleFunc("GET /api/v1/questions/{id}", middleware.AuthMiddlewareOptional(handle.HandleGetQuestionDetail(svcGetDetail)))
	mux.HandleFunc("POST /api/v1/questions/{id}/answers", middleware.AuthMiddleware(handle.HandleAddAnswer(svcAddAns)))
	mux.HandleFunc("POST /api/v1/questions/{id}/like", middleware.AuthMiddleware(handle.HandleToggleLike(svcLike)))
	mux.HandleFunc("POST /api/v1/questions/{id}/favorite", middleware.AuthMiddleware(handle.HandleToggleFav(svcFav)))

	mux.HandleFunc("POST /api/v1/alerts", middleware.AuthMiddleware(handle.HandleCreateAlert(pestService)))
	mux.HandleFunc("GET /api/v1/alerts/map", handle.HandleGetMapData(pestService))
	mux.HandleFunc("GET /api/v1/alerts/{id}", handle.HandleGetAlertDetail(pestService))
	mux.HandleFunc("POST /api/v1/alerts/{id}/verify", middleware.AuthMiddleware(handle.HandleVerifyAlert(pestService)))

	var finalHandler http.Handler = mux
	finalHandler = middleware.Logging(finalHandler)
	finalHandler = middleware.CORSMiddleware(finalHandler)

	log.Println("Server running at http://localhost:8080")
	http.ListenAndServe(":8080", finalHandler)
}