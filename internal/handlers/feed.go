package handlers

import (
	"encoding/xml"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/hellocommit/api/internal/services"
)

type FeedHandler struct {
	issueService *services.IssueService
}

func NewFeedHandler(issueService *services.IssueService) *FeedHandler {
	return &FeedHandler{issueService: issueService}
}

type rssChannel struct {
	XMLName       xml.Name  `xml:"channel"`
	Title         string    `xml:"title"`
	Link          string    `xml:"link"`
	Description   string    `xml:"description"`
	LastBuildDate string    `xml:"lastBuildDate"`
	Items         []rssItem `xml:"item"`
}

type rssItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	Author      string `xml:"author"`
	PubDate     string `xml:"pubDate"`
	GUID        string `xml:"guid"`
}

type rssFeed struct {
	XMLName xml.Name   `xml:"rss"`
	Version string     `xml:"version,attr"`
	Channel rssChannel `xml:"channel"`
}

func (h *FeedHandler) GetRSS(c fiber.Ctx) error {
	userID := c.Locals("userID").(int64)

	issues, err := h.issueService.GetGoodFirstIssues(c.Context(), userID)
	if err != nil {
		log.Printf("[GetRSS] error for user %d: %v", userID, err)
		return c.Status(fiber.StatusInternalServerError).SendString("failed to get issues")
	}

	items := make([]rssItem, 0, len(issues))
	for _, issue := range issues {
		items = append(items, rssItem{
			Title:       issue.Title,
			Link:        issue.HTMLURL,
			Description: issue.Body,
			Author:      issue.Author,
			PubDate:     issue.CreatedAt.UTC().Format(time.RFC1123Z),
			GUID:        issue.HTMLURL,
		})
	}

	feed := rssFeed{
		Version: "2.0",
		Channel: rssChannel{
			Title:         fmt.Sprintf("HelloCommit — Good First Issues (user %d)", userID),
			Link:          "https://hellocommit.app",
			Description:   "Good first issues from your starred GitHub repositories",
			LastBuildDate: time.Now().UTC().Format(time.RFC1123Z),
			Items:         items,
		},
	}

	output, err := xml.MarshalIndent(feed, "", "  ")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("failed to generate feed")
	}

	c.Set(fiber.HeaderContentType, "application/rss+xml; charset=utf-8")
	return c.Send(append([]byte(xml.Header), output...))
}
