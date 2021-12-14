package shortener

type URLCache interface {
	Set(shortURL, longURL string) error
	Get(shortURL string) (string, error)
}

type URLService interface {
	CreateShortURL(longURL string) (string, error)
	GetLongURL(shortURL string) (string, error)
}

type URLCreationRequest struct {
	LongURL string `json:"long_url" binding:"required"`
}

