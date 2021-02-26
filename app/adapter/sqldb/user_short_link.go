package sqldb

import (
	"database/sql"
	"fmt"

	"github.com/cross-team/clublink-backend/app/adapter/sqldb/table"
	"github.com/cross-team/clublink-backend/app/entity"
	"github.com/cross-team/clublink-backend/app/usecase/repository"
)

var _ repository.UserShortLink = (*UserShortLinkSQL)(nil)

// UserShortLinkSQL accesses UserShortLink information in user_short_link
// table.
type UserShortLinkSQL struct {
	db *sql.DB
}

// GetUserByShortLink fetches the user associated with a given ShortLink ID
func (u UserShortLinkSQL) GetUserByShortLink(shortLinkID string) (entity.User, error) {
	statement := fmt.Sprintf(`SELECT "%s", "%s", "%s" FROM "%s", "%s" WHERE "%s"=$1 AND "%s"="%s";`,
		table.User.ColumnID,
		table.User.ColumnEmail,
		table.User.ColumnName,
		table.User.TableName,
		table.UserShortLink.TableName,
		table.UserShortLink.ColumnShortLinkID,
		table.UserShortLink.ColumnUserID,
		table.User.ColumnID,
	)

	rows := u.db.QueryRow(statement, shortLinkID)

	user := entity.User{}
	err := rows.Scan(
		&user.ID,
		&user.Email,
		&user.Name,
	)
	if err != nil {
		return entity.User{}, err
	}

	return user, nil
}

// CreateRelation establishes bi-directional relationship between a user and a
// short link in user_short_link table.
func (u UserShortLinkSQL) CreateRelation(user entity.User, shortLinkInput entity.ShortLinkInput) error {
	statement := fmt.Sprintf(`
INSERT INTO "%s" ("%s","%s")
VALUES ($1,$2)
`,
		table.UserShortLink.TableName,
		table.UserShortLink.ColumnUserID,
		table.UserShortLink.ColumnShortLinkID,
	)

	_, err := u.db.Exec(statement, user.ID, shortLinkInput.GetID(""))
	return err
}

// FindAliasesByUser fetches the aliases of all the ShortLinks created by the given user.
// TODO(issue#260): allow API client to filter urls based on visibility.
func (u UserShortLinkSQL) FindAliasesByUser(user entity.User) ([]string, error) {
	statement := fmt.Sprintf(`SELECT "%s" FROM "%s" WHERE "%s"=$1;`,
		table.UserShortLink.ColumnShortLinkID,
		table.UserShortLink.TableName,
		table.UserShortLink.ColumnUserID,
	)

	var aliases []string
	rows, err := u.db.Query(statement, user.ID)
	// TODO(issue#711): errors should be checked before using defer
	defer rows.Close()
	if err != nil {
		return aliases, nil
	}

	for rows.Next() {
		var alias string
		err = rows.Scan(&alias)
		if err != nil {
			return aliases, err
		}

		aliases = append(aliases, alias)
	}

	return aliases, nil
}

// HasMapping checks whether a given short link is tied to a user.
func (u UserShortLinkSQL) HasMapping(user entity.User, alias string) (bool, error) {
	query := fmt.Sprintf(`SELECT "%s" FROM "%s" WHERE "%s"=$1 AND "%s"=$2`,
		table.UserShortLink.ColumnUserID,
		table.UserShortLink.TableName,
		table.UserShortLink.ColumnUserID,
		table.UserShortLink.ColumnShortLinkID,
	)

	var id string
	err := u.db.QueryRow(query, user.ID, alias).Scan(&id)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// NewUserShortLinkSQL creates UserShortLinkSQL
func NewUserShortLinkSQL(db *sql.DB) UserShortLinkSQL {
	return UserShortLinkSQL{
		db: db,
	}
}
