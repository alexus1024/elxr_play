-include .env
export

APPS := api

# --- RPi5 (k3s) ---
RPI_HOST  ?= root@elxr-rpi5
RPI_IMAGE  = alexus_test_apps/elxr_rp5

# --- EKS (AWS) ---
AWS_REGION ?= ca-central-1
ECR_URL    ?= UPDATE_ME.dkr.ecr.ca-central-1.amazonaws.com

.PHONY: build-rp5 build-eks load-rp5 deploy-rp5 deploy-eks apply-config-rp5 deploy-nats-rpi5 undeploy-nats-rpi5 clean $(APPS)

# === Pre-build/codegen ===

swagger:
	cd api && swag init

# === Builds ===

build-rp5: swagger
	@for app in $(APPS); do \
		docker build --platform linux/arm64 -t $(RPI_IMAGE)/$$app:latest $$app; \
	done

build-eks: swagger
	@for app in $(APPS); do \
		docker build --platform linux/amd64 -t $(ECR_URL)/$$app:latest $$app; \
	done

# === Helm (RPi5) ===

deploy-nats-rpi5:
	ssh $(RPI_HOST) helm upgrade --install nats nats/nats \
		--namespace elxr \
		--values - < deploy/helm/nats/values.yaml

undeploy-nats-rpi5:
	ssh $(RPI_HOST) helm uninstall nats --namespace elxr

# === RPi5 ===

load-rp5: build-rp5
	@for app in $(APPS); do \
		echo "==> Loading $(RPI_IMAGE)/$$app to $(RPI_HOST)"; \
		ssh $(RPI_HOST) k3s ctr images rm docker.io/$(RPI_IMAGE)/$$app:latest 2>/dev/null || true; \
		docker save $(RPI_IMAGE)/$$app:latest | ssh $(RPI_HOST) k3s ctr images import -; \
	done

apply-config-rp5:
	kubectl kustomize deploy/overlays/rpi | ssh $(RPI_HOST) kubectl apply -f -

deploy-rp5: load-rp5
	kubectl kustomize deploy/overlays/rpi | ssh $(RPI_HOST) kubectl apply -f -
	@for app in $(APPS); do \
		ssh $(RPI_HOST) kubectl rollout restart deployment/$$app -n elxr; \
		
	done
# ssh $(RPI_HOST) kubectl delete pods -l app=$$app -n elxr; \ - replace rollout if something wrong with deployment

# === EKS ===

push-eks: build-eks
	aws ecr get-login-password --region $(AWS_REGION) | docker login --username AWS --password-stdin $(ECR_URL)
	@for app in $(APPS); do \
		docker push $(ECR_URL)/$$app:latest; \
	done

deploy-eks: push-eks
	kubectl apply -k deploy/overlays/aws
	@for app in $(APPS); do \
		kubectl rollout restart deployment/$$app -n elxr; \
	done

watch-eks:
	kubectl get pods -n elxr -w

# === Clean ===

clean-rp5:
	@for app in $(APPS); do \
		docker rmi $(RPI_IMAGE)/$$app:latest 2>/dev/null || true; \
	done

clean-eks:
	@for app in $(APPS); do \
		docker rmi $(ECR_URL)/$$app:latest 2>/dev/null || true; \
	done
