package handle

import (
	"net/http"

	"github.com/cross-team/clublink-backend/app/adapter/request"
	"github.com/short-d/app/fw/router"
)

// Track records event happened in the API caller.
func Track(instrumentationFactory request.InstrumentationFactory) router.Handle {
	return func(w http.ResponseWriter, r *http.Request, params router.Params) {
		i := instrumentationFactory.NewHTTP(r)

		event := params["event"]
		i.Track(event)

		w.WriteHeader(http.StatusNoContent)
	}
}
