package client

import (
	"context"
	"errors"
	"net/http"
	"ratelimiter/internal/httpx"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc *Service
}

func NewClientHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) CreateClient(c *gin.Context) {
	var body CreateClientRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		httpx.BindError(c, err)
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	client, err := h.svc.CreateClientService(ctx, &body)
	if err != nil {
		if errors.Is(err, ErrClientAlreadyExists) {
			httpx.Error(c, http.StatusConflict, MsgClientAlreadyExists)
			return
		}
		httpx.Error(c, http.StatusInternalServerError, httpx.MsgInternalError)
		return
	}
	httpx.Success(c, http.StatusCreated, MsgClientCreated, client)
}

func (h *Handler) GetClientByID(c *gin.Context) {
	clientID := c.Param("id")
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	client, err := h.svc.GetClientByIDService(ctx, clientID)
	if err != nil {
		if errors.Is(err, ErrClientNotFound) {
			httpx.Error(c, http.StatusNotFound, MsgClientNotFound)
			return
		}
		httpx.Error(c, http.StatusInternalServerError, httpx.MsgInternalError)
		return
	}
	httpx.Success(c, http.StatusOK, MsgClientFetched, client)
}

func (h *Handler) GetAllClients(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	clients, err := h.svc.GetAllClientsService(ctx)
	if err != nil {
		httpx.Error(c, http.StatusInternalServerError, httpx.MsgInternalError)
		return
	}
	httpx.Success(c, http.StatusOK, MsgClientsFetched, clients)
}

func (h *Handler) UpdateClient(c *gin.Context) {
	clientID := c.Param("id")
	var body UpdateClientRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		httpx.BindError(c, err)
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	client, err := h.svc.UpdateClientService(ctx, clientID, &body)
	if err != nil {
		if errors.Is(err, ErrClientNotFound) {
			httpx.Error(c, http.StatusNotFound, MsgClientNotFound)
			return
		}
		httpx.Error(c, http.StatusInternalServerError, httpx.MsgInternalError)
		return
	}
	httpx.Success(c, http.StatusOK, MsgClientUpdated, client)
}

func (h *Handler) DeleteClient(c *gin.Context) {
	clientID := c.Param("id")
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	if err := h.svc.DeleteClientByIDService(ctx, clientID); err != nil {
		if errors.Is(err, ErrClientNotFound) {
			httpx.Error(c, http.StatusNotFound, MsgClientNotFound)
			return
		}
		httpx.Error(c, http.StatusInternalServerError, httpx.MsgInternalError)
		return
	}
	httpx.Success(c, http.StatusOK, MsgClientDeleted, nil)
}
