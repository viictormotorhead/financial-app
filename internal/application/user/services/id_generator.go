package services

type IDGenerator interface {
	New() (string, error)
}
