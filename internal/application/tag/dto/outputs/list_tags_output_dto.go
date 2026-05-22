package outputs

type ListTagsOutputDTO struct {
	Tags []TagOutputDTO
}

type TagOutputDTO struct {
	ID          uint
	Name        string
	Description string
}
