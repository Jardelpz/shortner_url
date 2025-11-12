package model

type UrlRepo interface {
	Insert(*Url) error
	Find(shortUrl string) (*Url, error)
}
