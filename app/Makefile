DOCKER_IMAGE = d-img-gen-demo
DOCKER_IMAGE_LOCAL = d-img-gen-demo-local:latest
ECR_REGISTRY = 
IMAGE_TAG = latest
LAMBDA_FUNCTION = dynamic-image-lambda
PUBLIC_ECR_REGISTRY = public.ecr.aws/a4a4z6p8

local-build:
	docker build --target local -t $(DOCKER_IMAGE_LOCAL) .

local-run:
	docker run --rm -it -p 8080:8080 $(DOCKER_IMAGE_LOCAL)

local-stop:
	docker stop $(shell docker ps -q --filter ancestor=$(DOCKER_IMAGE_LOCAL))

local-clean:
	docker rmi $(DOCKER_IMAGE_LOCAL)

lambda-build:
	docker build --target lambda -t $(DOCKER_IMAGE):$(IMAGE_TAG) .

ecr-push:
	aws ecr get-login-password --region ap-northeast-1 | docker login --username AWS --password-stdin $(ECR_REGISTRY)
	docker build --target lambda -t $(DOCKER_IMAGE) . --platform linux/x86_64
	docker tag $(DOCKER_IMAGE):$(IMAGE_TAG) $(ECR_REGISTRY)/$(DOCKER_IMAGE):$(IMAGE_TAG)
	docker push $(ECR_REGISTRY)/$(DOCKER_IMAGE):$(IMAGE_TAG)

public-ecr-push:
	aws ecr-public get-login-password --region us-east-1 | docker login --username AWS --password-stdin $(PUBLIC_ECR_REGISTRY)
	docker build --target lambda -t $(DOCKER_IMAGE) . --platform linux/x86_64
	docker tag $(DOCKER_IMAGE):$(IMAGE_TAG) $(PUBLIC_ECR_REGISTRY)/$(DOCKER_IMAGE):$(IMAGE_TAG)
	docker push $(PUBLIC_ECR_REGISTRY)/$(DOCKER_IMAGE):$(IMAGE_TAG)

lambda-deploy:
	aws lambda update-function-code --function-name $(LAMBDA_FUNCTION) --image-uri $(ECR_REGISTRY)/$(DOCKER_IMAGE):$(IMAGE_TAG)
	aws lambda wait function-updated --function-name $(LAMBDA_FUNCTION)

mockgen:
	mockgen -source=$(file) -destination=./internal/test/mock/$(patsubst internal/%,%,$(file))

test:
	go test ./... -v