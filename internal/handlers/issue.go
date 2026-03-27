package handlers

import (
	"log"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/hellocommit/api/internal/services"
)

type IssueHandler struct {
	issueService *services.IssueService
}

func NewIssueHandler(issueService *services.IssueService) *IssueHandler {
	return &IssueHandler{issueService: issueService}
}

func (h *IssueHandler) GetIssue(c fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid issue id"})
	}
	issue, err := h.issueService.GetIssueByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "issue not found"})
	}
	return c.JSON(issue)
}

func (h *IssueHandler) GetIssues(c fiber.Ctx) error {
	userID := c.Locals("userID").(int64)
	owner := c.Params("owner")
	repo := c.Params("name")
	labelFilter := c.Query("label")

	issues, err := h.issueService.GetIssues(c.Context(), userID, owner, repo, labelFilter)
	if err != nil {
		log.Printf("[GetIssues] error for user %d: %v", userID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get issues",
		})
	}

	return c.JSON(fiber.Map{
		"issues": issues,
		"count":  len(issues),
	})
}

func (h *IssueHandler) GetGoodFirstIssues(c fiber.Ctx) error {
	userID := c.Locals("userID").(int64)

	issues, err := h.issueService.GetGoodFirstIssues(c.Context(), userID)
	if err != nil {
		log.Printf("[GetGoodFirstIssues] error for user %d: %v", userID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get issues",
		})
	}

	return c.JSON(fiber.Map{
		"issues": issues,
		"count":  len(issues),
	})
}

func (h *IssueHandler) SyncIssues(c fiber.Ctx) error {
	userID := c.Locals("userID").(int64)
	owner := c.Params("owner")
	repo := c.Params("name")

	issues, err := h.issueService.SyncIssues(c.Context(), userID, owner, repo)
	if err != nil {
		log.Printf("[SyncIssues] error for user %d repo %s/%s: %v", userID, owner, repo, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to sync issues",
		})
	}

	return c.JSON(fiber.Map{
		"issues": issues,
		"count":  len(issues),
		"synced": true,
	})
}

func (h *IssueHandler) SyncAllIssues(c fiber.Ctx) error {
	userID := c.Locals("userID").(int64)

	if err := h.issueService.SyncAllIssues(c.Context(), userID); err != nil {
		log.Printf("[SyncAllIssues] error for user %d: %v", userID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to sync issues",
		})
	}

	return c.JSON(fiber.Map{
		"synced":  true,
		"message": "all issues synced successfully",
	})
}

func (h *IssueHandler) DismissIssue(c fiber.Ctx) error {
	userID := c.Locals("userID").(int64)
	id, err := strconv.ParseInt(c.Params("issueId"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid issue id"})
	}
	if err := h.issueService.DismissIssue(c.Context(), id, userID); err != nil {
		log.Printf("[DismissIssue] error for user %d issue %d: %v", userID, id, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to dismiss issue"})
	}
	return c.JSON(fiber.Map{"success": true})
}

func (h *IssueHandler) ToggleBookmark(c fiber.Ctx) error {
	userID := c.Locals("userID").(int64)
	issueID, err := strconv.ParseInt(c.Params("issueId"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid issue id",
		})
	}

	if err := h.issueService.ToggleBookmark(c.Context(), issueID, userID); err != nil {
		log.Printf("[ToggleBookmark] error for user %d issue %d: %v", userID, issueID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to toggle bookmark",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
	})
}
