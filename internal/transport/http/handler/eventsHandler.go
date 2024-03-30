package handler

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"permnews/internal/models/entity"
	"permnews/internal/storage/filters"
	"permnews/internal/transport/http/models/output"
	"strconv"
)

type EventsService interface {
	GetAll(filter *filters.EventFilter) ([]*entity.Event, error)
}

type EventsHandler struct {
	s EventsService
}

func NewEventsHandler(s EventsService) *EventsHandler {
	return &EventsHandler{s: s}
}

// GetEvents получить список событий
//
//	@Summary		получить список событий
//	@Description	получить список событий по фильтру
//	@Tags			events
//	@Accept			json
//	@Produce		json
//	@Param			id		query		string	true	"фильтр по id"
//	@Param			source	query		string	true	"фильтр по источнику"
//	@Param			id		category	string	true	"фильтр по категории"
//	@Success		200		{object}
func (h *EventsHandler) GetEvents(c *gin.Context) {
	byId := c.Query("id")
	bySource := c.Query("source")
	byCategory := c.Query("category")
	skip, err := strconv.Atoi(c.Query("skip"))
	if err != nil {
		skip = 0
	}
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		limit = 20
	}
	filter := &filters.EventFilter{
		Id:       byId,
		Source:   bySource,
		Category: byCategory,
		Skip:     skip,
		Limit:    limit,
	}
	res, err := h.s.GetAll(filter)
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	out := make([]output.Event, len(res))
	for i, e := range res {
		out[i] = output.Event{
			Id:       e.Id,
			Title:    e.Title,
			Notice:   e.Notice,
			Source:   e.Source,
			Image:    e.Image,
			Category: e.Category,
		}
	}
	c.JSON(http.StatusOK, out)
	return
}
