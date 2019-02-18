package main

import (
	"time"
)

type email struct {
	Email                      string
	UnsubscribeSecretHex       string
	LastEmailNotificationDate  time.Time
	PendingEmails              int
	SendReplyNotifications     bool
	SendModeratorNotifications bool
}
