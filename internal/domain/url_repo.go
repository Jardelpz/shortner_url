package domain

type UrlRepository interface {
	Insert(url *Url) error
	Find(shortUrl string) (*Url, error)
}
