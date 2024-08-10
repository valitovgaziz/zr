package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"zr/configs"
	"zr/models"

	"github.com/gofiber/fiber/v2"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/jackc/pgx/v4/pgxpool"
	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/dialects/postgresql"
)

var DB reform.DB
var Configs configs.Configs
var ServerIsClosed = make(chan bool)

func main() {
	// read config to env var
	err := cleanenv.ReadEnv(&Configs)
	if err != nil {
		log.Fatal("Can't read env var", err)
		os.Exit(2)
	}

	// init server and routing
	app := fiber.New()
	app.Get("/list", getNews)
	app.Post("/edit/:id", editNews)
	go func() {
		defer close(ServerIsClosed)
		log.Printf("app listen and serv: %s", app.Listen(":3000"))
	}()

	// init DB
	dbURL := fmt.Sprintf(
		"%s://%s:%s@%s:%s/%s",
		os.Getenv("DRIVER"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("HOST_DB"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)
	DBconfig, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		log.Fatal("Unable to parse config: %v", err)
		os.Exit(2)
	}

	pool, err := pgxpool.ConnectConfig(context.Background(), DBconfig)
	if err != nil {
		log.Fatal("Unable to connect to database: %v", err)
		os.Exit(2)
	}
	defer pool.Close()
	logger := log.New(os.Stderr, "SQL: ", log.Flags())
	DB := reform.NewDB(pool, postgresql.Dialect, reform.NewPrintfLogger(logger.Printf))

	// sqlDB, err := sql.Open("postgres", dbURL)
	// if err != nil {
	// 	log.Fatal("Can't connect to database", err)
	// }

	// logger := log.New(os.Stderr, "SQL: ", log.Flags())

	// DB := reform.NewDB(sqlDB, postgresql.Dialect, reform.NewPrintfLogger(logger.Printf))
	// log.Println("Db connected ", DB.String())

	slog.Info("Server is closed", "info", <-ServerIsClosed)
}

// get news list
func getNews(c *fiber.Ctx) error {
	newses, err := DB.FindAllFrom(models.NewsTable, "id")
	if err != nil {
		return c.JSON(http.StatusNotFound)
	}
	return c.JSON(newses)
}

// post news
func editNews(c *fiber.Ctx) error {
	editNews := new(models.News)
	if err := c.BodyParser(editNews); err != nil {
		log.Fatal("Can't parce news from body for edit", err)
	}
	err := DB.Save(editNews)
	if err != nil {
		return c.JSON(http.StatusNotModified)
	}
	return c.JSON(http.StatusOK)
}
