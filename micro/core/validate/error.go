package validate

type Error struct {
	err string
}

func (e *Error) String() string {
	return e.err
}

var Err = new(Error)
