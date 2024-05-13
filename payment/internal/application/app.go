package application

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/PickHD/LezPay/payment/internal/config"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

// App..
type App struct {
	Application   *echo.Echo
	Context       context.Context
	Config        *config.Configuration
	Logger        *logrus.Logger
	DB            *pgxpool.Pool
	Redis         *redis.Client
	KafkaProducer *kafka.Writer
	Tracer        *trace.TracerProvider
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
		app.Logger.Error("failed init Tracer", err)
		return app, nil
	}

	otel.SetTracerProvider(app.Tracer)

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

	err = initCreateTopics(app.Config)
	if err != nil {
		app.Logger.Error("failed create Kafka topics", err)
		return app, nil
	}

	app.KafkaProducer = initProducer(app.Config)

	app.Application = echo.New()
	app.Application.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

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

// initCreateTopics will create all defined kafka topic
func initCreateTopics(cfg *config.Configuration) error {
	conn, err := kafka.Dial("tcp", cfg.Kafka.FirstBrokerHost)
	if err != nil {
		panic(err.Error())
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		panic(err.Error())
	}
	var controllerConn *kafka.Conn
	controllerConn, err = kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		panic(err.Error())
	}
	defer controllerConn.Close()

	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             cfg.Kafka.TopicRequestPayment,
			NumPartitions:     1,
			ReplicationFactor: 1,
		},
	}

	err = controllerConn.CreateTopics(topicConfigs...)
	if err != nil {
		panic(err.Error())
	}

	return nil
}

// initProducer will create generic kafka producer configuration
func initProducer(cfg *config.Configuration) *kafka.Writer {
	return &kafka.Writer{
		Addr:         kafka.TCP(cfg.Kafka.FirstBrokerHost, cfg.Kafka.SecondBrokerHost),
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequiredAcks(cfg.Kafka.RequiredAcks),
	}
}
