package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/turnes/go-github/internal/core/domain"
	"github.com/turnes/go-github/internal/core/ports"
)

type HTTPHandler struct {
	service ports.RepositoryService
}

func NewHTTPHandler(service ports.RepositoryService) *HTTPHandler {
	return &HTTPHandler{service: service}
}

func (h *HTTPHandler) RegisterRoutes(r *gin.Engine) {
	r.GET("/repositories/:owner", h.ListRepositories)
	r.POST("/repositories/:owner", h.CreateRepository)
	r.GET("/repositories/:owner/:name", h.GetRepository)
	r.PUT("/repositories/:owner/:name", h.UpdateRepository)
	r.DELETE("/repositories/:owner/:name", h.DeleteRepository)
}

func (h *HTTPHandler) ListRepositories(c *gin.Context) {
	owner := c.Param("owner")

	repos, err := h.service.ListRepositories(c.Request.Context(), owner, "", "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, repos)
}

func (h *HTTPHandler) CreateRepository(c *gin.Context) {
	owner := c.Param("owner")

	var input domain.CreateRepositoryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	repo, err := h.service.CreateRepository(c.Request.Context(), owner, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, repo)
}

func (h *HTTPHandler) GetRepository(c *gin.Context) {
	owner := c.Param("owner")
	name := c.Param("name")

	repo, err := h.service.GetRepository(c.Request.Context(), owner, name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, repo)
}

func (h *HTTPHandler) UpdateRepository(c *gin.Context) {
	owner := c.Param("owner")
	name := c.Param("name")

	var input domain.UpdateRepositoryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	repo, err := h.service.UpdateRepository(c.Request.Context(), owner, name, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, repo)
}

func (h *HTTPHandler) DeleteRepository(c *gin.Context) {
	owner := c.Param("owner")
	name := c.Param("name")

	if err := h.service.DeleteRepository(c.Request.Context(), owner, name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
