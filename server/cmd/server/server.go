package server

import (
	"VyacheslavKuchumov/test-backend/service/auth"
	"VyacheslavKuchumov/test-backend/service/draft"
	"VyacheslavKuchumov/test-backend/service/equipment"
	"VyacheslavKuchumov/test-backend/service/equipmentindraft"
	"VyacheslavKuchumov/test-backend/service/equipmentinproject"
	"VyacheslavKuchumov/test-backend/service/equipmentset"
	"VyacheslavKuchumov/test-backend/service/project"
	"VyacheslavKuchumov/test-backend/service/projecttype"
	"VyacheslavKuchumov/test-backend/service/settype"
	"VyacheslavKuchumov/test-backend/service/tracker"
	"VyacheslavKuchumov/test-backend/service/user"
	"VyacheslavKuchumov/test-backend/service/warehouse"
	"database/sql"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

type Server struct {
	addr string
	db   *sql.DB
}

func NewServer(addr string, db *sql.DB) *Server {
	return &Server{
		addr: addr,
		db:   db,
	}
}

func (s *Server) Run() error {
	log.Println("Listening on", s.addr)
	return http.ListenAndServe(s.addr, s.router())
}

func (s *Server) router() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	userStore := user.NewStore(s.db)
	userHandler := user.NewHandler(userStore)

	trackerStore := tracker.NewStore(s.db)
	setTypeService := settype.NewService(trackerStore)
	projectTypeService := projecttype.NewService(trackerStore)
	warehouseService := warehouse.NewService(trackerStore)
	equipmentSetService := equipmentset.NewService(trackerStore)
	equipmentService := equipment.NewService(trackerStore)
	projectService := project.NewService(trackerStore)
	draftService := draft.NewService(trackerStore)
	equipmentInProjectService := equipmentinproject.NewService(trackerStore)
	equipmentInDraftService := equipmentindraft.NewService(trackerStore)
	authMiddleware := auth.JWTAuthMiddleware(userStore)
	apiAuthMiddleware := auth.JWTAuthMiddlewareWithExclusions(
		userStore,
		"/api/v1/login",
		"/api/v1/register",
	)

	r.With(authMiddleware).Handle("/swagger/*", httpSwagger.Handler())

	r.Route("/api/v1", func(api chi.Router) {
		api.Use(apiAuthMiddleware)
		user.RegisterRoutes(api, userHandler)
		settype.RegisterRoutes(api, setTypeService)
		projecttype.RegisterRoutes(api, projectTypeService)
		warehouse.RegisterRoutes(api, warehouseService)
		equipmentset.RegisterRoutes(api, equipmentSetService)
		equipment.RegisterRoutes(api, equipmentService)
		project.RegisterRoutes(api, projectService)
		draft.RegisterRoutes(api, draftService)
		equipmentinproject.RegisterRoutes(api, equipmentInProjectService)
		equipmentindraft.RegisterRoutes(api, equipmentInDraftService)
	})

	return r
}
