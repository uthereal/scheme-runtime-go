package single

import (
	"github.com/uthereal/scheme-runtime-go/pkg/contract"
	"github.com/uthereal/scheme-runtime-go/pkg/orm"
	"github.com/uthereal/scheme-runtime-go/pkg/orm/relation"
)

func init() {
	Schema.Public.User.Profile = relation.HasOne[
		User,
		Profile,
		ProfileMutator,
	]{
		Relation: relation.Relation[User, Profile, ProfileMutator]{
			ForeignKeyColumns: []contract.Column[Profile]{
				Schema.Public.Profile.UserID,
			},
			ChildQueryFactory: func(
				db contract.DB,
			) *orm.QueryBuilder[Profile, ProfileMutator] {
				return NewProfileQuery(db)
			},
			LocalKeyExtractor: func(p *User) any {
				return p.ID
			},
			ForeignKeyExtractor: func(c *Profile) any {
				return c.UserID
			},
			Hydrator: func(p *User, children []Profile) {
				if len(children) > 0 {
					p.Profile = &children[0]
				}
			},
		},
	}

	Schema.Public.User.Posts = relation.HasMany[
		User,
		Post,
		PostMutator,
	]{
		Relation: relation.Relation[User, Post, PostMutator]{
			ForeignKeyColumns: []contract.Column[Post]{
				Schema.Public.Post.UserID,
			},
			ChildQueryFactory: func(
				db contract.DB,
			) *orm.QueryBuilder[Post, PostMutator] {
				return NewPostQuery(db)
			},
			LocalKeyExtractor: func(p *User) any {
				return p.ID
			},
			ForeignKeyExtractor: func(c *Post) any {
				return c.UserID
			},
			Hydrator: func(p *User, children []Post) {
				ptrs := make([]*Post, len(children))
				for i := range children {
					ptrs[i] = &children[i]
				}
				p.Posts = ptrs
			},
		},
	}

	Schema.Public.User.Groups = relation.BelongsToMany[
		User,
		Group,
		GroupMutator,
	]{
		PivotTable:              "user_groups",
		PivotForeignKeyToParent: "user_id",
		PivotForeignKeyToChild:  "group_id",
		ChildForeignKeyColumn:   Schema.Public.Group.ID,
		ChildQueryFactory: func(
			db contract.DB,
		) *orm.QueryBuilder[Group, GroupMutator] {
			return NewGroupQuery(db)
		},
		LocalKeyExtractor: func(p *User) any {
			return p.ID
		},
		ChildKeyExtractor: func(c *Group) any {
			return c.ID
		},
		Hydrator: func(p *User, children []Group) {
			ptrs := make([]*Group, len(children))
			for i := range children {
				ptrs[i] = &children[i]
			}
			p.Groups = ptrs
		},
	}

	Schema.Public.Post.Comments = relation.HasMany[
		Post,
		Comment,
		CommentMutator,
	]{
		Relation: relation.Relation[Post, Comment, CommentMutator]{
			ForeignKeyColumns: []contract.Column[Comment]{
				Schema.Public.Comment.PostID,
			},
			ChildQueryFactory: func(
				db contract.DB,
			) *orm.QueryBuilder[Comment, CommentMutator] {
				return NewCommentQuery(db)
			},
			LocalKeyExtractor: func(p *Post) any {
				return p.ID
			},
			ForeignKeyExtractor: func(c *Comment) any {
				return c.PostID
			},
			Hydrator: func(p *Post, children []Comment) {
				ptrs := make([]*Comment, len(children))
				for i := range children {
					ptrs[i] = &children[i]
				}
				p.Comments = ptrs
			},
		},
	}

	Schema.Public.Comment.Post = relation.BelongsTo[
		Comment,
		Post,
		PostMutator,
	]{
		Relation: relation.Relation[Comment, Post, PostMutator]{
			ForeignKeyColumns: []contract.Column[Post]{
				Schema.Public.Post.ID,
			},
			ChildQueryFactory: func(
				db contract.DB,
			) *orm.QueryBuilder[Post, PostMutator] {
				return NewPostQuery(db)
			},
			LocalKeyExtractor: func(c *Comment) any {
				return c.PostID
			},
			ForeignKeyExtractor: func(p *Post) any {
				return p.ID
			},
			Hydrator: func(c *Comment, children []Post) {
				if len(children) > 0 {
					c.Post = &children[0]
				}
			},
		},
	}

	Schema.Public.Comment.User = relation.BelongsTo[
		Comment,
		User,
		UserMutator,
	]{
		Relation: relation.Relation[Comment, User, UserMutator]{
			ForeignKeyColumns: []contract.Column[User]{
				Schema.Public.User.ID,
			},
			ChildQueryFactory: func(
				db contract.DB,
			) *orm.QueryBuilder[User, UserMutator] {
				return NewUserQuery(db)
			},
			LocalKeyExtractor: func(c *Comment) any {
				return c.UserID
			},
			ForeignKeyExtractor: func(u *User) any {
				return u.ID
			},
			Hydrator: func(c *Comment, children []User) {
				if len(children) > 0 {
					c.User = &children[0]
				}
			},
		},
	}
}
