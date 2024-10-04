package notifiler

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/lidofinance/finding-forwarder/generated/databus"
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

const maxTelegramMessageLength = 4096
const warningTelegramMessage = "Warn: Msg >=4096, pls review description message"

func (u *telegram) SendFinding(ctx context.Context, alert *databus.FindingDtoJson) error {
	message := TruncateMessageWithAlertID(
		fmt.Sprintf("%s\n\n%s", alert.Name, FormatAlert(alert, u.source)),
		maxTelegramMessageLength,
		warningTelegramMessage,
	)

	if alert.Severity != databus.SeverityUnknown {
		m := EscapeMarkdownV1(message)

		if sendErr := u.send(ctx, m, true); sendErr != nil {
			message += "\n\nWarning: Could not send msg as markdown"
			return u.send(ctx, message, false)
		}

		return nil
	}

	return u.send(ctx, message, false)
}

func (u *telegram) send(ctx context.Context, message string, useMarkdown bool) error {
	requestURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage?disable_web_page_preview=true&disable_notification=true&chat_id=-%s&text=%s", u.botToken, u.chatID, url.QueryEscape(message))
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
		fmt.Println(resp)
		return fmt.Errorf("received from telegram non-200 response code: %v", resp.Status)
	}

	return nil
}
