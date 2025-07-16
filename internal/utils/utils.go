package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/inflection"
	"github.com/samvibes/vexop/auth-service/internal/models"
)

func GetCurrentUser(c *gin.Context) *models.User {
	userVar, exists := c.Get(UserContextKey)
	if !exists {
		return nil
	}
	user := userVar.(models.User)
	return &user
}

var irregularPlurals = map[string]string{
	"people":    "person",
	"data":      "data",
	"companies": "company",
}

func Singularize(word string) string {
	word = strings.ToLower(word)

	return inflection.Singular(word)

	// // Check irregular first
	// if singular, ok := irregularPlurals[word]; ok {
	// 	return singular
	// }

	// // Basic English rules
	// switch {
	// case strings.HasSuffix(word, "ies"):
	// 	// companies → company
	// 	return strings.TrimSuffix(word, "ies") + "y"
	// case strings.HasSuffix(word, "ses"):
	// 	// processes → process
	// 	return strings.TrimSuffix(word, "es")
	// case strings.HasSuffix(word, "s") && len(word) > 1:
	// 	// users → user, files → file
	// 	return strings.TrimSuffix(word, "s")
	// default:
	// 	return word
	// }
}

func GenerateInviteToken() (rawToken string, hashedToken string, err error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", "", err
	}

	rawToken = base64.RawURLEncoding.EncodeToString(b)

	// Hash the token
	hash := sha256.Sum256([]byte(rawToken))
	hashedToken = hex.EncodeToString(hash[:])
	err = nil

	return
}
