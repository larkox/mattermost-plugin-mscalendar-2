// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package http

import (
	"net/http"

	"github.com/mattermost/mattermost-plugin-mscalendar/server/api"
)

func (h *Handler) oauth2Connect(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("Mattermost-User-ID")
	if userID == "" {
		http.Error(w, "Not authorized", http.StatusUnauthorized)
		return
	}

	api := api.FromContext(r.Context())

	var url string
	var err error
	if r.URL.Query().Get("bot") == "true" {
		isAdmin, err := api.IsAuthorizedAdmin(userID)
		if err != nil || !isAdmin {
			http.Error(w, "Not authorized", http.StatusUnauthorized)
			return
		}
		url, err = api.InitOAuth2(h.Config.BotUserID)
	} else {
		url, err = api.InitOAuth2(userID)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, url, http.StatusFound)
}
