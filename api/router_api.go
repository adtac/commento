package main

import (
	"github.com/gorilla/mux"
)

func apiRouterInit(router *mux.Router) error {
	router.HandleFunc("/api/owner/new", rateLimitMiddleware(10, ownerNewHandler)).Methods("POST")
	router.HandleFunc("/api/owner/confirm-hex", ownerConfirmHexHandler).Methods("GET")
	router.HandleFunc("/api/owner/login", rateLimitMiddleware(20, ownerLoginHandler)).Methods("POST")
	router.HandleFunc("/api/owner/self", ownerSelfHandler).Methods("POST")
	router.HandleFunc("/api/owner/update", ownerUpdateHandler).Methods("POST")
	router.HandleFunc("/api/owner/delete", ownerDeleteHandler).Methods("POST")

	router.HandleFunc("/api/domain/new", domainNewHandler).Methods("POST")
	router.HandleFunc("/api/domain/delete", domainDeleteHandler).Methods("POST")
	router.HandleFunc("/api/domain/clear", domainClearHandler).Methods("POST")
	router.HandleFunc("/api/domain/sso/new", domainSsoSecretNewHandler).Methods("POST")
	router.HandleFunc("/api/domain/list", domainListHandler).Methods("POST")
	router.HandleFunc("/api/domain/update", domainUpdateHandler).Methods("POST")
	router.HandleFunc("/api/domain/moderator/new", domainModeratorNewHandler).Methods("POST")
	router.HandleFunc("/api/domain/moderator/delete", domainModeratorDeleteHandler).Methods("POST")
	router.HandleFunc("/api/domain/statistics", domainStatisticsHandler).Methods("POST")
	router.HandleFunc("/api/domain/import/disqus", domainImportDisqusHandler).Methods("POST")
	router.HandleFunc("/api/domain/import/commento", domainImportCommentoHandler).Methods("POST")
	router.HandleFunc("/api/domain/export/begin", domainExportBeginHandler).Methods("POST")
	router.HandleFunc("/api/domain/export/download", domainExportDownloadHandler).Methods("GET")

	router.HandleFunc("/api/commenter/token/new", commenterTokenNewHandler).Methods("GET")
	router.HandleFunc("/api/commenter/new", rateLimitMiddleware(10, commenterNewHandler)).Methods("POST")
	router.HandleFunc("/api/commenter/login", rateLimitMiddleware(20, commenterLoginHandler)).Methods("POST")
	router.HandleFunc("/api/commenter/self", commenterSelfHandler).Methods("POST")
	router.HandleFunc("/api/commenter/update", commenterUpdateHandler).Methods("POST")
	router.HandleFunc("/api/commenter/photo", commenterPhotoHandler).Methods("GET")

	router.HandleFunc("/api/forgot", forgotHandler).Methods("POST")
	router.HandleFunc("/api/reset", resetHandler).Methods("POST")

	router.HandleFunc("/api/email/get", emailGetHandler).Methods("POST")
	router.HandleFunc("/api/email/update", emailUpdateHandler).Methods("POST")
	router.HandleFunc("/api/email/moderate", emailModerateHandler).Methods("GET")

	router.HandleFunc("/api/oauth/google/redirect", googleRedirectHandler).Methods("GET")
	router.HandleFunc("/api/oauth/google/callback", googleCallbackHandler).Methods("GET")

	router.HandleFunc("/api/oauth/github/redirect", githubRedirectHandler).Methods("GET")
	router.HandleFunc("/api/oauth/github/callback", githubCallbackHandler).Methods("GET")

	router.HandleFunc("/api/oauth/twitter/redirect", twitterRedirectHandler).Methods("GET")
	router.HandleFunc("/api/oauth/twitter/callback", twitterCallbackHandler).Methods("GET")

	router.HandleFunc("/api/oauth/gitlab/redirect", gitlabRedirectHandler).Methods("GET")
	router.HandleFunc("/api/oauth/gitlab/callback", gitlabCallbackHandler).Methods("GET")

	router.HandleFunc("/api/oauth/sso/redirect", ssoRedirectHandler).Methods("GET")
	router.HandleFunc("/api/oauth/sso/callback", ssoCallbackHandler).Methods("GET")

	router.HandleFunc("/api/comment/new", rateLimitMiddleware(30, commentNewHandler)).Methods("POST")
	router.HandleFunc("/api/comment/edit", commentEditHandler).Methods("POST")
	router.HandleFunc("/api/comment/list", commentListHandler).Methods("POST")
	router.HandleFunc("/api/comment/count", commentCountHandler).Methods("POST")
	router.HandleFunc("/api/comment/vote", commentVoteHandler).Methods("POST")
	router.HandleFunc("/api/comment/approve", commentApproveHandler).Methods("POST")
	router.HandleFunc("/api/comment/delete", commentDeleteHandler).Methods("POST")

	router.HandleFunc("/api/page/update", pageUpdateHandler).Methods("POST")

	// Plan and subscription routes
	router.HandleFunc("/api/plan/list", planListHandler).Methods("POST")
	router.HandleFunc("/api/subscription/get", subscriptionGetHandler).Methods("POST")
	router.HandleFunc("/api/subscription/update", subscriptionUpdateHandler).Methods("POST")
	router.HandleFunc("/api/owner/usage", usageGetHandler).Methods("POST")

	// Webhook routes
	router.HandleFunc("/api/domain/webhook/list", webhookListHandler).Methods("POST")
	router.HandleFunc("/api/domain/webhook/new", webhookNewHandler).Methods("POST")
	router.HandleFunc("/api/domain/webhook/delete", webhookDeleteHandler).Methods("POST")

	// Admin routes
	router.HandleFunc("/api/admin/users/list", adminUsersListHandler).Methods("POST")
	router.HandleFunc("/api/admin/users/get", adminUsersGetHandler).Methods("POST")
	router.HandleFunc("/api/admin/domains/list", adminDomainsListHandler).Methods("POST")
	router.HandleFunc("/api/admin/stats", adminStatsHandler).Methods("POST")
	router.HandleFunc("/api/admin/subscription/update", adminSubscriptionUpdateHandler).Methods("POST")

	return nil
}
