package handlers

import (
	"encoding/json"
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/hellocommit/api/internal/services"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

type CreateUserRequest struct {
	Username    string `json:"username"`
	AccessToken string `json:"access_token"`
}

func (h *UserHandler) Create(c fiber.Ctx) error {
	var req CreateUserRequest
	if err := json.Unmarshal(c.Body(), &req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	user, err := h.userService.CreateOrUpdate(c.Context(), req.Username, req.AccessToken)
	if err != nil {
		log.Printf("[Create] error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to create user",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(user)
}

func (h *UserHandler) Get(c fiber.Ctx) error {
	userID := c.Locals("userID").(int64)

	user, err := h.userService.GetByID(userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "user not found",
		})
	}

	return c.JSON(user)
}

func (h *UserHandler) GetByGitHubID(c fiber.Ctx) error {
	// This handler is not registered in main.go, kept for reference
	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "not found"})
}

type SyncRequest struct {
	GitHubID    int64  `json:"github_id"`
	AccessToken string `json:"access_token"`
}

func (h *UserHandler) Delete(c fiber.Ctx) error {
	userID := c.Locals("userID").(int64)
	if err := h.userService.Delete(userID); err != nil {
		log.Printf("[Delete] error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to delete account"})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *UserHandler) Sync(c fiber.Ctx) error {
	var req SyncRequest
	if err := json.Unmarshal(c.Body(), &req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	user, err := h.userService.CreateOrUpdate(c.Context(), "", req.AccessToken)
	if err != nil {
		log.Printf("[Sync] error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to sync user",
		})
	}

	return c.JSON(fiber.Map{
		"user":    user,
		"message": "user synced successfully",
	})
}
