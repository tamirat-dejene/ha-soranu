package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// Directory structure aligned with services/auth-service
var structure = []string{
	"cmd",
	"docs",
	"internal/api/grpc/dto",
	"internal/api/grpc/handler",
	"internal/api/grpc/interceptor",
	"internal/domain",
	"internal/repository",
	"internal/usecase",
	"internal/util",
	"migrations",
	"test",
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run newservice.go <service-name>")
		os.Exit(1)
	}

	service := strings.TrimSpace(os.Args[1])
	if service == "" {
		fmt.Println("Invalid service name")
		os.Exit(1)
	}

	root := filepath.Join("services", service)

	// Prevent overwriting existing service
	if _, err := os.Stat(root); err == nil {
		fmt.Printf("Service '%s' already exists. Aborting.\n", service)
		os.Exit(1)
	}

	fmt.Printf("Creating new microservice: %s\n\n", service)

	// Create folder structure
	for _, dir := range structure {
		path := filepath.Join(root, dir)
		if err := os.MkdirAll(path, 0o755); err != nil {
			fmt.Printf("Failed to create %s: %v\n", path, err)
			os.Exit(1)
		}
		fmt.Printf("Created: %s\n", path)
	}

	// Create boilerplate files similar to auth-service
	modPath := modulePath()
	createFile(filepath.Join(root, "cmd", "main.go"), mainTemplate(service))
	createFile(filepath.Join(root, "env.go"), envTemplate(service))
	createFile(filepath.Join(root, "migrations", "migrate.go"), migrateTemplate(service, modPath))
	createFile(filepath.Join(root, "internal", "api", "grpc", "interceptor", "logging_interceptor.go"), interceptorTemplate())
	createFile(filepath.Join(root, "README.md"), readmeTemplate(service))

	fmt.Println("\nDone! Service skeleton created.")
	fmt.Printf("Next steps:\n  cd services/%s\n", service)
}

func createFile(path, content string) {
	if data, err := os.ReadFile(path); err == nil && len(data) > 0 {
		fmt.Printf("File exists, skipping: %s\n", path)
		return
	}

	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		fmt.Printf("Failed to write %s: %v\n", path, err)
		os.Exit(1)
	}
	fmt.Printf("Created file: %s\n", path)
}

func mainTemplate(service string) string {
	// derive env prefix from service (first token before '-')
	parts := strings.Split(service, "-")
	prefix := strings.ToUpper(parts[0])
	return fmt.Sprintf(`package main

import (
	"fmt"
	"net"

	"%s/shared/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	svc "%s/services/%s"
)

func main() {
	// 1. Load Configuration
	env, err := svc.GetEnv()
	if err != nil {
		panic(err)
	}

	// 2. Initialize Logger
	logger.InitLogger(env.SRV_ENV)
	defer logger.Log.Sync()
	logger.Info("%s service is starting...", zap.String("env", env.SRV_ENV))

	// 3. Start gRPC Server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%%s", env.%s_SRV_PORT))
	if err != nil {
		logger.Fatal("failed to listen", zap.Error(err))
	}

	s := grpc.NewServer()
	// TODO: Register protobuf services

	logger.Info("Service listening", zap.String("port", env.%s_SRV_PORT))
	if err := s.Serve(lis); err != nil {
		logger.Fatal("failed to serve", zap.Error(err))
	}
}
`, modulePath(), modulePath(), service, service, serviceEnvPrefix(prefix), serviceEnvPrefix(prefix))
}

func envTemplate(service string) string {
	pkg := strings.ReplaceAll(service, "-", "")
	parts := strings.Split(service, "-")
	prefix := strings.ToUpper(parts[0])
	return fmt.Sprintf(`package %s

import (
	"log"
	"os"
	"strconv"
)

type Env struct {
	// Service settings
	SRV_ENV       string `+"`mapstructure:\"SRV_ENV\"`"+`
	%s_SRV_NAME string `+"`mapstructure:\"%s_SRV_NAME\"`"+`
	%s_SRV_PORT string `+"`mapstructure:\"%s_SRV_PORT\"`"+`

	// Database settings
	DBHost     string `+"`mapstructure:\"POSTGRES_HOST\"`"+`
	DBPort     string `+"`mapstructure:\"POSTGRES_PORT\"`"+`
	DBUser     string `+"`mapstructure:\"POSTGRES_USER\"`"+`
	DBPassword string `+"`mapstructure:\"POSTGRES_PASSWORD\"`"+`
	DBName     string `+"`mapstructure:\"POSTGRES_DB\"`"+`

	// Redis settings
	RedisHOST     string `+"`mapstructure:\"REDIS_HOST\"`"+`
	RedisPort     int    `+"`mapstructure:\"REDIS_PORT\"`"+`
	RedisPassword string `+"`mapstructure:\"REDIS_PASSWORD\"`"+`
	RedisDB       int    `+"`mapstructure:\"REDIS_DB\"`"+`
}

func getString(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Printf("Invalid integer for %%s: %%s, using default %%d", key, valueStr, defaultValue)
		return defaultValue
	}
	return value
}

func GetEnv() (*Env, error) {
	env := Env{
		SRV_ENV:       getString("SRV_ENV", "development"),
		%s_SRV_NAME: getString("%s_SRV_NAME", "%s"),
		%s_SRV_PORT: getString("%s_SRV_PORT", "9090"),
		DBHost:        getString("POSTGRES_HOST", "postgres-db"),
		DBPort:        getString("POSTGRES_PORT", "5432"),
		DBUser:        getString("POSTGRES_USER", "postgres"),
		DBPassword:    getString("POSTGRES_PASSWORD", "password"),
		DBName:        getString("POSTGRES_DB", "%sdb"),
		RedisHOST:     getString("REDIS_HOST", "localhost"),
		RedisPort:     getInt("REDIS_PORT", 6379),
		RedisPassword: getString("REDIS_PASSWORD", ""),
		RedisDB:       getInt("REDIS_DB", 0),
	}
	return &env, nil
}
`, pkg, prefix, prefix, prefix, prefix, prefix, prefix, service, prefix, prefix, service)
}

func migrateTemplate(service, modPath string) string {
	return fmt.Sprintf(`package migrations

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/lib/pq"
	goose "github.com/pressly/goose/v3"
	svc "%s/services/%s"
)

type migrator struct {
	env svc.Env
}

func NewMigrator(env svc.Env) *migrator {
	return &migrator{env: env}
}

func (f *migrator) Migrate(ctx context.Context, dir string) error {
	dsn := fmt.Sprintf("postgres://%%s:%%s@%%s:%%s/postgres?sslmode=disable",
		f.env.DBUser, f.env.DBPassword, f.env.DBHost, f.env.DBPort)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to Postgres: %%w", err)
	}
	defer db.Close()
	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %%s", f.env.DBName))
	if err != nil {
		if !strings.Contains(err.Error(), "already exists") {
			return fmt.Errorf("failed to create database: %%w", err)
		}
	}
	db.Close()

	dsn = fmt.Sprintf("postgres://%%s:%%s@%%s:%%s/%%s?sslmode=disable",
		f.env.DBUser, f.env.DBPassword, f.env.DBHost, f.env.DBPort, f.env.DBName)
	db, err = sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to DB: %%w", err)
	}
	defer db.Close()

	if err := goose.Up(db, dir); err != nil {
		return err
	}
	return nil
}
`, modPath, service)
}

func interceptorTemplate() string {
	return `package interceptor

import (
	"context"
	"google.golang.org/grpc"
)

// LoggingInterceptor is a placeholder unary interceptor.
func LoggingInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	// TODO: add request logging here
	return handler(ctx, req)
}
`
}

func readmeTemplate(service string) string {
	title := cases.Title(language.Und).String(service)

	return fmt.Sprintf(
		"# %s Service\n\n"+
			"This is the **%s microservice**, generated using a Clean Architecture skeleton.\n\n"+
			"## Structure\n\n"+
			"- `cmd/` – service entrypoint (gRPC)\n"+
			"- `internal/` – API handlers, usecases, repositories\n"+
			"- `migrations/` – SQL migrations + runner\n"+
			"- `docs/` – documentation\n"+
			"- `test/` – tests\n\n"+
			"## Commands\n\n"+
			"```bash\n"+
			"# Customize run/build targets for your service\n"+
			"make test\n"+
			"```\n",
		title, service,
	)
}

// modulePath reads the go.mod in repo root to get the module path.
func modulePath() string {
	data, err := ioutil.ReadFile("go.mod")
	if err != nil {
		return "github.com/your/module"
	}
	lines := strings.Split(string(data), "\n")
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if strings.HasPrefix(l, "module ") {
			return strings.TrimSpace(strings.TrimPrefix(l, "module "))
		}
	}
	return "github.com/your/module"
}

// serviceEnvPrefix returns sanitized prefix for env variable fields.
func serviceEnvPrefix(prefix string) string { return prefix }
