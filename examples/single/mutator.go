package single

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/uthereal/scheme-runtime-go/pkg/contract"
)

// UserMutator manages partial payload mutations for User fields.
type UserMutator struct {
	ID          contract.Set[int64]
	Email       contract.Set[string]
	Age         contract.Set[*int]
	Tags        contract.Set[[]string]
	Preferences contract.Set[*map[string]any]
	Metadata    contract.Set[map[string]string]
	CreatedAt   contract.Set[time.Time]
}

// ProfileMutator manages partial payload mutations for Profile fields.
type ProfileMutator struct {
	ID             contract.Set[int64]
	UserID         contract.Set[int64]
	Bio            contract.Set[*string]
	Location       contract.Set[*pgtype.Point]
	ActiveDuration contract.Set[time.Duration]
	IsPublic       contract.Set[*bool]
}

// PostMutator manages partial payload mutations for Post fields.
type PostMutator struct {
	ID      contract.Set[int64]
	UserID  contract.Set[int64]
	Title   contract.Set[string]
	Content contract.Set[string]
	Rating  contract.Set[pgtype.Numeric]
}

// CommentMutator manages partial payload mutations for Comment fields.
type CommentMutator struct {
	ID     contract.Set[int64]
	PostID contract.Set[int64]
	UserID contract.Set[int64]
	Text   contract.Set[string]
}

// GroupMutator manages partial payload mutations for Group fields.
type GroupMutator struct {
	ID   contract.Set[string]
	Name contract.Set[string]
}
