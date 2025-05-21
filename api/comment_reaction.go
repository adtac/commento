package main

import (
	"net/http"
	"time"
)

func commentReactionHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		CommenterToken *string `json:"commenterToken"`
		CommentHex     *string `json:"commentHex"`
		ReactionType   *string `json:"reactionType"`
		Remove         *bool   `json:"remove"`
	}

	var x request
	if err := bodyUnmarshal(r, &x); err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	c, err := commenterGetByCommenterToken(*x.CommenterToken)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	domain, _, err := commentDomainPathGet(*x.CommentHex)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	isFrozen, err := domainIsFrozen(domain)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}
	if isFrozen {
		bodyMarshal(w, response{"success": false, "message": errorDomainFrozen.Error()})
		return
	}

	if *x.Remove {
		// Remove the reaction
		if err := commentReactionRemove(*x.CommentHex, c.CommenterHex, *x.ReactionType); err != nil {
			bodyMarshal(w, response{"success": false, "message": err.Error()})
			return
		}
	} else {
		// Add the reaction
		if err := commentReactionAdd(*x.CommentHex, c.CommenterHex, *x.ReactionType); err != nil {
			bodyMarshal(w, response{"success": false, "message": err.Error()})
			return
		}
	}

	bodyMarshal(w, response{"success": true})
}

func commentReactionAdd(commentHex string, commenterHex string, reactionType string) error {
	statement := `
		INSERT INTO
		reactions (commentHex, commenterHex, type, reactionDate)
		VALUES   ($1,         $2,           $3,   $4)
		ON CONFLICT (commentHex, commenterHex, type) DO NOTHING;
	`
	if _, err := db.Exec(statement, commentHex, commenterHex, reactionType, time.Now().UTC()); err != nil {
		logger.Errorf("cannot insert reaction: %v", err)
		return errorInternal
	}

	return nil
}

func commentReactionRemove(commentHex string, commenterHex string, reactionType string) error {
	statement := `
		DELETE FROM reactions
		WHERE commentHex = $1 AND commenterHex = $2 AND type = $3;
	`
	if _, err := db.Exec(statement, commentHex, commenterHex, reactionType); err != nil {
		logger.Errorf("cannot delete reaction: %v", err)
		return errorInternal
	}

	return nil
}

// commentReactionsGet returns a map of reaction types to their counts for a given comment
func commentReactionsGet(commentHex string) (map[string]int, error) {
	statement := `
		SELECT type, COUNT(*) 
		FROM reactions 
		WHERE commentHex = $1 
		GROUP BY type;
	`
	rows, err := db.Query(statement, commentHex)
	if err != nil {
		logger.Errorf("error selecting reactions: %v", err)
		return nil, errorInternal
	}
	defer rows.Close()

	reactions := make(map[string]int)
	for rows.Next() {
		var reactionType string
		var count int

		if err = rows.Scan(&reactionType, &count); err != nil {
			logger.Errorf("error scanning reaction row: %v", err)
			return nil, errorInternal
		}

		reactions[reactionType] = count
	}

	return reactions, nil
}

// domainIsFrozen checks if a domain is frozen (read-only mode)
func domainIsFrozen(domain string) (bool, error) {
	d, err := domainGet(domain)
	if err != nil {
		return false, err
	}

	return d.State == "frozen", nil
}

func commentSetSpoiler(commentHex string, spoiler bool) error {
	statement := `
		UPDATE comments
		SET spoiler = $2
		WHERE commentHex = $1;
	`
	if _, err := db.Exec(statement, commentHex, spoiler); err != nil {
		logger.Errorf("cannot update comment spoiler status: %v", err)
		return errorInternal
	}

	return nil
}

// A struct to represent a spoiler request
type commentSpoilerRequest struct {
	CommenterToken string `json:"commenterToken"`
	CommentHex     string `json:"commentHex"`
	Spoiler        bool   `json:"spoiler"`
}

func commentSpoilerHandler(w http.ResponseWriter, r *http.Request) {
	var req commentSpoilerRequest

	if err := bodyUnmarshal(r, &req); err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	c, err := commenterGetByCommenterToken(req.CommenterToken)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}
	domain, _, err := commentDomainPathGet(req.CommentHex)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	// Check if commenter owns the comment or is a moderator
	isOwner, err := commentOwnershipVerify(req.CommentHex, c.CommenterHex)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}
	isModerator, err := isDomainModerator(domain, c.Email)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}
	if !isOwner && !isModerator {
		bodyMarshal(w, response{"success": false, "message": errorNotAuthorised.Error()})
		return
	}

	// Check if the domain is frozen
	isFrozen, err := domainIsFrozen(domain)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}
	if isFrozen {
		bodyMarshal(w, response{"success": false, "message": errorDomainFrozen.Error()})
		return
	}
	if err := commentSetSpoiler(req.CommentHex, req.Spoiler); err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	bodyMarshal(w, response{"success": true})
}
