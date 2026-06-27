package credentials

import "errors"

type fileStore struct{}

func NewStore() Store {
	return &fileStore{}
}

func (f *fileStore) Save(username, password string) error {
	return errors.New("credential store is not implemented yet")
}

func (f *fileStore) Load(username string) (string, error) {
	return "", errors.New("credential store is not implemented yet")
}

func (f *fileStore) Delete(username string) error {
	return errors.New("credential store is not implemented yet")
}
