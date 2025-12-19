.PHONY: imports tools

# 安装工具
tools:
	go install github.com/daixiang0/gci@latest

# 整理 imports
# -s standard: 标准库
# -s "prefix(trpc-go-note)": 本地包 (二方包)
# -s default: 第三方包
imports:
	gci write --custom-order -s standard -s "prefix(trpc-go-note)" -s default .
