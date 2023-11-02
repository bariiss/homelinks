docker-build:
	docker build --no-cache -t homelinks .

docker-run-tty: docker-build
	docker run -it --rm -p 8080:8080 homelinks

docker-run-daemon: docker-stop docker-build
	docker run -d --rm -p 8080:8080 homelinks

docker-stop:
	docker stop homelinks || echo "Container not running."

docker-image-rm:
	docker image rm homelinks || echo "Image not found."