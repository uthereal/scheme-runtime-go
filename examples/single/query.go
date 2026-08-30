package single

import (
	"github.com/uthereal/scheme-runtime-go/pkg/contract"
	"github.com/uthereal/scheme-runtime-go/pkg/grammar"
	"github.com/uthereal/scheme-runtime-go/pkg/orm"
)

// NewUserQuery constructs a QueryBuilder for User.
func NewUserQuery(
	db contract.DB,
) *orm.QueryBuilder[User, UserMutator] {
	compiler := grammar.NewPostgresGrammar()
	res := orm.NewQueryBuilder[User, UserMutator](
		db,
		compiler,
		Table.Public.User,
		Hydrate.User,
		Dehydrate.User,
	)
	return res
}

// NewProfileQuery constructs a QueryBuilder for Profile.
func NewProfileQuery(
	db contract.DB,
) *orm.QueryBuilder[Profile, ProfileMutator] {
	compiler := grammar.NewPostgresGrammar()
	res := orm.NewQueryBuilder[Profile, ProfileMutator](
		db,
		compiler,
		Table.Public.Profile,
		Hydrate.Profile,
		Dehydrate.Profile,
	)
	return res
}

// NewPostQuery constructs a QueryBuilder for Post.
func NewPostQuery(
	db contract.DB,
) *orm.QueryBuilder[Post, PostMutator] {
	compiler := grammar.NewPostgresGrammar()
	res := orm.NewQueryBuilder[Post, PostMutator](
		db,
		compiler,
		Table.Public.Post,
		Hydrate.Post,
		Dehydrate.Post,
	)
	return res
}

// NewCommentQuery constructs a QueryBuilder for Comment.
func NewCommentQuery(
	db contract.DB,
) *orm.QueryBuilder[Comment, CommentMutator] {
	compiler := grammar.NewPostgresGrammar()
	res := orm.NewQueryBuilder[Comment, CommentMutator](
		db,
		compiler,
		Table.Public.Comment,
		Hydrate.Comment,
		Dehydrate.Comment,
	)
	return res
}

// NewGroupQuery constructs a QueryBuilder for Group.
func NewGroupQuery(
	db contract.DB,
) *orm.QueryBuilder[Group, GroupMutator] {
	compiler := grammar.NewPostgresGrammar()
	res := orm.NewQueryBuilder[Group, GroupMutator](
		db,
		compiler,
		Table.Public.Group,
		Hydrate.Group,
		Dehydrate.Group,
	)
	return res
}
