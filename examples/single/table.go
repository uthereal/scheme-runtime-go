package single

import (
	"github.com/uthereal/scheme-runtime-go/pkg/contract"
)

// Table defines the table metadata.
var Table = struct {
	Public struct {
		User    contract.TableMetadata[User]
		Profile contract.TableMetadata[Profile]
		Post    contract.TableMetadata[Post]
		Comment contract.TableMetadata[Comment]
		Group   contract.TableMetadata[Group]
	}
}{
	Public: struct {
		User    contract.TableMetadata[User]
		Profile contract.TableMetadata[Profile]
		Post    contract.TableMetadata[Post]
		Comment contract.TableMetadata[Comment]
		Group   contract.TableMetadata[Group]
	}{
		User: contract.TableMetadata[User]{
			SchemaName: "public",
			TableName:  "users",
			DefaultColumns: []contract.Column[User]{
				Schema.Public.User.ID,
				Schema.Public.User.Email,
				Schema.Public.User.Age,
				Schema.Public.User.Tags,
				Schema.Public.User.Preferences,
				Schema.Public.User.Metadata,
				Schema.Public.User.CreatedAt,
			},
		},
		Profile: contract.TableMetadata[Profile]{
			SchemaName: "public",
			TableName:  "profiles",
			DefaultColumns: []contract.Column[Profile]{
				Schema.Public.Profile.ID,
				Schema.Public.Profile.UserID,
				Schema.Public.Profile.Bio,
				Schema.Public.Profile.Location,
				Schema.Public.Profile.ActiveDuration,
				Schema.Public.Profile.IsPublic,
			},
		},
		Post: contract.TableMetadata[Post]{
			SchemaName: "public",
			TableName:  "posts",
			DefaultColumns: []contract.Column[Post]{
				Schema.Public.Post.ID,
				Schema.Public.Post.UserID,
				Schema.Public.Post.Title,
				Schema.Public.Post.Content,
				Schema.Public.Post.Rating,
			},
		},
		Comment: contract.TableMetadata[Comment]{
			SchemaName: "public",
			TableName:  "comments",
			DefaultColumns: []contract.Column[Comment]{
				Schema.Public.Comment.ID,
				Schema.Public.Comment.PostID,
				Schema.Public.Comment.UserID,
				Schema.Public.Comment.Text,
			},
		},
		Group: contract.TableMetadata[Group]{
			SchemaName: "public",
			TableName:  "groups",
			DefaultColumns: []contract.Column[Group]{
				Schema.Public.Group.ID,
				Schema.Public.Group.Name,
			},
		},
	},
}
