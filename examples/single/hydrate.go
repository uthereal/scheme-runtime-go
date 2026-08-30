package single

// Hydrate maps table columns dynamically into pointers inside model structs.
var Hydrate = struct {
	User    func(*User, []string) []any
	Profile func(*Profile, []string) []any
	Post    func(*Post, []string) []any
	Comment func(*Comment, []string) []any
	Group   func(*Group, []string) []any
}{
	User: func(model *User, columns []string) []any {
		pointers := make([]any, len(columns))
		for i, col := range columns {
			switch col {
			case Schema.Public.User.ID.ColumnName():
				pointers[i] = &model.ID
			case Schema.Public.User.Email.ColumnName():
				pointers[i] = &model.Email
			case Schema.Public.User.Age.ColumnName():
				pointers[i] = &model.Age
			case Schema.Public.User.Tags.ColumnName():
				pointers[i] = &model.Tags
			case Schema.Public.User.Preferences.ColumnName():
				pointers[i] = &model.Preferences
			case Schema.Public.User.Metadata.ColumnName():
				pointers[i] = &model.Metadata
			case Schema.Public.User.CreatedAt.ColumnName():
				pointers[i] = &model.CreatedAt
			default:
				var discard any
				pointers[i] = &discard
			}
		}
		return pointers
	},
	Profile: func(model *Profile, columns []string) []any {
		pointers := make([]any, len(columns))
		for i, col := range columns {
			switch col {
			case Schema.Public.Profile.ID.ColumnName():
				pointers[i] = &model.ID
			case Schema.Public.Profile.UserID.ColumnName():
				pointers[i] = &model.UserID
			case Schema.Public.Profile.Bio.ColumnName():
				pointers[i] = &model.Bio
			case Schema.Public.Profile.Location.ColumnName():
				pointers[i] = &model.Location
			case Schema.Public.Profile.ActiveDuration.ColumnName():
				pointers[i] = &model.ActiveDuration
			case Schema.Public.Profile.IsPublic.ColumnName():
				pointers[i] = &model.IsPublic
			default:
				var discard any
				pointers[i] = &discard
			}
		}
		return pointers
	},
	Post: func(model *Post, columns []string) []any {
		pointers := make([]any, len(columns))
		for i, col := range columns {
			switch col {
			case Schema.Public.Post.ID.ColumnName():
				pointers[i] = &model.ID
			case Schema.Public.Post.UserID.ColumnName():
				pointers[i] = &model.UserID
			case Schema.Public.Post.Title.ColumnName():
				pointers[i] = &model.Title
			case Schema.Public.Post.Content.ColumnName():
				pointers[i] = &model.Content
			case Schema.Public.Post.Rating.ColumnName():
				pointers[i] = &model.Rating
			default:
				var discard any
				pointers[i] = &discard
			}
		}
		return pointers
	},
	Comment: func(model *Comment, columns []string) []any {
		pointers := make([]any, len(columns))
		for i, col := range columns {
			switch col {
			case Schema.Public.Comment.ID.ColumnName():
				pointers[i] = &model.ID
			case Schema.Public.Comment.PostID.ColumnName():
				pointers[i] = &model.PostID
			case Schema.Public.Comment.UserID.ColumnName():
				pointers[i] = &model.UserID
			case Schema.Public.Comment.Text.ColumnName():
				pointers[i] = &model.Text
			default:
				var discard any
				pointers[i] = &discard
			}
		}
		return pointers
	},
	Group: func(model *Group, columns []string) []any {
		pointers := make([]any, len(columns))
		for i, col := range columns {
			switch col {
			case Schema.Public.Group.ID.ColumnName():
				pointers[i] = &model.ID
			case Schema.Public.Group.Name.ColumnName():
				pointers[i] = &model.Name
			default:
				var discard any
				pointers[i] = &discard
			}
		}
		return pointers
	},
}
