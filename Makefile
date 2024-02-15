.PHONY: build clean install install-binary install-resources

BINARY_PATH=./build/cloak
INSTALL_PATH=/usr/local/bin/cloak
WWW_PATH=/var/www/cloak
LOG_PATH=/var/log/cloak
ETC_PATH=/etc/cloak

# Check if we are the root user. If UID is 0 (root), SUDO will be an empty string; otherwise, it will be `sudo`.
UID := $(shell id -u)
SUDO := $(if $(filter 0,$(UID)),,sudo)

build: clean
	@echo "Building cloak..."
	@go build -o $(BINARY_PATH) .

clean:
	@echo "Cleaning..."
	@rm -f $(BINARY_PATH)

install: install-binary install-resources

install-binary: build
	@echo "Do you wish to install the binary to $(INSTALL_PATH)? [y/N]" && read ans && [ $${ans:-N} = y ]
	@$(SUDO) mv $(BINARY_PATH) $(INSTALL_PATH)
	@echo "Installed cloak to $(INSTALL_PATH)"

install-resources:
	@echo "Do you wish to install resources to /var and /etc? [y/N]" && read ans && [ $${ans:-N} = y ]
	@$(SUDO) mkdir -p $(WWW_PATH) && $(SUDO) cp -r stuff/www/* $(WWW_PATH)
	@$(SUDO) mkdir -p $(LOG_PATH) && $(SUDO) cp -r stuff/log/* $(LOG_PATH)
	@$(SUDO) mkdir -p $(ETC_PATH)
	@$(SUDO) cp stuff/dictionary.txt $(ETC_PATH)/dictionary.txt
	@$(SUDO) cp stuff/domains.txt $(ETC_PATH)/domains.txt
	@$(SUDO) cp stuff/cloak.cfg $(ETC_PATH)/cloak.cfg
	@echo "Resources installed to /var and /etc"
