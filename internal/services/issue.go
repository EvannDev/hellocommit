package services

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/hellocommit/api/internal/models"
	"github.com/hellocommit/api/internal/repositories"
	"github.com/hellocommit/api/pkg/github"
)

type IssueService struct {
	issueRepo *repositories.IssueRepository
	repoRepo  *repositories.RepoRepository
	userRepo  *repositories.UserRepository
}

func NewIssueService(issueRepo *repositories.IssueRepository, repoRepo *repositories.RepoRepository, userRepo *repositories.UserRepository) *IssueService {
	return &IssueService{
		issueRepo: issueRepo,
		repoRepo:  repoRepo,
		userRepo:  userRepo,
	}
}

func (s *IssueService) GetIssueByID(ctx context.Context, id int64) (*models.Issue, error) {
	return s.issueRepo.GetByID(id)
}

func (s *IssueService) GetIssues(ctx context.Context, userID int64, owner, repoName string, labelFilter string) ([]*models.Issue, error) {
	repos, err := s.repoRepo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	var repoID int64
	for _, r := range repos {
		if r.Owner == owner && r.Name == repoName {
			repoID = r.ID
			break
		}
	}

	if repoID == 0 {
		return nil, fmt.Errorf("repo not found")
	}

	return s.issueRepo.GetByRepoID(repoID, labelFilter)
}

func (s *IssueService) GetGoodFirstIssues(ctx context.Context, userID int64) ([]*models.Issue, error) {
	return s.issueRepo.GetGoodFirstIssues(userID)
}

func (s *IssueService) SyncIssues(ctx context.Context, userID int64, owner, repoName string) ([]*models.Issue, error) {
	repos, err := s.repoRepo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	var targetRepo *models.Repo
	for _, r := range repos {
		if r.Owner == owner && r.Name == repoName {
			targetRepo = r
			break
		}
	}

	if targetRepo == nil {
		return nil, fmt.Errorf("repo not found in user's starred list")
	}

	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	client := github.NewClient(user.AccessToken)

	since := targetRepo.FetchedAt
	isFirstSync := since.IsZero()
	state := "all"
	if isFirstSync {
		state = "open"
	}

	log.Printf("[issues] fetching issues for %s/%s from GitHub (state=%s, since=%v)", owner, repoName, state, since)
	ghIssues, err := client.ListAllIssues(ctx, owner, repoName, state, since)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch issues from github: %w", err)
	}

	var openIssues []*models.Issue
	var closedNumbers []int

	for _, gi := range ghIssues {
		if gi.IsPullRequest() {
			continue
		}
		if gi.GetState() == "closed" {
			closedNumbers = append(closedNumbers, gi.GetNumber())
			continue
		}
		var labels []string
		for _, l := range gi.Labels {
			labels = append(labels, l.GetName())
		}
		openIssues = append(openIssues, &models.Issue{
			RepoID:      targetRepo.ID,
			IssueNumber: gi.GetNumber(),
			Title:       gi.GetTitle(),
			Body:        gi.GetBody(),
			HTMLURL:     gi.GetHTMLURL(),
			Labels:      strings.Join(labels, ","),
			State:       gi.GetState(),
			Author:      gi.GetUser().GetLogin(),
			AuthorURL:   gi.GetUser().GetHTMLURL(),
			Assignee:    gi.GetAssignee().GetLogin(),
			Comments:    gi.GetComments(),
			CreatedAt:   gi.GetCreatedAt().Time,
		})
	}

	if len(openIssues) > 0 {
		if err := s.issueRepo.UpsertBatch(openIssues); err != nil {
			return nil, err
		}
	}

	if isFirstSync {
		// Full prune on first sync: remove any DB rows not in the open set.
		openNumbers := make([]int, len(openIssues))
		for i, iss := range openIssues {
			openNumbers[i] = iss.IssueNumber
		}
		if err := s.issueRepo.PruneClosedIssues(targetRepo.ID, openNumbers); err != nil {
			return nil, fmt.Errorf("pruning closed issues: %w", err)
		}
	} else if len(closedNumbers) > 0 {
		// Incremental: only delete the specific issues GitHub told us closed.
		if err := s.issueRepo.DeleteByNumbers(targetRepo.ID, closedNumbers); err != nil {
			return nil, fmt.Errorf("deleting closed issues: %w", err)
		}
	}

	if err := s.repoRepo.UpdateFetchedAt(targetRepo.ID, time.Now()); err != nil {
		return nil, err
	}

	log.Printf("[issues] synced %d open issues for %s/%s (%d closed removed)", len(openIssues), owner, repoName, len(closedNumbers))
	return openIssues, nil
}

func (s *IssueService) SyncAllIssues(ctx context.Context, userID int64) error {
	repos, err := s.repoRepo.GetByUserID(userID)
	if err != nil {
		return err
	}

	const maxConcurrency = 20
	log.Printf("[issues] syncing issues for %d repos (user %d) with concurrency %d", len(repos), userID, maxConcurrency)

	sem := make(chan struct{}, maxConcurrency)
	errCh := make(chan error, len(repos))
	var wg sync.WaitGroup

	for _, repo := range repos {
		wg.Add(1)
		sem <- struct{}{}
		go func(owner, name string) {
			defer wg.Done()
			defer func() { <-sem }()
			if _, err := s.SyncIssues(ctx, userID, owner, name); err != nil {
				errCh <- err
			}
		}(repo.Owner, repo.Name)
	}

	wg.Wait()
	close(errCh)

	for err := range errCh {
		return err
	}

	log.Printf("[issues] finished syncing all issues for user %d", userID)
	return nil
}

func (s *IssueService) ToggleBookmark(ctx context.Context, issueID, userID int64) error {
	return s.issueRepo.ToggleBookmark(issueID, userID)
}

func (s *IssueService) DismissIssue(ctx context.Context, issueID, userID int64) error {
	return s.issueRepo.Dismiss(issueID, userID)
}
