package err

import (
	"fmt"
)

type DalRepositoryError struct {
	repo string
	msg  string
	err  error
}

var ErrDalRepository *DalRepositoryError

func NewDalRepositoryError(repo, msg string, err error) *DalRepositoryError {
	return &DalRepositoryError{
		repo: repo,
		msg:  msg,
		err:  err,
	}
}

func (e *DalRepositoryError) Error() string {
	return fmt.Sprintf("DAL repository [%s] error with message [%s] with error [%v]", e.repo, e.msg, e.err)
}

func (e *DalRepositoryError) Unwrap() error {
	return e.err
}
