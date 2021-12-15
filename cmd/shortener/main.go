package main

import (
	"fmt"
	"io/fs"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"github.com/quantonganh/shortener"
	"github.com/quantonganh/shortener/cassandra"
	"github.com/quantonganh/shortener/http"
	"github.com/quantonganh/shortener/redis"
	"github.com/quantonganh/shortener/ui"
)

func main() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AutomaticEnv()
	viper.SetDefault("redis.addr", "localhost:6379")
	viper.SetDefault("cassandra.hosts", []string{"localhost"})
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Fatal(err)
		}
	}

	var cfg *shortener.Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatal(err)
	}

	rdb, err := redis.NewDB(cfg.Redis.Addr)
	if err != nil {
		log.Fatal(err)
	}
	urlCache := redis.NewURLCache(rdb)

	cdb := cassandra.NewDB(cfg.Cassandra.Hosts...)
	if err := cdb.Open(); err != nil {
		log.Fatal(err)
	}

	urlService := cassandra.NewURLService(cdb)
	s := http.NewServer(urlCache, urlService)

	publicFS, err := fs.Sub(ui.Public, "public")
	if err != nil {
		log.Fatal(err)
	}

	_ = fs.WalkDir(publicFS, ".", func(path string, d fs.DirEntry, err error) error {
		fmt.Println(path)
		return nil
	})

	r := gin.Default()
	r.GET("/", gin.WrapH(http.UIHandler(publicFS)))
	r.GET("/favicon.png", gin.WrapH(http.UIHandler(publicFS)))
	r.GET("/global.css", gin.WrapH(http.UIHandler(publicFS)))
	r.GET("/build/bundle.css", gin.WrapH(http.UIHandler(publicFS)))
	r.GET("/build/bundle.js", gin.WrapH(http.UIHandler(publicFS)))
	r.GET("/build/bundle.js.map", gin.WrapH(http.UIHandler(publicFS)))
	r.POST("/create-short-url", s.CreateShortURL)
	r.GET("/:shortURL", s.RedirectToOriginalURL)

	if err := r.Run(); err != nil {
		log.Fatal(err)
	}
}