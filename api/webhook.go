package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

type webhook struct {
	WebhookHex  string    `json:"webhookHex"`
	Domain      string    `json:"domain"`
	Url         string    `json:"url"`
	Events      string    `json:"events"`
	Secret      string    `json:"secret"`
	CreatedDate time.Time `json:"createdDate"`
}

func webhookList(domain string) ([]webhook, error) {
	if domain == "" {
		return []webhook{}, errorMissingField
	}

	statement := `
		SELECT webhookHex, domain, url, events, secret, createdDate
		FROM webhooks
		WHERE domain = $1;
	`
	rows, err := db.Query(statement, domain)
	if err != nil {
		logger.Errorf("cannot query webhooks: %v", err)
		return nil, errorInternal
	}
	defer rows.Close()

	webhooks := []webhook{}
	for rows.Next() {
		var w webhook
		if err := rows.Scan(&w.WebhookHex, &w.Domain, &w.Url, &w.Events, &w.Secret, &w.CreatedDate); err != nil {
			logger.Errorf("cannot scan webhook: %v", err)
			return nil, errorInternal
		}
		webhooks = append(webhooks, w)
	}

	return webhooks, rows.Err()
}

func webhookNew(domain string, url string, events string) (webhook, error) {
	if domain == "" || url == "" || events == "" {
		return webhook{}, errorMissingField
	}

	webhookHex, err := randomHex(32)
	if err != nil {
		return webhook{}, errorInternal
	}

	secret, err := randomHex(32)
	if err != nil {
		return webhook{}, errorInternal
	}

	statement := `
		INSERT INTO webhooks (webhookHex, domain, url, events, secret, createdDate)
		VALUES ($1, $2, $3, $4, $5, $6);
	`
	_, err = db.Exec(statement, webhookHex, domain, url, events, secret, time.Now().UTC())
	if err != nil {
		logger.Errorf("cannot insert webhook: %v", err)
		return webhook{}, errorInternal
	}

	return webhook{WebhookHex: webhookHex, Domain: domain, Url: url, Events: events, Secret: secret, CreatedDate: time.Now().UTC()}, nil
}

func webhookDelete(webhookHex string) error {
	if webhookHex == "" {
		return errorMissingField
	}

	statement := `DELETE FROM webhooks WHERE webhookHex = $1;`
	_, err := db.Exec(statement, webhookHex)
	if err != nil {
		logger.Errorf("cannot delete webhook: %v", err)
		return errorInternal
	}

	return nil
}

func webhookFire(domain string, event string, payload interface{}) {
	webhooks, err := webhookList(domain)
	if err != nil {
		return
	}

	jsonPayload, err := json.Marshal(map[string]interface{}{
		"event":   event,
		"domain":  domain,
		"payload": payload,
	})
	if err != nil {
		logger.Errorf("cannot marshal webhook payload: %v", err)
		return
	}

	for _, w := range webhooks {
		if !strings.Contains(w.Events, event) {
			continue
		}

		go func(wh webhook) {
			mac := hmac.New(sha256.New, []byte(wh.Secret))
			mac.Write(jsonPayload)
			signature := hex.EncodeToString(mac.Sum(nil))

			req, err := http.NewRequest("POST", wh.Url, bytes.NewBuffer(jsonPayload))
			if err != nil {
				logger.Warningf("cannot create webhook request for %s: %v", wh.Url, err)
				return
			}

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("X-Webhook-Signature", signature)

			client := &http.Client{Timeout: 10 * time.Second}
			resp, err := client.Do(req)
			if err != nil {
				logger.Warningf("webhook delivery failed for %s: %v", wh.Url, err)
				return
			}
			resp.Body.Close()
		}(w)
	}
}

func webhookListHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		OwnerToken *string `json:"ownerToken"`
		Domain     *string `json:"domain"`
	}

	var x request
	if err := bodyUnmarshal(r, &x); err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	o, err := ownerGetByOwnerToken(*x.OwnerToken)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	domain := domainStrip(*x.Domain)
	isOwner, err := domainOwnershipVerify(o.OwnerHex, domain)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}
	if !isOwner {
		bodyMarshal(w, response{"success": false, "message": errorNotAuthorised.Error()})
		return
	}

	webhooks, err := webhookList(domain)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	bodyMarshal(w, response{"success": true, "webhooks": webhooks})
}

func webhookNewHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		OwnerToken *string `json:"ownerToken"`
		Domain     *string `json:"domain"`
		Url        *string `json:"url"`
		Events     *string `json:"events"`
	}

	var x request
	if err := bodyUnmarshal(r, &x); err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	o, err := ownerGetByOwnerToken(*x.OwnerToken)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	domain := domainStrip(*x.Domain)
	isOwner, err := domainOwnershipVerify(o.OwnerHex, domain)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}
	if !isOwner {
		bodyMarshal(w, response{"success": false, "message": errorNotAuthorised.Error()})
		return
	}

	wh, err := webhookNew(domain, *x.Url, *x.Events)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	bodyMarshal(w, response{"success": true, "webhook": wh})
}

func webhookDeleteHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		OwnerToken *string `json:"ownerToken"`
		Domain     *string `json:"domain"`
		WebhookHex *string `json:"webhookHex"`
	}

	var x request
	if err := bodyUnmarshal(r, &x); err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	o, err := ownerGetByOwnerToken(*x.OwnerToken)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	domain := domainStrip(*x.Domain)
	isOwner, err := domainOwnershipVerify(o.OwnerHex, domain)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}
	if !isOwner {
		bodyMarshal(w, response{"success": false, "message": errorNotAuthorised.Error()})
		return
	}

	if err := webhookDelete(*x.WebhookHex); err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	bodyMarshal(w, response{"success": true})
}
