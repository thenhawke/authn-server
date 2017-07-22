package sessions

import (
	"net/http"

	"github.com/keratin/authn-server/api"
	"github.com/keratin/authn-server/models"
)

func getSessionRefresh(app *api.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// check for valid session with live token
		accountId := api.GetSessionAccountId(r)
		if accountId == 0 {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// refresh the refresh token
		session := api.GetSession(r)
		err := app.RefreshTokenStore.Touch(models.RefreshToken(session.Subject), accountId)
		if err != nil {
			panic(err)
		}

		// generate the requested identity token
		identityToken, err := api.IdentityForSession(app.KeyStore, app.Config, session, accountId)
		if err != nil {
			panic(err)
		}

		api.WriteData(w, http.StatusCreated, map[string]string{
			"id_token": identityToken,
		})
	}
}
