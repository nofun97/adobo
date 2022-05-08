ARRAI_FILES = $(shell find . -name '*.arrai')
TEST_FILES  = $(shell find ./testdata/ -name '*.json')
ARRAI		= go run github.com/arr-ai/arrai/cmd/arrai

ifdef $(LOCAL)
ARRAI=arrai
endif

build: assets
	go build -o adobo main.go

.PHONY: assets
assets: pkg/bundles/generator.arraiz

pkg/bundles/generator.arraiz: $(ARRAI_FILES)
	$(ARRAI) bundle --out $@ pkg/arrai/generate_from_json_path.arrai

testdata: $(TEST_FILES) testdata/testdata.arrai
	$(ARRAI) r --out=dir:testdata testdata/testdata.arrai

test: testdata $(ARRAI_FILES)
	$(ARRAI) test

check-clean: assets testdata
	git --no-pager diff HEAD && test -z "$$(git status --porcelain)"
