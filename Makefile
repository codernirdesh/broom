APP         := broom
OUT_DIR     := bin
INSTALL_DIR ?= /mnt/c/Windows/System32

# Go build flags
GOOS        := windows
GOARCH      := amd64
LDFLAGS     := -s -w
GUI_FLAGS   := -H windowsgui

.PHONY: all build build-console install uninstall clean tidy

all: build

build:
	@mkdir -p $(OUT_DIR)
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags="$(LDFLAGS) $(GUI_FLAGS)" -o $(OUT_DIR)/$(APP).exe ./cmd/broom
	@echo "Built: $(OUT_DIR)/$(APP).exe"

build-console:
	@mkdir -p $(OUT_DIR)
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags="$(LDFLAGS)" -o $(OUT_DIR)/$(APP).exe ./cmd/broom
	@echo "Built (console): $(OUT_DIR)/$(APP).exe"

tidy:
	go mod tidy

install: build
	cp $(OUT_DIR)/$(APP).exe $(INSTALL_DIR)/$(APP).exe
	@echo "Installed: $(INSTALL_DIR)/$(APP).exe"

uninstall:
	rm -f $(INSTALL_DIR)/$(APP).exe
	@echo "Uninstalled: $(INSTALL_DIR)/$(APP).exe"

clean:
	rm -rf $(OUT_DIR)
