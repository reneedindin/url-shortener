package model

type URL struct {
	ID string `json:"id"`
	ShortUrl string `json:"shortUrl"`
}

type ShortURLInfo struct {
	URL string `json:"url"`
	ExpireAt string `json:"expireAt"`
}
