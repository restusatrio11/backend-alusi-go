package worker

import (
	"context"
	"crypto/tls"
	"net/http"
	"sync"
	"time"

	"backend-alusi-go/internal/domain"
	"backend-alusi-go/internal/repository/postgres"

	"github.com/rs/zerolog/log"
)

type HealthProbeWorker struct {
	appRepo         *postgres.AppRepo
	statusCheckRepo *postgres.StatusCheckRepo
	httpClient      *http.Client
	interval        time.Duration
	ctx             context.Context
	cancel          context.CancelFunc
}

func NewHealthProbeWorker(
	appRepo *postgres.AppRepo,
	statusCheckRepo *postgres.StatusCheckRepo,
	interval time.Duration,
) *HealthProbeWorker {
	if interval <= 0 {
		interval = 5 * time.Minute
	}

	ctx, cancel := context.WithCancel(context.Background())

	// Custom HTTP client with timeout and TLS ignore for internal self-signed certs
	tr := &http.Transport{
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
		MaxIdleConns:        50,
		IdleConnTimeout:     30 * time.Second,
		DisableKeepAlives:   false,
	}

	httpClient := &http.Client{
		Transport: tr,
		Timeout:   5 * time.Second,
	}

	return &HealthProbeWorker{
		appRepo:         appRepo,
		statusCheckRepo: statusCheckRepo,
		httpClient:      httpClient,
		interval:        interval,
		ctx:             ctx,
		cancel:          cancel,
	}
}

// Start initiates the recurring probe loop
func (w *HealthProbeWorker) Start() {
	go func() {
		// Run initial probe 5 seconds after startup
		select {
		case <-w.ctx.Done():
			return
		case <-time.After(5 * time.Second):
			w.ProbeAll(w.ctx)
		}

		ticker := time.NewTicker(w.interval)
		defer ticker.Stop()

		for {
			select {
			case <-w.ctx.Done():
				return
			case <-ticker.C:
				w.ProbeAll(w.ctx)
			}
		}
	}()
}

// ProbeSingleApp probes a specific application and records status
func (w *HealthProbeWorker) ProbeSingleApp(ctx context.Context, app *domain.App) (*domain.StatusCheck, error) {
	start := time.Now()

	req, err := http.NewRequestWithContext(ctx, http.MethodHead, app.URL, nil)
	var resp *http.Response
	var statusCode *int
	var errMsg *string
	var statusResult string

	if err != nil {
		msg := err.Error()
		errMsg = &msg
		statusResult = "down"
	} else {
		req.Header.Set("User-Agent", "PortalBPS-HealthProbe/1.0")
		resp, err = w.httpClient.Do(req)

		// Fallback to GET if HEAD method is not allowed (405 Method Not Allowed)
		if resp != nil && resp.StatusCode == http.StatusMethodNotAllowed {
			resp.Body.Close()
			reqGet, _ := http.NewRequestWithContext(ctx, http.MethodGet, app.URL, nil)
			reqGet.Header.Set("User-Agent", "PortalBPS-HealthProbe/1.0")
			resp, err = w.httpClient.Do(reqGet)
		}

		latency := time.Since(start)
		latencyMS := int(latency.Milliseconds())

		if err != nil {
			msg := err.Error()
			errMsg = &msg
			statusResult = "down"
		} else {
			defer resp.Body.Close()
			code := resp.StatusCode
			statusCode = &code

			// Determine status
			if code >= 200 && code < 400 {
				if latencyMS >= 2000 {
					statusResult = "degraded"
				} else {
					statusResult = "online"
				}
			} else if code == http.StatusUnauthorized || code == http.StatusForbidden {
				// Auth walls / SSO redirects are considered reachable (online)
				statusResult = "online"
			} else if code == http.StatusServiceUnavailable {
				statusResult = "maintenance"
			} else {
				statusResult = "down"
			}
		}
	}

	latencyMS := int(time.Since(start).Milliseconds())

	check := &domain.StatusCheck{
		AppID:          app.ID,
		Status:         statusResult,
		HTTPStatusCode: statusCode,
		LatencyMS:      latencyMS,
		ErrorMessage:   errMsg,
		CheckedAt:      time.Now(),
	}

	// Record to database
	if w.statusCheckRepo != nil {
		_ = w.statusCheckRepo.Record(ctx, check)
	}

	// Update app service status if not locked in maintenance
	if w.appRepo != nil && app.StatusLayanan != "maintenance" {
		_ = w.appRepo.UpdateStatus(ctx, app.ID, statusResult)
	}

	return check, nil
}

// ProbeAll probes all active apps concurrently
func (w *HealthProbeWorker) ProbeAll(ctx context.Context) {
	if w.appRepo == nil {
		return
	}

	filter := domain.AppFilter{
		AktifOnly: true,
		Limit:     500,
		Offset:    0,
	}

	apps, _, err := w.appRepo.List(ctx, filter)
	if err != nil || len(apps) == 0 {
		return
	}

	log.Debug().Int("count", len(apps)).Msg("Starting automated health probe for active applications...")

	sem := make(chan struct{}, 5) // Max 5 concurrent probes
	var wg sync.WaitGroup

	for i := range apps {
		app := apps[i]
		wg.Add(1)

		go func(a domain.App) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			probeCtx, cancel := context.WithTimeout(ctx, 6*time.Second)
			defer cancel()

			_, _ = w.ProbeSingleApp(probeCtx, &a)
		}(app)
	}

	wg.Wait()
	log.Debug().Msg("Automated health probe completed")
}

// Stop terminates the probe worker
func (w *HealthProbeWorker) Stop() {
	w.cancel()
}
