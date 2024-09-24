package notifiler

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/lidofinance/finding-forwarder/generated/forta/models"
	"github.com/lidofinance/finding-forwarder/generated/proto"
	"github.com/lidofinance/finding-forwarder/internal/connectors/metrics"
)

type telegram struct {
	botToken   string
	chatID     string
	httpClient *http.Client
	metrics    *metrics.Store
	source     string
}

func NewTelegram(botToken, chatID string, httpClient *http.Client, metricsStore *metrics.Store, source string) *telegram {
	return &telegram{
		botToken:   botToken,
		chatID:     chatID,
		httpClient: httpClient,
		metrics:    metricsStore,
		source:     source,
	}
}

func (u *telegram) SendFinding(ctx context.Context, alert *proto.Finding) error {
	message := fmt.Sprintf("%s\n\n%s\n\nAlertId: %s\nSource: %s", alert.Name, alert.Description, alert.GetAlertId(), u.source)

	if alert.Severity != proto.Finding_UNKNOWN {
		return u.send(ctx, message, true)
	}

	return u.send(ctx, message, false)
}

func (u *telegram) SendAlert(ctx context.Context, alert *models.Alert) error {
	message := fmt.Sprintf("%s\n\n%s\n\nAlertId: %s\nSource: %s", alert.Name, alert.Description, alert.AlertID, u.source)

	if alert.Severity != models.AlertSeverityUNKNOWN {
		return u.send(ctx, message, true)
	}

	return u.send(ctx, message, false)
}

func (u *telegram) send(ctx context.Context, message string, useMarkdown bool) error {
	//nolint
	requestURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage?chat_id=-%s&text=%s", u.botToken, u.chatID, url.QueryEscape(message))
	if useMarkdown {
		requestURL += `&parse_mode=markdown`
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, http.NoBody)
	if err != nil {
		return fmt.Errorf("could not create telegram request: %w", err)
	}

	start := time.Now()
	resp, err := u.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("could not send telegram request: %w", err)
	}
	defer func() {
		resp.Body.Close()
		duration := time.Since(start).Seconds()
		u.metrics.SummaryHandlers.With(prometheus.Labels{metrics.Channel: `telegram`}).Observe(duration)
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received from telegram non-200 response code: %v", resp.Status)
	}

	return nil
}
