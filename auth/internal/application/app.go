package application

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"github.com/PickHD/LezPay/auth/internal/config"
	"github.com/PickHD/LezPay/auth/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	trace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gopkg.in/gomail.v2"
)

// App ..
type App struct {
	Application  *gin.Engine
	Context      context.Context
	Config       *config.Configuration
	Logger       *logrus.Logger
	DB           *pgxpool.Pool
	Redis        *redis.Client
	Tracer       *trace.TracerProvider
	Mailer       *gomail.Dialer
	CustomerGRPC *grpc.ClientConn
	MerchantGRPC *grpc.ClientConn
	WalletGRPC   *grpc.ClientConn
}

// SetupApplication configuring dependencies app needed
func SetupApplication(ctx context.Context) (*App, error) {
	var err error

	app := &App{}
	app.Context = context.TODO()
	app.Config = config.NewConfig()

	// custom log app with logrus
	logWithLogrus := logrus.New()
	logWithLogrus.Formatter = &logrus.JSONFormatter{}
	logWithLogrus.ReportCaller = true
	app.Logger = logWithLogrus

	// initialize tracers
	app.Tracer, err = initJaegerTracerProvider(app.Config)
	if err != nil {
		app.Logger.Error("failed init Jaeger Tracer", err)
		return app, nil
	}

	otel.SetTracerProvider(app.Tracer)

	// initialize mailer
	app.Mailer = initSMTPMailDialer(app.Config)

	// "postgres://username:password@localhost:5432/database_name"
	dbpool, err := pgxpool.Connect(context.Background(), fmt.Sprintf("postgres://%s:%s@%s:%d/%s", app.Config.Database.Username, app.Config.Database.Password, app.Config.Database.Host, app.Config.Database.Port, app.Config.Database.Name))
	if err != nil {
		app.Logger.Error("failed create pool connection Postgres", err)
		return app, nil
	}

	app.DB = dbpool

	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", app.Config.Redis.Host, app.Config.Redis.Port),
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	app.Redis = redisClient

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	app.CustomerGRPC, err = grpc.Dial(app.Config.Service.GRPCCustomerHost, opts...)
	if err != nil {
		app.Logger.Error("failed Dial CustomerGRPC, error :", err)
		return app, err
	}

	app.MerchantGRPC, err = grpc.Dial(app.Config.Service.GRPCMerchantHost, opts...)
	if err != nil {
		app.Logger.Error("failed Dial MerchantGRPC, error :", err)
		return app, err
	}

	app.WalletGRPC, err = grpc.Dial(app.Config.Service.GRPCWalletHost, opts...)
	if err != nil {
		app.Logger.Error("failed Dial WalletGRPC, error :", err)
		return app, err
	}

	app.Application = gin.New()
	app.Application.Use(middleware.CORSMiddleware())

	app.Logger.Info("APP RUN SUCCESSFULLY")

	return app, nil
}

// Close method will close any instances before app terminated
func (a *App) Close(ctx context.Context) {
	a.Logger.Info("APP CLOSED SUCCESSFULLY")

	defer func(ctx context.Context) {
		// DB
		a.DB.Close()

		// Redis
		a.Redis.Close()

		// TRACER
		ctx, cancel := context.WithTimeout(ctx, time.Second*5)
		defer cancel()
		if err := a.Tracer.Shutdown(ctx); err != nil {
			panic(err)
		}

	}(ctx)
}

// initJaegerTracerProvider returns an OpenTelemetry TracerProvider configured to use
// the OTLP Jaeger exporter that will send spans to the provided url. The returned
// TracerProvider will also use a Resource configured with all the information
// about the application.
func initJaegerTracerProvider(cfg *config.Configuration) (*trace.TracerProvider, error) {
	// Create the HTTP Client exporter with endpoint refer to jaeger URL
	client := otlptracehttp.NewClient(otlptracehttp.WithEndpoint(cfg.Tracer.JaegerURL), otlptracehttp.WithInsecure())
	exp, err := otlptrace.New(context.Background(), client)
	if err != nil {
		return nil, err
	}
	tp := trace.NewTracerProvider(
		// Always be sure to batch in production.
		trace.WithBatcher(exp),
		// Record information about this application in a Resource.
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(cfg.Server.AppName),
			attribute.String("environment", cfg.Server.AppEnv),
			attribute.String("ID", cfg.Server.AppID),
		)),
	)
	return tp, nil
}

// initSMTPMailDialer returns an Dialer configured to use
func initSMTPMailDialer(cfg *config.Configuration) *gomail.Dialer {
	d := gomail.NewDialer(cfg.Mailer.Host, cfg.Mailer.Port, cfg.Mailer.Username, cfg.Mailer.Password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	return d
}
