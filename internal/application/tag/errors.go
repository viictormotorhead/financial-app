package tag

import "errors"

var ErrTagNameAlreadyExists = errors.New("a tag with this name already exists")

type TagsNotFoundError struct {
	Missing []string
}

func (e *TagsNotFoundError) Error() string {
	return "one or more tags do not exist"
}

func NewTagsNotFoundError(missing []string) *TagsNotFoundError {
	return &TagsNotFoundError{Missing: missing}
}
