package auth

import "golang.org/x/crypto/bcrypt"

func (svc *Service) hashPassword(raw string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(raw), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func (svc *Service) passwordMatches(rawPassword string, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(rawPassword)) == nil
}
