package provider

import (
	"time"

	"github.com/cross-team/clublink-backend/app/usecase/repository"
	"github.com/cross-team/clublink-backend/app/usecase/search"
	"github.com/short-d/app/fw/logger"
)

// SearchTimeout represents timeout duration of a search request.
type SearchTimeout time.Duration

// NewSearch creates Search given its dependencies.
func NewSearch(
	logger logger.Logger,
	shortLinkRepo repository.ShortLink,
	userShortLinkRepo repository.UserShortLink,
	timeout SearchTimeout,
) search.Search {
	return search.NewSearch(logger, shortLinkRepo, userShortLinkRepo, time.Duration(timeout))
}
