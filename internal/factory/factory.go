package factory

import (
	"fmt"
	"log"

	"settlements/internal/config"
	"settlements/internal/db"
	"settlements/internal/repo"
	"settlements/internal/service"
	"settlements/internal/transport/http/controller"
	"settlements/internal/transport/http/router"

	"gorm.io/gorm"
)

// ApplicationFactory is a factory for creating and initializing application components
// This implements the Factory Pattern to centralize object creation and dependency management
type ApplicationFactory struct {
	config *config.Config
	db     *gorm.DB
	repo   *repo.CityRepo
}

// NewApplicationFactory creates a new ApplicationFactory with loaded configuration
// and established database connection
func NewApplicationFactory(cfg *config.Config) (*ApplicationFactory, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	dbConnection, err := db.Connect(&cfg.Database)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return &ApplicationFactory{
		config: cfg,
		db:     dbConnection,
	}, nil
}

// CreateRouter creates and returns a new Router instance
// Satisfies the router creation contract
func (f *ApplicationFactory) CreateRouter() *router.Router {
	return router.New()
}

// CreateRepository creates and returns a new CityRepository instance
// Lazy initialization pattern: repository is created once and cached
func (f *ApplicationFactory) CreateRepository() (*repo.CityRepo, error) {
	if f.repo != nil {
		return f.repo, nil
	}

	if f.db == nil {
		return nil, fmt.Errorf("database connection not initialized")
	}

	f.repo = repo.New(f.db)
	return f.repo, nil
}

// CreateService creates and returns a new Service instance
// Dependencies (repository) are automatically resolved via factory
func (f *ApplicationFactory) CreateService() (*service.Service, error) {
	repo, err := f.CreateRepository()
	if err != nil {
		return nil, fmt.Errorf("failed to create repository: %w", err)
	}

	return service.New(repo), nil
}

// CreateController creates and returns a new MainController instance
// Dependencies (service) are automatically resolved via factory
func (f *ApplicationFactory) CreateController() (*controller.MainController, error) {
	svc, err := f.CreateService()
	if err != nil {
		return nil, fmt.Errorf("failed to create service: %w", err)
	}

	return controller.New(svc), nil
}

// GetDatabase returns the underlying database connection
// Useful for migrations and advanced operations
func (f *ApplicationFactory) GetDatabase() *gorm.DB {
	return f.db
}

// GetConfig returns the loaded configuration
// Useful for accessing server settings
func (f *ApplicationFactory) GetConfig() *config.Config {
	return f.config
}

// ApplicationBootstrapper encapsulates complete application initialization
// Coordinates factory methods in proper sequence
type ApplicationBootstrapper struct {
	factory *ApplicationFactory
}

// NewApplicationBootstrapper creates a bootstrapper with a factory
func NewApplicationBootstrapper(factory *ApplicationFactory) *ApplicationBootstrapper {
	return &ApplicationBootstrapper{
		factory: factory,
	}
}

// InitializeApplication sets up all application components in correct order
// Returns a fully initialized application context
func (b *ApplicationBootstrapper) InitializeApplication() (*ApplicationContext, error) {
	// Create router
	router := b.factory.CreateRouter()
	log.Println("Router initialized")

	// Create controller with all dependencies
	controller, err := b.factory.CreateController()
	if err != nil {
		return nil, fmt.Errorf("failed to create controller: %w", err)
	}
	log.Println("Controller initialized with service and repository")

	// Register routes
	router.GET("/", controller.GetMainPage)
	log.Println("Routes registered")

	return &ApplicationContext{
		Router:     router,
		Controller: controller,
		Config:     b.factory.GetConfig(),
	}, nil
}

// ApplicationContext holds the fully initialized application components
type ApplicationContext struct {
	Router     *router.Router
	Controller *controller.MainController
	Config     *config.Config
}

// GetServerAddress returns the formatted server address (host:port)
func (ctx *ApplicationContext) GetServerAddress() string {
	return ":" + ctx.Config.Server.Port
}
