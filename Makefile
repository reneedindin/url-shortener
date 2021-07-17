docker_build:
	docker build -t url-shortener .

deploy:
	docker build -t url-shortener ./
	docker-compose up
