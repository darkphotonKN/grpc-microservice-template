package notification

import "context"

type NotificationService interface {
	SendMessage(ctx context.Context, message string) error
}
