package notifiler

import (
	"fmt"
	"strings"
	"time"

	"github.com/lidofinance/onchain-mon/generated/databus"
)

func FormatAlert(alert *databus.FindingDtoJson, source string) string {
	var (
		body   string
		footer string
	)

	if alert.Description != "" {
		body = alert.Description
		footer += "\n\n"
	}

	footer += fmt.Sprintf("Alert Id: %s", alert.AlertId)
	footer += fmt.Sprintf("\nBot name: %s", alert.BotName)
	footer += fmt.Sprintf("\nTeam: %s", alert.Team)

	if alert.BlockNumber != nil {
		footer += fmt.Sprintf("\nBlock number: [%d](https://etherscan.io/block/%d/)", *alert.BlockNumber, *alert.BlockNumber)
	}

	if alert.TxHash != nil {
		footer += fmt.Sprintf("\nTx hash: [%s](https://etherscan.io/tx/%s/)", shortenHex(*alert.TxHash), *alert.TxHash)
	}

	footer += fmt.Sprintf("\nSource: %s", source)
	footer += fmt.Sprintf("\nServer timestamp: %s", time.Now().Format("15:04:05.000 MST"))
	if alert.BlockTimestamp != nil {
		footer += fmt.Sprintf("\nBlock timestamp:   %s", time.Unix(int64(*alert.BlockTimestamp), 0).Format("15:04:05.000 MST"))
	}

	return fmt.Sprintf("%s%s", body, footer)
}

func shortenHex(input string) string {
	if len(input) <= 5 {
		return input
	}
	return fmt.Sprintf("x%s...%s", input[2:5], input[len(input)-3:])
}

func TruncateMessageWithAlertID(message string, stringLimit int, warnMessage string) string {
	if len(message) <= stringLimit {
		return message
	}

	alertIndex := strings.LastIndex(message, "Alert Id:")
	if alertIndex == -1 {
		return fmt.Sprintf("%s\n%s", warnMessage, message[:stringLimit-len(warnMessage)-1])
	}

	alertText := message[alertIndex:]

	const formatSpecialCharsLength = 9
	maxTextLength := stringLimit - len(warnMessage) - len(alertText) - formatSpecialCharsLength

	if maxTextLength > 0 && alertIndex > maxTextLength {
		return fmt.Sprintf("%s\n...\n\n*%s*\n%s", message[:maxTextLength], warnMessage, alertText)
	}

	return fmt.Sprintf("%s\n%s", warnMessage, alertText)
}
