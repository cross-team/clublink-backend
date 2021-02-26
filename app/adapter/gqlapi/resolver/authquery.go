package resolver

import (
	"time"

	"github.com/cross-team/clublink-backend/app/adapter/gqlapi/scalar"
	"github.com/cross-team/clublink-backend/app/entity"
	"github.com/cross-team/clublink-backend/app/usecase/authenticator"
	"github.com/cross-team/clublink-backend/app/usecase/changelog"
	"github.com/cross-team/clublink-backend/app/usecase/repository"
	"github.com/cross-team/clublink-backend/app/usecase/shortlink"
)

// AuthQuery represents GraphQL query resolver that acts differently based
// on the identify of the user
type AuthQuery struct {
	authToken          *string
	authenticator      authenticator.Authenticator
	changeLog          changelog.ChangeLog
	shortLinkRetriever shortlink.Retriever
	userShortLinkRepo  repository.UserShortLink
}

// ShortLinkArgs represents possible parameters for ShortLink endpoint
type ShortLinkArgs struct {
	ID string
}

// ShortLink retrieves a ShortLink given the ID
func (v AuthQuery) ShortLink(args *ShortLinkArgs) (*ShortLink, error) {
	s, err := v.shortLinkRetriever.GetShortLink(args.ID)
	if err != nil {
		return nil, err
	}

	return &ShortLink{shortLink: s}, nil
}

// ActiveShortLinkArgs represents possible parameters for ActiveShortLink endpoint
type ActiveShortLinkArgs struct {
	Alias       string
	ExpireAfter *scalar.Time
}

// ActiveShortLink retrieves an ShortLink persistent storage given alias and expiration time.
func (v AuthQuery) ActiveShortLink(args *ActiveShortLinkArgs) (*ShortLink, error) {
	var expireAt *time.Time
	if args.ExpireAfter != nil {
		expireAt = &args.ExpireAfter.Time
	} else {
		today := time.Now()
		expireAt = &today
	}

	s, err := v.shortLinkRetriever.GetActiveShortLink(args.Alias, expireAt)
	if err != nil {
		return nil, err
	}
	return &ShortLink{shortLink: s}, nil
}

// ShortLinkArgs represents possible parameters for UserByShortLink endpoint
type UserByShortLinkArgs struct {
	ID string
}

// UserByShortLink retrieves a User
func (v AuthQuery) UserByShortLink(args *UserByShortLinkArgs) (*User, error) {
	user, err := v.userShortLinkRepo.GetUserByShortLink(args.ID)
	return &User{user: user}, err
}

// ChangeLog retrieves full ChangeLog from persistent storage
func (v AuthQuery) ChangeLog() (ChangeLog, error) {
	user, err := viewer(v.authToken, v.authenticator)
	if err != nil {
		return newChangeLog([]entity.Change{}, nil), ErrInvalidAuthToken{}
	}

	changeLog, err := v.changeLog.GetChangeLog()
	if err != nil {
		return ChangeLog{}, err
	}

	lastViewedAt, err := v.changeLog.GetLastViewedAt(user)
	return newChangeLog(changeLog, lastViewedAt), err
}

// AllChanges retrieves all the changes that exists in the persistent storage.
func (v AuthQuery) AllChanges() ([]Change, error) {
	user, err := viewer(v.authToken, v.authenticator)
	if err != nil {
		return []Change{}, ErrInvalidAuthToken{}
	}

	changes, err := v.changeLog.GetAllChanges(user)
	if err != nil {
		return []Change{}, err
	}

	var gqlChanges []Change
	for _, change := range changes {
		gqlChanges = append(gqlChanges, newChange(change))
	}

	return gqlChanges, nil
}

// ShortLinks retrieves short links created by a given user from persistent storage
func (v AuthQuery) ShortLinks() ([]ShortLink, error) {
	user, err := viewer(v.authToken, v.authenticator)
	if err != nil {
		return []ShortLink{}, ErrInvalidAuthToken{}
	}

	shortLinks, err := v.shortLinkRetriever.GetShortLinksByUser(user)
	if err != nil {
		return []ShortLink{}, err
	}

	var gqlShortLinks []ShortLink
	for _, v := range shortLinks {
		gqlShortLinks = append(gqlShortLinks, newShortLink(v))
	}

	return gqlShortLinks, nil
}

func newAuthQuery(
	authToken *string,
	authenticator authenticator.Authenticator,
	changeLog changelog.ChangeLog,
	shortLinkRetriever shortlink.Retriever,
	userShortLinkRepo repository.UserShortLink,
) AuthQuery {
	return AuthQuery{
		authToken:          authToken,
		authenticator:      authenticator,
		changeLog:          changeLog,
		shortLinkRetriever: shortLinkRetriever,
		userShortLinkRepo:  userShortLinkRepo,
	}
}
