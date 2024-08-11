package main

import (
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"
	"zr/configs"
	"zr/models"

	_ "github.com/lib/pq"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/dialects/postgresql"
)

var DB reform.DB
var Configs configs.Configs
var ServerIsClosed = make(chan bool)

func main() {
	// read config to env var
	viper.SetConfigName("configs")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Can't read config message: ", err)
	}
	err = viper.Unmarshal(&Configs)
	if err != nil {
		log.Fatal("Can't unmarshal config", err)
	}

	// init server and routing
	app := fiber.New()
	app.Get("/list", getNews)
	app.Post("/edit", editNews)
	go func() {
		defer close(ServerIsClosed)
		log.Printf("app listen and serv: %s", app.Listen(":3000"))
	}()

	// init DB
	dbURL := fmt.Sprintf(
		"%s://%s:%s@%s:%s/%s?sslmode=%s",
		Configs.DRIVER,
		Configs.DB_USER,
		Configs.DB_PASSWORD,
		Configs.HOST_DB,
		Configs.DB_PORT,
		Configs.DB_NAME,
		Configs.SSLmode,
	)
	poolSqlDB, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Can't connect to database", err)
	}
	errP := poolSqlDB.Ping()
	if errP != nil {
		log.Fatal("sqlDB not connected message: ", errP)
		os.Exit(2)
	} else {
		defer poolSqlDB.Close()
	}
	// set max pool connections 10 and maxIdle con 10, and one hour lifetime conn
	poolSqlDB.SetMaxOpenConns(10)
	poolSqlDB.SetConnMaxIdleTime(10)
	poolSqlDB.SetConnMaxLifetime(time.Hour)
	logger := log.New(os.Stderr, "SQL: ", log.Flags())
	// init glob var with New reform DB orm
	DB := reform.NewDB(poolSqlDB, postgresql.Dialect, reform.NewPrintfLogger(logger.Printf))
	log.Println("Db connected ", DB.String())

	// end the programm server close
	slog.Info("Server is closed", "info", <-ServerIsClosed)
}

// get news list
func getNews(c *fiber.Ctx) error {
	news_dbs, err := DB.Query("SELECT id, title, content FROM news;")
	if err != nil {
		c.Append("Secces", "false")
		return c.JSON(http.StatusNotFound)
	}
	defer news_dbs.Close()
	var newses []models.News
	for news_dbs.Next()  {
		var news models.News
		if err:= news_dbs.Scan(&news.Id, &news.Title, &news.Content); err != nil {
			c.Append("Secces", "false")
			return c.JSON(http.StatusNotFound)
		}
		catQuery := fmt.Sprintf("SELECT newsId, categoriesId FROM categories WHERE newsId = %d", news.Id)
		cotegoryRow, err := DB.Query(catQuery);
		if err != nil {
			c.Append("Secces", "false")
			return c.JSON(http.StatusNotFound)
		}
		news.Categories = make([]int, 0)
		for cotegoryRow.Next() {
			var cat int
			var some int
			err := cotegoryRow.Scan(some, &cat)
			if err != nil {
				c.Append("Secces", "false")
				return c.JSON(http.StatusNotFound)
			}
			news.Categories = append(news.Categories, cat)
			cat = 0
			some = 0
		}

	}
	c.Append("Sacces", "true")
	return c.JSON(newses)
}

// post news
func editNews(c *fiber.Ctx) error {
	tx, err := DB.Begin()
	if err != nil {
		log.Fatal("Can't exect by transaction", err)
		return err
	}
	editNews := new(models.News)
	if err := c.BodyParser(editNews); err != nil {
		log.Fatal("Can't parce news from body for edit", err)
	}
	news_id := editNews.Id
	newsDB := new(models.News_db)
	newsDB.Id = editNews.Id
	newsDB.Title = editNews.Title
	newsDB.Content = editNews.Content
	query := "INSERT INTO newscategories VALUES(?, ?)"
	errs := DB.Save(newsDB)
	if errs != nil {
		tx.Rollback()
		c.Append("error", errs.Error())
		return c.JSON(http.StatusNotModified)
	}
	for _, item := range editNews.Categories {
		_, err := DB.Exec(query, news_id, item)
		if err != nil {
			log.Fatal("didn't save: ", query, news_id, item)
			tx.Rollback()
			c.Append("error", err.Error())
			return c.JSON(http.StatusNotModified)
		}
	}
	return c.JSON(http.StatusOK)
}
