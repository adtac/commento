package main

import ()

var domainsRowColumns = `
	domains.domain,
	domains.ownerHex,
	domains.name,
	domains.creationDate,
	domains.state,
	domains.importedComments,
	domains.autoSpamFilter,
	domains.requireModeration,
	domains.requireIdentification,
	domains.moderateAllAnonymous,
	domains.emailNotificationPolicy,
	domains.commentoProvider,
	domains.googleProvider,
	domains.twitterProvider,
	domains.githubProvider,
	domains.gitlabProvider,
	domains.ssoProvider,
	domains.ssoSecret,
	domains.ssoUrl,
	domains.defaultSortPolicy
`

func domainsRowScan(s sqlScanner, d *domain) error {
	return s.Scan(
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
		&d.SsoUrl,
		&d.DefaultSortPolicy,
	)
}

func domainGet(dmn string) (domain, error) {
	if dmn == "" {
		return domain{}, errorMissingField
	}

	statement := `
		SELECT ` + domainsRowColumns + `
		FROM domains
		WHERE domain = $1;
	`
	row := db.QueryRow(statement, dmn)

	var err error
	d := domain{}
	if err = domainsRowScan(row, &d); err != nil {
		return d, errorNoSuchDomain
	}

	d.Moderators, err = domainModeratorList(d.Domain)
	if err != nil {
		return domain{}, err
	}

	return d, nil
}
