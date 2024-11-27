package model

type User struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Hash  string `json:"hash"`
}

func NewUser(name string, email string, password string) (*User, error) {
	hash, err := GenerateHash(&password)
	if err != nil {
		return nil, err
	}

	return &User{
		Name:  name,
		Email: email,
		Hash:  *hash,
	}, nil
}
