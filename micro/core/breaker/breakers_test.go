package breaker

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

var (
	errTestCustom = errors.New("自定义错误")
	errTestOther  = errors.New("其他错误")
	errTestFail   = errors.New("降级错误")
)

func opt() {
	Opt("any",
		WithName("any"),
		WithTimeout(3*time.Second),
	)
}

func TestGetBreaker(t *testing.T) {
	opt()

	for i := 0; i < 25; i++ {
		fmt.Println("=============>>>>>>>>>> in start", i)

		testDoWithFallbackAcceptable("any", i)

		fmt.Println("============>>>>>>>>>>> in end", i)
		fmt.Println("")
	}
}

func testDo(name string, i int) {
	Do(name, func() error {
		fmt.Println("testDo 测试:", i)
		if i < 15 {
			return errTestCustom
		}

		return nil
	})
}

func testDoWithAcceptable(name string, i int) {
	DoWithAcceptable(name, func() error {
		fmt.Println("testDoWithAcceptable 测试:", i)
		if i%3 == 0 {
			return errTestCustom
		}

		if i%4 == 0 {
			return errTestOther
		}

		return nil
	}, func(err error) bool {
		fmt.Println(err, i, err == errTestCustom)
		return err != errTestCustom
	})
}
func testDoWithFallback(name string, i int) {
	err := DoWithFallback(name, func() error {
		fmt.Println("testDoWithFallback 测试:", i)
		if i < 15 {
			return errTestCustom
		}

		return nil
	}, func(err error) error {
		time.Sleep(time.Second * 1)
		fmt.Println("降级")
		return errTestFail
	})

	fmt.Println("testDoWithFallback:", err)
}

func testDoWithFallbackAcceptable(name string, i int) {
	err := DoWithFallbackAcceptable(name, func() error {
		fmt.Println("testDoWithFallback 测试:", i)
		if i < 8 {
			return errTestCustom
		}

		if i < 14 {
			return errTestOther
		}

		if i < 22 {
			return errTestFail
		}

		return nil
	}, func(err error) error {
		time.Sleep(time.Second * 1)
		fmt.Println("降级")
		return errTestFail
	}, func(err error) bool {
		fmt.Println(err, i, err == errTestCustom)
		return err != errTestCustom
	})

	fmt.Println("testDoWithFallbackAcceptable:", err)
}
