package main

import ()

func emailNotificationModerator(d domain, path string, title string, commenterHex string, commentHex string, html string, state string) {
	if d.EmailNotificationPolicy == "none" {
		return
	}

	if d.EmailNotificationPolicy == "pending-moderation" && state == "approved" {
		return
	}

	var commenterName string
	var commenterEmail string
	if commenterHex == "anonymous" {
		commenterName = "Anonymous"
	} else {
		c, err := commenterGetByHex(commenterHex)
		if err != nil {
			logger.Errorf("cannot get commenter to send email notification: %v", err)
			return
		}

		commenterName = c.Name
		commenterEmail = c.Email
	}

	kind := d.EmailNotificationPolicy
	if state != "approved" {
		kind = "pending-moderation"
	}

	for _, m := range d.Moderators {
		// Do not email the commenting moderator their own comment.
		if commenterHex != "anonymous" && m.Email == commenterEmail {
			continue
		}

		e, err := emailGet(m.Email)
		if err != nil {
			// No such email.
			continue
		}

		if !e.SendModeratorNotifications {
			continue
		}

		statement := `
			SELECT name
			FROM commenters
			WHERE email = $1;
		`
		row := db.QueryRow(statement, m.Email)
		var name string
		if err := row.Scan(&name); err != nil {
			// The moderator has probably not created a commenter account.
			// We should only send emails to people who signed up, so skip.
			continue
		}

		if err := smtpEmailNotification(m.Email, name, kind, d.Domain, path, commentHex, commenterName, title, html, e.UnsubscribeSecretHex); err != nil {
			logger.Errorf("error sending email to %s: %v", m.Email, err)
			continue
		}
	}
}

func emailNotificationReply(d domain, path string, title string, commenterHex string, commentHex string, html string, parentHex string, state string) {
	// No reply notifications for root comments.
	if parentHex == "root" {
		return
	}

	// No reply notification emails for unapproved comments.
	if state != "approved" {
		return
	}

	statement := `
		SELECT commenterHex
		FROM comments
		WHERE commentHex = $1;
	`
	row := db.QueryRow(statement, parentHex)

	var parentCommenterHex string
	err := row.Scan(&parentCommenterHex)
	if err != nil {
		logger.Errorf("cannot scan commenterHex and parentCommenterHex: %v", err)
		return
	}

	// No reply notification emails for anonymous users.
	if parentCommenterHex == "anonymous" {
		return
	}

	// No reply notification email for self replies.
	if parentCommenterHex == commenterHex {
		return
	}

	pc, err := commenterGetByHex(parentCommenterHex)
	if err != nil {
		logger.Errorf("cannot get commenter to send email notification: %v", err)
		return
	}

	var commenterName string
	if commenterHex == "anonymous" {
		commenterName = "Anonymous"
	} else {
		c, err := commenterGetByHex(commenterHex)
		if err != nil {
			logger.Errorf("cannot get commenter to send email notification: %v", err)
			return
		}
		commenterName = c.Name
	}

	epc, err := emailGet(pc.Email)
	if err != nil {
		// No such email.
		return
	}

	if !epc.SendReplyNotifications {
		return
	}

	smtpEmailNotification(pc.Email, pc.Name, "reply", d.Domain, path, commentHex, commenterName, title, html, epc.UnsubscribeSecretHex)
}

func emailNotificationNew(d domain, path string, commenterHex string, commentHex string, html string, parentHex string, state string) {
	p, err := pageGet(d.Domain, path)
	if err != nil {
		logger.Errorf("cannot get page to send email notification: %v", err)
		return
	}

	if p.Title == "" {
		p.Title, err = pageTitleUpdate(d.Domain, path)
		if err != nil {
			// Not being able to update a page title isn't serious enough to skip an
			// email notification.
			p.Title = d.Domain
		}
	}

	emailNotificationModerator(d, path, p.Title, commenterHex, commentHex, html, state)
	emailNotificationReply(d, path, p.Title, commenterHex, commentHex, html, parentHex, state)
}
