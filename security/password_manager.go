package security

import (
	"crypto/sha1"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"github.com/adwitiyaio/arka/dependency"
	"github.com/adwitiyaio/arka/logger"
)

const DependencyPasswordManager = "password_manager"

type PasswordManager interface {
	// HashPassword ... Create a hash for a plain text password
	HashPassword(password string) string
	// VerifyPassword ... Compare a hash with a plain text password
	VerifyPassword(password string, hash string) error

	// CreateSha1Hash ... Create a sha1 hash
	CreateSha1Hash(text string) string
}

func BootstrapPasswordManager() {
	c := dependency.GetManager()
	pm := &simplePasswordManager{}
	c.Register(DependencyPasswordManager, pm)
}

type simplePasswordManager struct{}

func (s simplePasswordManager) HashPassword(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		logger.Log.Error().Err(err).Stack().Msg("failed to hash password")
	}
	return string(hash)
}

func (s simplePasswordManager) VerifyPassword(password string, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return err
	}
	return nil
}

func (s simplePasswordManager) CreateSha1Hash(text string) string {
	h := sha1.New()
	h.Write([]byte(text))
	bs := h.Sum(nil)
	return fmt.Sprintf("%x", bs)
}
