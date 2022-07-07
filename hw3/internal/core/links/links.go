package links

import (
	"context"
	"fmt"

	"github.com/golang-jwt/jwt"
)

// Repository holds datastore
type Repository interface {
	Create(ctx context.Context, link Link) error
	FindByToken(ctx context.Context, token string) (*Link, error)
	FindAllByToken(ctx context.Context, token string) (StatList, error)
	SaveStat(ctx context.Context, id int64, ip string) error
}

// Links struct holds the Repository instance
type Links struct {
	repo      Repository
	token     Token
	jwtSecret []byte
}

// Token holds method for generate token
type Token interface {
	Generate() string
}

// Payload struct holds JWT settings
type Payload struct {
	jwt.StandardClaims

	Name string
}

// NewService creates a new Workspaces service
func NewService(repo Repository, t Token, jwtSecret string) *Links {
	return &Links{repo: repo, token: t, jwtSecret: []byte(jwtSecret)}
}

// NewModel create new instance Link model
func (l Links) NewModel() *Link {
	return &Link{Token: l.token.Generate()}
}

// Create save entity
func (l Links) Create(ctx context.Context, link Link) error {
	if err := link.Validate(); err != nil {
		return fmt.Errorf("link's validate failed: %w", err)
	}

	if err := l.repo.Create(ctx, link); err != nil {
		return fmt.Errorf("failed to create link in repo: %w", err)
	}

	return nil
}

// FindByToken find by token link
func (l Links) FindByToken(ctx context.Context, token string) (*Link, error) {
	mlink, err := l.repo.FindByToken(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("link not found in repo: %w", err)
	}

	return mlink, nil
}

// SaveStat save link statistics
func (l Links) SaveStat(ctx context.Context, id int64, ip string) error {
	return l.repo.SaveStat(ctx, id, ip)
}

// GetJWTToken create JwtToken
func (l Links) GetJWTToken() (*string, error) {
	payload := Payload{
		StandardClaims: jwt.StandardClaims{},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &payload)

	signedToken, err := token.SignedString(l.jwtSecret)
	if err != nil {
		return nil, err
	}

	return &signedToken, nil
}

func (l Links) FindAllByToken(ctx context.Context, token string) (StatList, error) {
	return l.repo.FindAllByToken(ctx, token)
}
