package credentials

type Store interface {
	Save(username, password string) error
	Load(username string) (string, error)
	Delete(username string) error
}
