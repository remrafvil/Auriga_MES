package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/remrafvil/Auriga_API/config"
	"github.com/remrafvil/Auriga_API/databases"
	"github.com/remrafvil/Auriga_API/internal/httpapi"
	"github.com/remrafvil/Auriga_API/internal/httpapi/handlers"
	"github.com/remrafvil/Auriga_API/internal/httpapi/middlewares"
	"github.com/remrafvil/Auriga_API/internal/repositories"
	"github.com/remrafvil/Auriga_API/internal/services"
	"github.com/remrafvil/Auriga_API/internal/utils"
	"gorm.io/gorm"

	//Utilizamos la librería para la inyección de dependencias
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Params struct {
	fx.In

	Lc            fx.Lifecycle
	Config        *config.Settings
	DB            *gorm.DB
	InfluxManager *databases.InfluxClientManager // Cambiado de InfluxDB a InfluxManager
	Echo          *echo.Echo
	Handlers      []handlers.Handler `group:"handlers"`
	Logger        *zap.Logger
}

func main() {
	// Crear un contexto raíz para la aplicación
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	app := fx.New(
		/* 	Con Provide e Invoke, FX podrá crear todas las dependencias y ejecutarlas
		Provide se pasan todos los métodos que devuelven un struct
		Dentro del fx.Provide no es necesario utilziar las paréntesis para las funciones			*/
		fx.Provide(
			func() context.Context { return ctx }, // Proveer el contexto raíz
			//context.Background,
			config.New,
			databases.New,
			databases.NewInfluxClientManager,
			databases.NewWorkeraClient,
			echo.New,
			zap.NewProduction, // Cambiar a NewProduction en entorno de producción
			utils.NewCustomValidator,
		),
		repositories.Module,
		services.Module,
		httpapi.Module,
		fx.Invoke(
			setLifeCycle,
		),
	)

	// app.Run() // ejecutamos nuestra aplicación para que inicilice
	// Manejar señales de sistema para shutdown graceful
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		cancel() // Cancelar el contexto cuando llegue señal
	}()

	if err := app.Start(ctx); err != nil {
		panic(err)
	}

	// Esperar a que la aplicación termine
	<-ctx.Done()

	// Shutdown graceful con timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := app.Stop(shutdownCtx); err != nil {
		panic(err)
	}
}

func setLifeCycle(p Params) {
	// Crear un contexto para el monitoreo que viva durante toda la aplicación
	monitorCtx := context.Background()
	p.Lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			// El manager ya se inicializó automáticamente en el constructor
			// Solo necesitamos verificar el estado
			availableConnections := p.InfluxManager.GetAvailableClients()

			p.Logger.Info("Estado inicial de conexiones InfluxDB",
				zap.Int("total_configuradas", len(p.InfluxManager.GetAllConfigs())),
				zap.Int("conexiones_activas", len(availableConnections)),
				zap.Strings("conexiones_disponibles", availableConnections),
			)

			// Iniciar monitoreo de conexiones en background
			go p.InfluxManager.StartConnectionMonitor(monitorCtx)

			// Configurar el validador desde utils
			validator := utils.NewCustomValidator()
			p.Echo.Validator = validator

			// Configurar middlewares según el entorno
			if p.Config.IsProduction() {
				p.Logger.Info("Starting in PRODUCTION mode",
					zap.String("port", p.Config.App.EchoPort),
					zap.String("environment", p.Config.App.Zap),
				)
				// Configuraciones específicas de producción
				p.Echo.HideBanner = true
				p.Echo.Debug = false
			} else {
				p.Logger.Info("Starting in DEVELOPMENT mode",
					zap.String("port", p.Config.App.EchoPort),
					zap.String("environment", p.Config.App.Zap),
					zap.Bool("debug", true),
				)
				// Configuraciones de desarrollo
				p.Echo.Debug = true
				p.Echo.HideBanner = false
			}

			// ✅ Configurar middlewares globales (sin auth)
			middlewares.MainMiddlewares(p.Echo, p.Config, p.Logger)
			for _, h := range p.Handlers {
				h.RegisterRoutes(p.Echo, p.Config)
			}
			p.Logger.Info("Server configuration completed",
				zap.String("app_name", p.Config.App.Name),
				zap.String("version", p.Config.App.Version),
				zap.String("environment", p.Config.App.Zap),
				zap.Int("influx_connections_available", len(availableConnections)),
			)

			// Iniciar servidor en goroutine separada
			go func() {
				p.Logger.Info("Starting server",
					zap.String("port", p.Config.App.EchoPort),
					zap.Int("influx_connections", len(availableConnections)),
				)
				if err := p.Echo.Start(p.Config.App.EchoPort); err != nil {
					p.Logger.Fatal("Failed to start server", zap.Error(err))
				}
			}()

			return nil
		},

		OnStop: func(ctx context.Context) error {
			p.Logger.Info("Shutting down server")

			// Cerrar todas las conexiones InfluxDB
			p.Logger.Info("Closing InfluxDB connections")
			p.InfluxManager.CloseAll()

			// Shutdown del servidor Echo con timeout
			shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
			defer cancel()

			if err := p.Echo.Shutdown(shutdownCtx); err != nil {
				p.Logger.Error("Error shutting down server", zap.Error(err))
				return err
			}

			p.Logger.Info("Server shutdown completed")
			return nil
		},
	})
}
