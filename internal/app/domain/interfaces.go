package domain

type RepositoryStore interface {
	Save(id, url string) error
	Get(id string) (string, bool)
}
