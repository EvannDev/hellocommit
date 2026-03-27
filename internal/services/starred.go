package services

import (
	"context"
	"fmt"
	"log"

	"github.com/hellocommit/api/internal/models"
	"github.com/hellocommit/api/internal/repositories"
	"github.com/hellocommit/api/pkg/github"
)

type StarredService struct {
	userRepo     *repositories.UserRepository
	repoRepo     *repositories.RepoRepository
	githubClient *github.Client
}

func NewStarredService(userRepo *repositories.UserRepository, repoRepo *repositories.RepoRepository, githubClient *github.Client) *StarredService {
	return &StarredService{
		userRepo:     userRepo,
		repoRepo:     repoRepo,
		githubClient: githubClient,
	}
}

func (s *StarredService) GetStarredRepos(ctx context.Context, userID int64) ([]*models.Repo, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	client := github.NewClient(user.AccessToken)
	log.Printf("[starred] fetching starred repos for user %d (%s) from GitHub", userID, user.Username)
	starred, err := client.ListStarredRepos(ctx, user.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch starred repos: %w", err)
	}
	log.Printf("[starred] got %d starred repos for user %d", len(starred), userID)

	var repos []*models.Repo
	var seenIDs []int64
	for _, sr := range starred {
		r := sr.GetRepository()
		repo := &models.Repo{
			UserID:   userID,
			RepoID:   r.GetID(),
			Owner:    r.GetOwner().GetLogin(),
			Name:     r.GetName(),
			FullName: r.GetFullName(),
			HTMLURL:  r.GetHTMLURL(),
			Stars:    r.GetStargazersCount(),
			Language: r.GetLanguage(),
		}
		if err := s.repoRepo.Upsert(repo); err != nil {
			return nil, err
		}
		seenIDs = append(seenIDs, repo.RepoID)
		repos = append(repos, repo)
	}

	if err := s.repoRepo.DeleteNotInList(userID, seenIDs); err != nil {
		return nil, err
	}

	return repos, nil
}

func (s *StarredService) GetUser(userID int64) (*models.User, error) {
	return s.userRepo.GetByID(userID)
}

func (s *StarredService) GetCachedRepos(userID int64) ([]*models.Repo, error) {
	return s.repoRepo.GetByUserID(userID)
}

func (s *StarredService) SyncStarredRepos(ctx context.Context, userID int64) ([]*models.Repo, error) {
	log.Printf("[starred] starting full sync for user %d", userID)

	if err := s.userRepo.UpdateLastSync(userID); err != nil {
		return nil, err
	}

	return s.GetStarredRepos(ctx, userID)
}
