package errors

import (
	"errors"
	"html/template"
	"net/http"

	"github.com/calvindc/Web3RpcHub/db"
	"github.com/calvindc/Web3RpcHub/models/web/i18n"
	"go.mindeco.de/http/auth"
	"go.mindeco.de/http/render"

	"go.mindeco.de/log/level"
	"go.mindeco.de/logging"
)

type ErrorHandler struct {
	locHelper *i18n.LangHelper
	render    *render.Renderer
	flashes   *FlashHelper
}

func NewErrorHandler(locHelper *i18n.LangHelper, flashes *FlashHelper) *ErrorHandler {
	return &ErrorHandler{
		locHelper: locHelper,
		flashes:   flashes,
	}
}

// SetRenderer needs to update the renderer later since we need to pass ErrorHandler into render.New (i.e. before we get the pointer for r)
func (eh *ErrorHandler) SetRenderer(r *render.Renderer) {
	eh.render = r
}

func (eh *ErrorHandler) Handle(rw http.ResponseWriter, req *http.Request, code int, err error) {
	log := logging.FromContext(req.Context())
	level.Error(log).Log("event", "handling error", "path", req.URL.Path, "err", err)
	var redirectErr ErrRedirect
	if errors.As(err, &redirectErr) {
		if redirectErr.Reason != nil {
			eh.flashes.AddError(rw, req, redirectErr.Reason)
		}
		http.Redirect(rw, req, redirectErr.Path, http.StatusSeeOther)
		return
	}

	var ih = eh.locHelper.FromRequest(req)

	code, msg := localizeError(ih, err)

	data := errorTemplateData{
		Err: template.HTML(msg),
		// TODO: localize status codes? might be fine with a few
		Status:     http.StatusText(code),
		StatusCode: code,

		BackURL: req.URL.Path,
	}

	if code == http.StatusNotFound {
		data.BackURL = "/"
		referer := req.Header.Get("Referer")
		if referer != "" {
			data.BackURL = referer
		}
	}

	renderErr := eh.render.Render(rw, req, "error.tmpl", code, data)
	if renderErr != nil {
		logger := logging.FromContext(req.Context())
		level.Error(logger).Log("event", "error template renderfailed",
			"orig-err", err,
			"render-err", renderErr,
		)
	}
}

type errorTemplateData struct {
	StatusCode int
	Status     string
	Err        template.HTML

	BackURL string
}

func localizeError(ih *i18n.Localizer, err error) (int, template.HTML) {

	// default, unlocalized message
	msg := template.HTML(err.Error())

	// localize some specific error messages
	var (
		aa  db.ErrAlreadyAdded
		pnf PageNotFound
		br  ErrBadRequest
		f   ErrForbidden
		gl  ErrGenericLocalized
	)

	code := http.StatusInternalServerError

	switch {

	case err == ErrNotAuthorized:
		code = http.StatusForbidden
		msg = ih.LocalizeSimple("ErrorNotAuthorized")

	case err == auth.ErrBadLogin:
		msg = ih.LocalizeSimple("ErrorAuthBadLogin")

	case errors.As(err, &gl):
		msg = ih.LocalizeSimple(gl.Label)

	case errors.Is(err, db.ErrNotFound):
		code = http.StatusNotFound
		msg = ih.LocalizeSimple("ErrorNotFound")

	case errors.As(err, &aa):
		msg = ih.LocalizeWithData("ErrorAlreadyAdded", "Feed", aa.Ref.String())

	case errors.As(err, &pnf):
		code = http.StatusNotFound
		msg = ih.LocalizeWithData("ErrorPageNotFound", "Path", pnf.Path)

	case errors.As(err, &br):
		code = http.StatusBadRequest
		// TODO: we could localize all the "Where:" as labels, too
		// buttt it feels like overkill right now
		msg = ih.LocalizeWithData("ErrorBadRequest", "Where", br.Where, "Details", br.Details.Error())

	case errors.As(err, &f):
		code = http.StatusForbidden
		msg = ih.LocalizeWithData("ErrorForbidden", "Details", f.Details.Error())
	}

	return code, msg
}
