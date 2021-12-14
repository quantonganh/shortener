package http

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/quantonganh/shortener"
)

type Server struct {
	URLCache   shortener.URLCache
	URLService shortener.URLService
}

func NewServer(redisService shortener.URLCache, cassandraService shortener.URLService) *Server {
	return &Server{
		URLCache:   redisService,
		URLService: cassandraService,
	}
}

func (s *Server) CreateShortURL(c *gin.Context) {
	var urlCreationRequest shortener.URLCreationRequest
	if err := c.ShouldBindJSON(&urlCreationRequest); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"message": "bad request",
		})
		return
	}

	shortURL, err := s.URLService.CreateShortURL(urlCreationRequest.LongURL)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"message": err,
		})
	}

	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	c.IndentedJSON(http.StatusOK, gin.H{
		"message": "created short URL",
		"short_url": fmt.Sprintf("%s://%s/%s", scheme, c.Request.Host, shortURL),
	})
}

func (s *Server) RedirectToOriginalURL(c *gin.Context) {
	shortURL := c.Param("shortURL")
	longURL, _ := s.URLCache.Get(shortURL)
	if longURL == "" {
		longURL, err := s.URLService.GetLongURL(shortURL)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{
				"message": err,
			})
		}

		if err := s.URLCache.Set(shortURL, longURL); err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{
				"message": err,
			})
		}

		c.Redirect(http.StatusFound, longURL)
	} else {
		c.Redirect(http.StatusFound, longURL)
	}
}