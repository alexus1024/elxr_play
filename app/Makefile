PLATFORM ?= linux/arm64
APPS := api
RPI_HOST ?= root@elxr-rpi5
OVERLAY ?= rpi

.PHONY: build load deploy clean $(APPS)

build: $(APPS)

$(APPS):
	docker build --platform $(PLATFORM) -t alexus_test_apps/elxr_rp5/$@:latest $@

load:
	@for app in $(APPS); do \
		echo "==> Loading alexus_test_apps/elxr_rp5/$$app to $(RPI_HOST)"; \
		docker save alexus_test_apps/elxr_rp5/$$app:latest | ssh $(RPI_HOST) k3s ctr images import -; \
	done

deploy: load
	kubectl kustomize deploy/overlays/$(OVERLAY) | ssh $(RPI_HOST) kubectl apply -f -

clean:
	@for app in $(APPS); do \
		docker rmi alexus_test_apps/elxr_rp5/$$app:latest 2>/dev/null || true; \
	done
