package action

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/lbryio/notifica/app/config"
	"github.com/lbryio/notifica/app/metrics"

	"github.com/lbryio/lbry.go/v2/extras/api"
	"github.com/lbryio/lbry.go/v2/extras/errors"
	"github.com/lbryio/lbry.go/v2/extras/util"
	v "github.com/lbryio/ozzo-validation"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
)

func init() {
	// make validation always return json-style names for fields
	f := func(str string) string {
		return util.Underscore(str)
	}
	v.ErrorTagFunc = &f
}

func configureAPIServer() {
	api.TraceEnabled = config.IsDebugMode

	api.IgnoredFormFields = []string{}

	hs := make(map[string]string)

	hs["Server"] = "lbry.com"
	hs["Content-Type"] = "application/json; charset=utf-8"

	hs["Access-Control-Allow-Methods"] = "GET, PUT, POST, DELETE, OPTIONS"
	hs["Access-Control-Allow-Origin"] = "*"

	hs["X-Content-Type-Options"] = "nosniff"
	hs["X-Frame-Options"] = "deny"
	hs["Content-Security-Policy"] = "default-src 'none'"
	hs["X-XSS-Protection"] = "1; mode=block"
	hs["Referrer-Policy"] = "same-origin"
	if !config.IsDebugMode {
		hs["Strict-Transport-Security"] = "max-age=31536000; preload"
	}

	api.ResponseHeaders = hs

	api.Log = func(request *http.Request, response *api.Response, err error) {
		consoleText := request.RemoteAddr + " [" + strconv.Itoa(response.Status) + "]: " + request.Method + " " + request.URL.Path
		path := strings.TrimLeft(request.URL.Path, "/")
		metrics.StatusErrors.WithLabelValues(path, strconv.Itoa(response.Status)).Inc()
		if err == nil {
			logrus.Debug(color.GreenString(consoleText))
		} else {
			logrus.Debug(color.RedString(consoleText + ": " + err.Error()))
			if response.Status >= http.StatusInternalServerError {
				metrics.ServerErrors.WithLabelValues(path).Inc()
				if config.IsDebugMode {
					err := util.SendToSlack(strconv.Itoa(response.Status) + " " + request.Method + " " + request.URL.Path + ": " + errors.FullTrace(response.Error))
					if err != nil {
						logrus.Error(err)
					}
				}
			}
		}
	}
}
