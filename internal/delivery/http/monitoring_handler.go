package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"backend-alusi-go/internal/delivery/http/response"
	"backend-alusi-go/internal/usecase"
	"backend-alusi-go/pkg/realtime"

	"github.com/gin-gonic/gin"
)

type MonitoringHandler struct {
	monitoringUsecase *usecase.MonitoringUsecase
	sseHub            *realtime.SSEHub
}

func NewMonitoringHandler(monitoringUsecase *usecase.MonitoringUsecase, sseHub *realtime.SSEHub) *MonitoringHandler {
	return &MonitoringHandler{
		monitoringUsecase: monitoringUsecase,
		sseHub:            sseHub,
	}
}

// GetAppStatusHistory godoc
// @Summary      Riwayat status kesehatan aplikasi
// @Description  Mengambil riwayat log ping & health check uptime sebuah aplikasi
// @Tags         Monitoring
// @Produce      json
// @Param        slug   path      string  true  "Slug atau ID aplikasi"
// @Param        limit  query     int     false "Batas jumlah riwayat (default 30)"
// @Success      200    {object}  response.StandardResponse
// @Failure      400    {object}  response.StandardResponse
// @Failure      404    {object}  response.StandardResponse
// @Router       /apps/{slug}/status-history [get]
func (h *MonitoringHandler) GetAppStatusHistory(c *gin.Context) {
	slugOrID := c.Param("slug")
	if slugOrID == "" {
		slugOrID = c.Param("id")
	}
	if slugOrID == "" {
		response.BadRequest(c, "Identifier aplikasi tidak valid.", nil)
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "30"))
	history, err := h.monitoringUsecase.GetAppStatusHistoryByIdentifier(c.Request.Context(), slugOrID, limit)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Riwayat status aplikasi berhasil dimuat", history, nil)
}

// GetServiceUptimeSummary godoc
// @Summary      Ringkasan uptime seluruh layanan (Service Status)
// @Description  Mengambil status operasional, SLA uptime persentase, dan latensi rata-rata seluruh aplikasi
// @Tags         Monitoring
// @Produce      json
// @Param        days  query     int  false "Jumlah hari rekapitulasi (default 30)"
// @Success      200   {object}  response.StandardResponse
// @Failure      500   {object}  response.StandardResponse
// @Router       /services/status [get]
func (h *MonitoringHandler) GetServiceUptimeSummary(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "30"))
	summary, err := h.monitoringUsecase.GetServiceUptimeSummary(c.Request.Context(), days)
	if err != nil {
		response.InternalServerError(c, "Gagal memuat ringkasan uptime layanan: "+err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Ringkasan uptime layanan berhasil dimuat", summary, nil)
}

// ManualProbeApp godoc
// @Summary      Trigger manual probe kesehatan aplikasi
// @Description  Memaksa pemeriksaan langsung (live HTTP ping) pada endpoint health check aplikasi (Admin)
// @Tags         Admin - Monitoring
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "ID Aplikasi"
// @Success      200  {object}  response.StandardResponse
// @Failure      400  {object}  response.StandardResponse
// @Failure      401  {object}  response.StandardResponse
// @Failure      403  {object}  response.StandardResponse
// @Failure      500  {object}  response.StandardResponse
// @Router       /admin/apps/{id}/probe [post]
func (h *MonitoringHandler) ManualProbeApp(c *gin.Context) {
	appID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "ID aplikasi tidak valid.", nil)
		return
	}

	check, err := h.monitoringUsecase.ManualProbeApp(c.Request.Context(), appID)
	if err != nil {
		response.InternalServerError(c, "Gagal melakukan health probe: "+err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Health probe selesai dijalankan", check, nil)
}

// StreamStatusEvents godoc
// @Summary      Realtime Server-Sent Events (SSE) Stream Status Layanan
// @Description  Membuka koneksi HTTP streaming Server-Sent Events (SSE) untuk menerima notifikasi realtime perubahan status layanan aplikasi secara instan
// @Tags         Monitoring
// @Produce      text/event-stream
// @Success      200  {string}  string "text/event-stream data stream"
// @Router       /services/realtime-status [get]
func (h *MonitoringHandler) StreamStatusEvents(c *gin.Context) {
	if h.sseHub == nil {
		response.InternalServerError(c, "Realtime SSE Hub is not initialized")
		return
	}

	w := c.Writer
	flusher, ok := w.(http.Flusher)
	if !ok {
		response.InternalServerError(c, "Streaming is not supported by the client or proxy")
		return
	}

	// Set headers for Server-Sent Events
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Transfer-Encoding", "chunked")
	w.Header().Set("X-Accel-Buffering", "no")
	w.WriteHeader(http.StatusOK)
	flusher.Flush()

	// Register client to Hub
	clientID := fmt.Sprintf("client-%d-%s", time.Now().UnixNano(), c.ClientIP())
	client := &realtime.Client{
		ID:       clientID,
		SendChan: make(chan []byte, 64),
	}
	h.sseHub.RegisterClient(client)
	defer h.sseHub.UnregisterClient(client)

	// Send initial snapshot
	if summary, err := h.monitoringUsecase.GetServiceUptimeSummary(c.Request.Context(), 30); err == nil {
		snapshotBytes, _ := json.Marshal(map[string]interface{}{
			"type":      "snapshot",
			"timestamp": time.Now(),
			"services":  summary,
		})
		_, _ = fmt.Fprintf(w, "event: snapshot\ndata: %s\n\n", string(snapshotBytes))
		flusher.Flush()
	}

	// Keep connection alive with periodic pings every 15s
	pingTicker := time.NewTicker(15 * time.Second)
	defer pingTicker.Stop()

	ctx := c.Request.Context()
	for {
		select {
		case <-ctx.Done():
			return
		case <-pingTicker.C:
			_, err := fmt.Fprintf(w, ": ping\n\n")
			if err != nil {
				return
			}
			flusher.Flush()
		case msg, ok := <-client.SendChan:
			if !ok {
				return
			}
			_, err := w.Write(msg)
			if err != nil {
				return
			}
			flusher.Flush()
		}
	}
}

