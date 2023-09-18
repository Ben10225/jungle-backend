build-docker-image:
	docker build -t jungle:latest . \
	&& docker tag jungle benpeng/jungle \
	&& docker push benpeng/jungle

.PHONY: build-docker-image