package handlers

import (
	"net/http"

	"github.com/calvindc/Web3RpcHub/db"
	"github.com/calvindc/Web3RpcHub/hubstat"
	"github.com/calvindc/Web3RpcHub/internal/network"
	"github.com/calvindc/Web3RpcHub/internal/repository"
	"github.com/calvindc/Web3RpcHub/internal/signalbridge"
	"github.com/calvindc/Web3RpcHub/models/web/i18n"
	"go.mindeco.de/logging"
)

var HTMLTemplates = []string{
	"landing/index.tmpl",
	"alias.tmpl",

	"change-member-password.tmpl",

	"invite/consumed.tmpl",
	"invite/facade.tmpl",
	"invite/facade-fallback.tmpl",
	"invite/insert-id.tmpl",

	"notice/list.tmpl",
	"notice/show.tmpl",

	"error.tmpl",
}

// Databases is an options stuct for the required databases of the web handlers
// db.interface
type Databases struct {
	Aliases       db.AliasesService
	AuthFallback  db.AuthFallbackService
	AuthWithToken db.AuthWithTokenService
	Config        db.HubConfig
	DeniedKeys    db.DeniedKeysService
	Invites       db.InvitesService
	Notices       db.NoticesService
	Members       db.MembersService
	PinnedNotices db.PinnedNoticesService
}

// New initializes the whole web stack for hubs, with all the sub-modules and routing.
/*
func StartHubServ(hMembers db.MembersService, hDeniedKeys db.DeniedKeysService, hAlias db.AliasesService,
	hAuthWithToken db.AuthWithTokenService, hAuthWithBirdge *signalbridge.SignalBridge,
	hConfig db.HubConfig, hNetInfo network.HubEndpoint, opts ...Option) (*HubServe, error) {
*/
func NewWebHandler(
	logger logging.Interface,
	repo repository.Interface,
	netInfo network.HubEndpoint,
	hubState *hubstat.HubNetManager,
	hubEndpoints network.Endpoints,
	bridge *signalbridge.SignalBridge,
	dbs Databases,
) (http.Handler, error) {

	_, err := i18n.New(repo, dbs.Config)
	if err != nil {
		return nil, err
	}

	//cookieCodec := locHelper.CookieCodec
	//cookieStore := locHelper.CookieStore

	//flashHelper := errors.NewFlashHelper(cookieStore, locHelper)

	//errh := errors.NewErrorHandler(locHelper, flashHelper)

	return nil, nil
}

// utils
func concatTemplates(lst ...[]string) []string {
	var catted []string

	for _, tpls := range lst {
		for _, t := range tpls {
			catted = append(catted, t)
		}

	}
	return catted
}
