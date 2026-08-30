package contract

// ComparisonOperator holds standardized comparison logic strings for querying.
type ComparisonOperator string

// BooleanOperator represents a logical connector used to combine multiple
// SQL WHERE conditions, such as AND or OR.
type BooleanOperator string

// OpEqual represents the equal comparison operator.
const OpEqual ComparisonOperator = "="

// OpNotEqual represents the not-equal comparison operator.
const OpNotEqual ComparisonOperator = "!="

// OpLessThan represents the less-than comparison operator.
const OpLessThan ComparisonOperator = "<"

// OpGreaterThan represents the greater-than comparison operator.
const OpGreaterThan ComparisonOperator = ">"

// OpLessThanOrEqual represents the less-than-or-equal comparison operator.
const OpLessThanOrEqual ComparisonOperator = "<="

// OpGreaterThanOrEqual represents the greater-than-or-equal comparison
// operator.
const OpGreaterThanOrEqual ComparisonOperator = ">="

// OpIsDistinctFrom represents the IS DISTINCT FROM comparison operator.
const OpIsDistinctFrom ComparisonOperator = "IS DISTINCT FROM"

// OpIsNotDistinctFrom represents the IS NOT DISTINCT FROM comparison operator.
const OpIsNotDistinctFrom ComparisonOperator = "IS NOT DISTINCT FROM"

// OpLike represents the LIKE string comparison operator.
const OpLike ComparisonOperator = "LIKE"

// OpNotLike represents the NOT LIKE string comparison operator.
const OpNotLike ComparisonOperator = "NOT LIKE"

// OpILike represents the ILIKE case-insensitive string comparison operator.
const OpILike ComparisonOperator = "ILIKE"

// OpNotILike represents the NOT ILIKE case-insensitive string comparison
// operator.
const OpNotILike ComparisonOperator = "NOT ILIKE"

// OpRegexMatch represents the regular expression match operator.
const OpRegexMatch ComparisonOperator = "~"

// OpRegexIMatch represents the case-insensitive regular expression match
// operator.
const OpRegexIMatch ComparisonOperator = "~*"

// OpRegexNotMatch represents the regular expression non-matching operator.
const OpRegexNotMatch ComparisonOperator = "!~"

// OpRegexNotIMatch represents the case-insensitive regular expression
// non-matching operator.
const OpRegexNotIMatch ComparisonOperator = "!~*"

// OpSimilarTo represents the SIMILAR TO SQL string comparison operator.
const OpSimilarTo ComparisonOperator = "SIMILAR TO"

// OpNotSimilarTo represents the NOT SIMILAR TO SQL string comparison
// operator.
const OpNotSimilarTo ComparisonOperator = "NOT SIMILAR TO"

// OpBetween represents the BETWEEN range comparison operator.
const OpBetween ComparisonOperator = "BETWEEN"

// OpNotBetween represents the NOT BETWEEN range comparison operator.
const OpNotBetween ComparisonOperator = "NOT BETWEEN"

// OpIn represents the IN set membership comparison operator.
const OpIn ComparisonOperator = "IN"

// OpNotIn represents the NOT IN set non-membership comparison operator.
const OpNotIn ComparisonOperator = "NOT IN"

// OpJsonGetField represents the JSON get field (->) operator.
const OpJsonGetField ComparisonOperator = "->"

// OpJsonGetTextField represents the JSON get text field (->>) operator.
const OpJsonGetTextField ComparisonOperator = "->>"

// OpJsonGetPath represents the JSON get path (#>) operator.
const OpJsonGetPath ComparisonOperator = "#>"

// OpJsonGetPathText represents the JSON get path text (#>>) operator.
const OpJsonGetPathText ComparisonOperator = "#>>"

// OpJsonContains represents the JSON contains (@>) operator.
const OpJsonContains ComparisonOperator = "@>"

// OpJsonContainedBy represents the JSON contained by (<@) operator.
const OpJsonContainedBy ComparisonOperator = "<@"

// OpJsonKeyExist represents the JSON key exists (?) operator.
const OpJsonKeyExist ComparisonOperator = "?"

// OpJsonKeyAnyExist represents the JSON key any exists (?|) operator.
const OpJsonKeyAnyExist ComparisonOperator = "?|"

// OpJsonKeyAllExist represents the JSON key all exists (?&) operator.
const OpJsonKeyAllExist ComparisonOperator = "?&"

// OpJsonDeletePath represents the JSON delete path (#-) operator.
const OpJsonDeletePath ComparisonOperator = "#-"

// OpArrayContains represents the array contains (@>) operator.
const OpArrayContains ComparisonOperator = "@>"

// OpArrayContainedBy represents the array contained by (<@) operator.
const OpArrayContainedBy ComparisonOperator = "<@"

// OpArrayOverlap represents the array overlap (&&) operator.
const OpArrayOverlap ComparisonOperator = "&&"

// OpArrayConcat represents the array concatenation (||) operator.
const OpArrayConcat ComparisonOperator = "||"

// OpGeoOverlaps represents the geometric overlaps (&&) operator.
const OpGeoOverlaps ComparisonOperator = "&&"

// OpGeoContains represents the geometric contains (@>) operator.
const OpGeoContains ComparisonOperator = "@>"

// OpGeoContainedBy represents the geometric contained by (<@) operator.
const OpGeoContainedBy ComparisonOperator = "<@"

// OpGeoStrictLeft represents the geometric strictly to the left (<<) operator.
const OpGeoStrictLeft ComparisonOperator = "<<"

// OpGeoStrictRight represents the geometric strictly to the right (>>)
// operator.
const OpGeoStrictRight ComparisonOperator = ">>"

// OpGeoBelow represents the geometric strictly below (<^) operator.
const OpGeoBelow ComparisonOperator = "<^"

// OpGeoAbove represents the geometric strictly above (>^) operator.
const OpGeoAbove ComparisonOperator = ">^"

// OpGeoDistance represents the geometric distance (<->) operator.
const OpGeoDistance ComparisonOperator = "<->"

// OpGeoClosestProx represents the geometric closest proximity (##) operator.
const OpGeoClosestProx ComparisonOperator = "##"

// OpTextSearchMatch represents the text search match (@@) operator.
const OpTextSearchMatch ComparisonOperator = "@@"

// BoolAnd represents the logical AND boolean operator.
const BoolAnd BooleanOperator = "AND"

// BoolOr represents the logical OR boolean operator.
const BoolOr BooleanOperator = "OR"
