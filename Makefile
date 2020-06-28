.PHONY: clean dev
clean:
	rm -rf ./tmp
	rm -rf ./dist

dev: clean
	goreleaser --snapshot --skip-publish --rm-dist
