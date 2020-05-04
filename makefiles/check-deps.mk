.PHONY: check-deps
check-deps:
	env GOPRIVATE=${GOPRIVATE} go list -u -m -json all | go-mod-outdated -update -direct