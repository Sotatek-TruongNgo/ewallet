package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/truongnqse05461/ewallet/internal/log"
	"go.uber.org/atomic"
)

type HealthHandler struct {
	isShuttingDown *atomic.Bool
	logger         log.Logger
	db             *sqlx.DB
}

func NewHealthHandler(
	logger log.Logger,
	db *sqlx.DB,
) *HealthHandler {
	return &HealthHandler{
		isShuttingDown: atomic.NewBool(false),
		logger:         logger,
		db:             db,
	}
}

// @Tags SRE
// @Summary readiness probe
// @Id readiness
// @Router /readyz [get]
// @Description Is it a good idea to send traffic to this Pod right now?
// @version 1.0
// @Success 200 string string
// @Failure 503 string string
func (h *HealthHandler) Ready(c *gin.Context) {
	// Don't depend on downstream dependencies, such as other services or databases in your probe.
	// If the dependency has a hickup, or for example, a database is restarted,
	// removing your Pods from the Load Balancers rotation will likely only make the downtime worse.
	if err := h.db.PingContext(c.Request.Context()); err != nil {
		h.logger.WithField("downstream", "database").WithErr(err).Error("ping database failed")
	}

	// The kubelet uses readiness probes to know when a container is ready to start accepting traffic.
	// A Pod is considered ready when all of its containers are ready.
	// One use of this signal is to control which Pods are used as backends for Services.
	// When a Pod is not ready, it is removed from Service load balancers.
	if h.isShuttingDown.Load() {
		c.JSON(http.StatusServiceUnavailable, http.StatusText(http.StatusServiceUnavailable))
	} else {
		c.JSON(http.StatusOK, http.StatusText(http.StatusOK))
	}
}

// @Tags SRE
// @Summary liveness probe
// @Id liveness
// @Router /livez [get]
// @Description Is the container healthy right now, or do we need to restart it?
// @version 1.0
// @Success 200 string string
func (h *HealthHandler) Live(c *gin.Context) {
	// The kubelet uses liveness probes to know when to restart a container.
	// For example, liveness probes could catch a deadlock,
	// where an application is running, but unable to make progress.
	// Restarting a container in such a state can help to make the application more available despite bugs.
	c.JSON(http.StatusOK, http.StatusText(http.StatusOK))
}

func (h *HealthHandler) SetShuttingDown() {
	h.isShuttingDown.Store(true)
}
