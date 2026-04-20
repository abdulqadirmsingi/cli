# DevPulse build targets
#
# make build   → builds ./pulse binary
# make install → builds + moves to /usr/local/bin
# make clean   → removes the binary

BINARY   = pulse
VERSION  = $(shell grep AppVersion internal/config/config.go | sed 's/.*"\(.*\)"/\1/')
LDFLAGS  = -ldflags="-X github.com/devpulse-cli/devpulse/internal/config.AppVersion=$(VERSION) -s -w"

.PHONY: build install clean run

build:
	go build $(LDFLAGS) -o $(BINARY) .

install: build
	mv $(BINARY) /usr/local/bin/$(BINARY)
	@echo "✓ installed pulse to /usr/local/bin"

clean:
	rm -f $(BINARY)

# Seed some test data so `pulse stats` shows real output during development
seed:
	@echo "seeding test commands..."
	@for i in $$(seq 1 20); do \
		pulse log --cmd "git status" --exit 0 --ms 120 --dir "$$(pwd)"; \
		pulse log --cmd "npm run dev" --exit 0 --ms 3200 --dir "$$(pwd)"; \
		pulse log --cmd "go build ./..." --exit 0 --ms 800 --dir "$$(pwd)"; \
		pulse log --cmd "vim main.go" --exit 0 --ms 45000 --dir "$$(pwd)"; \
	done
	@echo "✓ done — run 'pulse stats' to see your data"
