package handlers

import (
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/hellocommit/api/internal/services"
)

type StarredHandler struct {
	starredService *services.StarredService
	issueService   *services.IssueService
}

func NewStarredHandler(starredService *services.StarredService, issueService *services.IssueService) *StarredHandler {
	return &StarredHandler{starredService: starredService, issueService: issueService}
}

func (h *StarredHandler) GetStarred(c fiber.Ctx) error {
	userID := c.Locals("userID").(int64)

	repos, err := h.starredService.GetCachedRepos(userID)
	if err != nil {
		log.Printf("[GetStarred] error for user %d: %v", userID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get starred repos",
		})
	}

	return c.JSON(fiber.Map{
		"repos": repos,
		"count": len(repos),
	})
}

func (h *StarredHandler) SyncStarred(c fiber.Ctx) error {
	userID := c.Locals("userID").(int64)

	repos, err := h.starredService.SyncStarredRepos(c.Context(), userID)
	if err != nil {
		log.Printf("[SyncStarred] error for user %d: %v", userID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to sync starred repos",
		})
	}

	return c.JSON(fiber.Map{
		"repos":  repos,
		"count":  len(repos),
		"synced": true,
	})
}

func (h *StarredHandler) SyncAll(c fiber.Ctx) error {
	userID := c.Locals("userID").(int64)

	log.Printf("[sync] full sync triggered for user %d", userID)

	repos, err := h.starredService.SyncStarredRepos(c.Context(), userID)
	if err != nil {
		log.Printf("[SyncAll] starred repos error for user %d: %v", userID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to sync starred repos",
		})
	}

	if err := h.issueService.SyncAllIssues(c.Context(), userID); err != nil {
		log.Printf("[SyncAll] issues error for user %d: %v", userID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to sync issues",
		})
	}

	return c.JSON(fiber.Map{
		"repos":   repos,
		"count":   len(repos),
		"synced":  true,
		"message": "starred repos and issues synced successfully",
	})
}
