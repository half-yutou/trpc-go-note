package main

import (
	"errors"
	"fmt"

	"trpc.group/trpc-go/trpc-go/errs"
)

// 定义业务错误码
const (
	ErrCodeUserNotFound = 10001
	ErrCodeBalanceLow   = 10002
)

func main() {
	// 1. 开启堆栈跟踪 (通常在 init 或 main 开头)
	errs.SetTraceable(true)

	fmt.Println("=== 场景 1: 业务逻辑返回错误 ===")
	if err := findUser("999"); err != nil {
		handleError(err)
	}

	fmt.Println("\n=== 场景 2: 模拟框架错误 (如超时) ===")
	if err := mockFrameworkCall(); err != nil {
		handleError(err)
	}

	fmt.Println("\n=== 场景 3: 错误包装 (Wrap) ===")
	if err := checkBalance("123"); err != nil {
		// 打印详细信息 (包含堆栈，因为开启了 SetTraceable)
		fmt.Printf("Detailed Error: %+v\n", err)
	}
}

// 模拟业务函数
func findUser(uid string) error {
	// 直接返回一个带错误码的业务错误
	return errs.Newf(ErrCodeUserNotFound, "user %s not found", uid)
}

// 模拟框架行为
func mockFrameworkCall() error {
	// 返回一个标准的框架错误 (例如 Server Timeout)
	return errs.NewFrameError(errs.RetServerTimeout, "upstream request timeout")
}

// 模拟底层错误包装
func checkBalance(uid string) error {
	// 假设这是底层 DB 返回的原生 error
	dbErr := errors.New("connection reset")

	// 包装它，赋予业务含义
	return errs.Wrapf(dbErr, ErrCodeBalanceLow, "failed to check balance for %s", uid)
}

// 统一错误处理逻辑 (模拟 Client 端或上层调用者)
func handleError(err error) {
	// 1. 提取错误码
	code := errs.Code(err)

	// 2. 提取错误信息
	msg := errs.Msg(err)

	fmt.Printf("[Error Handled] Code: %d, Msg: %s\n", code, msg)

	// 3. 根据错误码做特定逻辑
	switch code {
	case ErrCodeUserNotFound:
		fmt.Println("-> Action: Redirect to registration page")
	case errs.RetServerTimeout:
		fmt.Println("-> Action: Retry request (Network issue)")
	default:
		fmt.Println("-> Action: Report to admin")
	}
}
