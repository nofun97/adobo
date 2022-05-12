ARRAI_FILES = $(shell find . -name '*.arrai')
TEST_FILES  = $(shell find ./testdata/ -name '*.json')
ARRAI		= docker run --rm -w $$(pwd) -v $$(pwd):$$(pwd) anzbank/arrai

ifdef LOCAL_ARRAI
ARRAI=arrai
endif

.PHONY: build
build: assets
	go build -o adobo main.go

.PHONY: assets
assets: pkg/bundles/generator.arraiz

pkg/bundles/generator.arraiz: $(ARRAI_FILES)
	$(ARRAI) bundle --out $@ pkg/arrai/generate_from_json_path.arrai

.PHONY: testdata
testdata: $(TEST_FILES) testdata/testdata.arrai
	$(ARRAI) r --out=dir:testdata testdata/testdata.arrai

.PHONY: test
test: testdata $(ARRAI_FILES)
	$(ARRAI) test

.PHONY: check-clean
check-clean: assets testdata
	go mod tidy
	git --no-pager diff HEAD && test -z "$$(git status --porcelain)"
