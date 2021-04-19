package main

import (
	stdErr "errors"
	"fmt"

	"github.com/cockroachdb/errors"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
)

var (
	ErrCloud = stdErr.New("cloud provider failed")
	// ErrCloud = errutil.NewWithDepth(1, "cloud provider failed")
	ErrServer = stdErr.New("server failed")
	// ErrServer = errutil.NewWithDepth(1, "server failed")
)

func main() {
	// fmt.Printf("err=%v", useCaseWithCombine())
	fmt.Printf("err=%+v", errors.Unwrap(useCaseWithMsg()))
	// fmt.Printf("err=%+v", useCaseWithMsg())
	// fmt.Printf("err=%+v", errors.Unwrap(useCaseWithCombine()))
}

func useCaseWithCombine() error {
	err := domain()
	if err != nil {
		Err := errors.CombineErrors(ErrServer, err)
		return errors.Wrap(Err, "use case failed")
	}
	return nil
}

func useCaseWithMark1() error {
	err := repo()
	if err != nil {
		Err := errors.Mark(ErrServer, err)
		return errors.Wrap(Err, "use case failed")
	}
	return nil
}
func useCaseWithMsg() error {
	err := domain()
	if err != nil {
		return errors.WithMessage(err, "use case failed")
	}
	return nil
}

func domain() error {
	return errors.Wrap(ErrCloud, "user key=XXX is not exist")
}

func useCaseWithMark2() error {
	err := domain()
	if err != nil {
		Err := errors.Mark(ErrServer, err)
		return errors.Wrap(Err, "use case failed")
	}
	return nil
}

func repo() error {
	return mysql.ErrNoTLS
}
