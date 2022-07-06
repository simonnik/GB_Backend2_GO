package links

import (
	pkgToken "github.com/simonnik/GB_Backend2_GO/hw3/internal/pkg/token"
)

type token struct {
	HashMinLength int
	HashSalt      string
}

func (t token) Generate() string {
	return pkgToken.GenerateToken(t.HashMinLength, t.HashSalt)
}

func NewToken(minL int, salt string) Token {
	return token{
		HashMinLength: minL,
		HashSalt:      salt,
	}
}
