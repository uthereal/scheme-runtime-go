package single

// Dehydrate maps mutator payloads to column slices and parameter slices.
var Dehydrate = struct {
	User    func(*UserMutator) ([]string, []any)
	Profile func(*ProfileMutator) ([]string, []any)
	Post    func(*PostMutator) ([]string, []any)
	Comment func(*CommentMutator) ([]string, []any)
	Group   func(*GroupMutator) ([]string, []any)
}{
	User: func(mutator *UserMutator) ([]string, []any) {
		var cols []string
		var vals []any
		if mutator.ID.IsSet {
			cols = append(cols, Schema.Public.User.ID.ColumnName())
			vals = append(vals, mutator.ID.Value)
		}
		if mutator.Email.IsSet {
			cols = append(cols, Schema.Public.User.Email.ColumnName())
			vals = append(vals, mutator.Email.Value)
		}
		if mutator.Age.IsSet {
			cols = append(cols, Schema.Public.User.Age.ColumnName())
			vals = append(vals, mutator.Age.Value)
		}
		if mutator.Tags.IsSet {
			cols = append(cols, Schema.Public.User.Tags.ColumnName())
			vals = append(vals, mutator.Tags.Value)
		}
		if mutator.Preferences.IsSet {
			cols = append(cols, Schema.Public.User.Preferences.ColumnName())
			vals = append(vals, mutator.Preferences.Value)
		}
		if mutator.Metadata.IsSet {
			cols = append(cols, Schema.Public.User.Metadata.ColumnName())
			vals = append(vals, mutator.Metadata.Value)
		}
		if mutator.CreatedAt.IsSet {
			cols = append(cols, Schema.Public.User.CreatedAt.ColumnName())
			vals = append(vals, mutator.CreatedAt.Value)
		}
		return cols, vals
	},
	Profile: func(mutator *ProfileMutator) ([]string, []any) {
		var cols []string
		var vals []any
		if mutator.ID.IsSet {
			cols = append(cols, Schema.Public.Profile.ID.ColumnName())
			vals = append(vals, mutator.ID.Value)
		}
		if mutator.UserID.IsSet {
			cols = append(cols, Schema.Public.Profile.UserID.ColumnName())
			vals = append(vals, mutator.UserID.Value)
		}
		if mutator.Bio.IsSet {
			cols = append(cols, Schema.Public.Profile.Bio.ColumnName())
			vals = append(vals, mutator.Bio.Value)
		}
		if mutator.Location.IsSet {
			cols = append(cols, Schema.Public.Profile.Location.ColumnName())
			vals = append(vals, mutator.Location.Value)
		}
		if mutator.ActiveDuration.IsSet {
			cols = append(
				cols,
				Schema.Public.Profile.ActiveDuration.ColumnName(),
			)
			vals = append(vals, mutator.ActiveDuration.Value)
		}
		if mutator.IsPublic.IsSet {
			cols = append(cols, Schema.Public.Profile.IsPublic.ColumnName())
			vals = append(vals, mutator.IsPublic.Value)
		}
		return cols, vals
	},
	Post: func(mutator *PostMutator) ([]string, []any) {
		var cols []string
		var vals []any
		if mutator.ID.IsSet {
			cols = append(cols, Schema.Public.Post.ID.ColumnName())
			vals = append(vals, mutator.ID.Value)
		}
		if mutator.UserID.IsSet {
			cols = append(cols, Schema.Public.Post.UserID.ColumnName())
			vals = append(vals, mutator.UserID.Value)
		}
		if mutator.Title.IsSet {
			cols = append(cols, Schema.Public.Post.Title.ColumnName())
			vals = append(vals, mutator.Title.Value)
		}
		if mutator.Content.IsSet {
			cols = append(cols, Schema.Public.Post.Content.ColumnName())
			vals = append(vals, mutator.Content.Value)
		}
		if mutator.Rating.IsSet {
			cols = append(cols, Schema.Public.Post.Rating.ColumnName())
			vals = append(vals, mutator.Rating.Value)
		}
		return cols, vals
	},
	Comment: func(mutator *CommentMutator) ([]string, []any) {
		var cols []string
		var vals []any
		if mutator.ID.IsSet {
			cols = append(cols, Schema.Public.Comment.ID.ColumnName())
			vals = append(vals, mutator.ID.Value)
		}
		if mutator.PostID.IsSet {
			cols = append(cols, Schema.Public.Comment.PostID.ColumnName())
			vals = append(vals, mutator.PostID.Value)
		}
		if mutator.UserID.IsSet {
			cols = append(cols, Schema.Public.Comment.UserID.ColumnName())
			vals = append(vals, mutator.UserID.Value)
		}
		if mutator.Text.IsSet {
			cols = append(cols, Schema.Public.Comment.Text.ColumnName())
			vals = append(vals, mutator.Text.Value)
		}
		return cols, vals
	},
	Group: func(mutator *GroupMutator) ([]string, []any) {
		var cols []string
		var vals []any
		if mutator.ID.IsSet {
			cols = append(cols, Schema.Public.Group.ID.ColumnName())
			vals = append(vals, mutator.ID.Value)
		}
		if mutator.Name.IsSet {
			cols = append(cols, Schema.Public.Group.Name.ColumnName())
			vals = append(vals, mutator.Name.Value)
		}
		return cols, vals
	},
}
