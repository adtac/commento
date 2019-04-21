package main

import ()

func domainGet(dmn string) (domain, error) {
	if dmn == "" {
		return domain{}, errorMissingField
	}

	statement := `
		SELECT
			domain,
			ownerHex,
			name,
			creationDate,
			state,
			importedComments,
			autoSpamFilter,
			requireModeration,
			requireIdentification,
			moderateAllAnonymous,
			emailNotificationPolicy,
			commentoProvider,
			googleProvider,
			twitterProvider,
			githubProvider,
			gitlabProvider,
			ssoProvider,
			ssoSecret,
			ssoUrl
		FROM domains
		WHERE domain = $1;
	`
	row := db.QueryRow(statement, dmn)

	var err error
	d := domain{}
	if err = row.Scan(
		&d.Domain,
		&d.OwnerHex,
		&d.Name,
		&d.CreationDate,
		&d.State,
		&d.ImportedComments,
		&d.AutoSpamFilter,
		&d.RequireModeration,
		&d.RequireIdentification,
		&d.ModerateAllAnonymous,
		&d.EmailNotificationPolicy,
		&d.CommentoProvider,
		&d.GoogleProvider,
		&d.TwitterProvider,
		&d.GithubProvider,
		&d.GitlabProvider,
		&d.SsoProvider,
		&d.SsoSecret,
		&d.SsoUrl); err != nil {
		return d, errorNoSuchDomain
	}

	d.Moderators, err = domainModeratorList(d.Domain)
	if err != nil {
		return domain{}, err
	}

	return d, nil
}
