# Pulse — build targets
#
# make build      → build ./pulse binary
# make install    → build + install to ~/.local/bin (no sudo needed)
# make uninstall  → remove the binary and shell hook
# make clean      → remove the local binary artifact

BINARY      = pulse
VERSION     = $(shell grep AppVersion internal/config/config.go | sed 's/.*"\(.*\)"/\1/')
LDFLAGS     = -ldflags="-X github.com/devpulse-cli/devpulse/internal/config.AppVersion=$(VERSION) -s -w"
INSTALL_DIR = $(HOME)/.local/bin

.PHONY: build install uninstall clean seed

build:
	go build $(LDFLAGS) -o $(BINARY) .

install: build
	@mkdir -p $(INSTALL_DIR)
	@cp $(BINARY) $(INSTALL_DIR)/$(BINARY)
	@echo "✓ installed pulse to $(INSTALL_DIR)/$(BINARY)"
	@echo ""
	@echo "  run: pulse init"
	@echo "  then: source ~/.zshrc"

uninstall:
	@rm -f $(INSTALL_DIR)/$(BINARY)
	@echo "✓ removed $(INSTALL_DIR)/$(BINARY)"
	@echo ""
	@echo "  to clean up your shell config too, run: pulse uninstall"

clean:
	@rm -f $(BINARY)

seed:
	@echo "seeding test data..."
	@for i in $$(seq 1 20); do \
		$(INSTALL_DIR)/$(BINARY) log --cmd "git status"    --exit 0 --ms 120   --dir "$$(pwd)"; \
		$(INSTALL_DIR)/$(BINARY) log --cmd "npm run dev"   --exit 0 --ms 3200  --dir "$$(pwd)"; \
		$(INSTALL_DIR)/$(BINARY) log --cmd "go build ./..."--exit 0 --ms 800   --dir "$$(pwd)"; \
		$(INSTALL_DIR)/$(BINARY) log --cmd "vim main.go"   --exit 0 --ms 45000 --dir "$$(pwd)"; \
	done
	@echo "✓ done — run 'pulse stats'"
