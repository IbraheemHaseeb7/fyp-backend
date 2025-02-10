package types

type RequestMethod int

const (
	CREATE = iota
	READ_ONE
	READ_ALL
	UPDATE
	DELETE
)
