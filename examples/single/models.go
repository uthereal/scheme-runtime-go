package single

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

// User represents a central user entity.
type User struct {
	ID          int64
	Email       string
	Age         *int
	Tags        []string
	Preferences *map[string]any
	Metadata    map[string]string
	CreatedAt   time.Time

	Profile *Profile
	Posts   []*Post
	Groups  []*Group
}

// Profile represents an isolated user profile.
type Profile struct {
	ID             int64
	UserID         int64
	Bio            *string
	Location       *pgtype.Point
	ActiveDuration time.Duration
	IsPublic       *bool
}

// Post represents a blog post or discussion thread.
type Post struct {
	ID       int64
	UserID   int64
	Title    string
	Content  string
	Rating   pgtype.Numeric
	Comments []*Comment
}

// Comment represents a reply to a specific Post.
type Comment struct {
	ID     int64
	PostID int64
	UserID int64
	Text   string
	Post   *Post
	User   *User
}

// Group represents a logical collection of users.
type Group struct {
	ID   string
	Name string
}
