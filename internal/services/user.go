package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/hellocommit/api/internal/models"
	"github.com/hellocommit/api/internal/repositories"
	"github.com/hellocommit/api/pkg/github"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type UserService struct {
	repo         *repositories.UserRepository
	githubClient *github.Client
}

func NewUserService(repo *repositories.UserRepository, githubClient *github.Client) *UserService {
	return &UserService{repo: repo, githubClient: githubClient}
}

func (s *UserService) CreateOrUpdate(ctx context.Context, username, accessToken string) (*models.User, error) {
	client := github.NewClient(accessToken)
	ghUser, _, err := client.GetAuthenticatedUser(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get github user: %w", err)
	}

	user := &models.User{
		GitHubID:    ghUser.GetID(),
		Username:    ghUser.GetLogin(),
		AccessToken: accessToken,
		LastSyncAt:  time.Now(),
	}

	existingUser, err := s.repo.GetByGitHubID(user.GitHubID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	if existingUser != nil {
		existingUser.AccessToken = accessToken
		existingUser.LastSyncAt = time.Now()
		if err := s.repo.UpdateToken(existingUser.ID, accessToken); err != nil {
			return nil, err
		}
		if err := s.repo.UpdateLastSync(existingUser.ID); err != nil {
			return nil, err
		}
		existingUser.AccessToken = accessToken
		return existingUser, nil
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) GetByID(id int64) (*models.User, error) {
	return s.repo.GetByID(id)
}

func (s *UserService) GetByGitHubID(githubID int64) (*models.User, error) {
	return s.repo.GetByGitHubID(githubID)
}

func (s *UserService) Delete(id int64) error {
	return s.repo.Delete(id)
}
