package main

import (
	"time"
)

type domain struct {
	Domain                  string      `json:"domain"`
	OwnerHex                string      `json:"ownerHex"`
	Name                    string      `json:"name"`
	CreationDate            time.Time   `json:"creationDate"`
	State                   string      `json:"state"`
	ImportedComments        bool        `json:"importedComments"`
	AutoSpamFilter          bool        `json:"autoSpamFilter"`
	RequireModeration       bool        `json:"requireModeration"`
	RequireIdentification   bool        `json:"requireIdentification"`
	ModerateAllAnonymous    bool        `json:"moderateAllAnonymous"`
	Moderators              []moderator `json:"moderators"`
	EmailNotificationPolicy string      `json:"emailNotificationPolicy"`
}
