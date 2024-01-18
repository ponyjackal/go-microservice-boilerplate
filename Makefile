help:
	@echo ''
	@echo 'Usage: make [TARGET] [EXTRA_ARGUMENTS]'
	@echo 'Targets:'
	@echo 'make dev: make dev for development work'
	@echo 'make staging: make staging for test work'
	@echo 'make build: make build container'
	@echo 'make production: docker production build'
	@echo 'clean: clean for all clear docker images'

doc:
	swag init

protobuf:
	buf generate proto

update-proto:
	buf mod update proto

dev:
	cp .env.local .env
	docker-compose -f docker-compose-dev.yml down
	docker-compose -f docker-compose-dev.yml up

build:
	docker-compose -f docker-compose-prod.yml build
	docker-compose -f docker-compose-dev.yml down build

prod:
	cp .env.prod .env
	docker-compose -f docker-compose-staging.yml down
	docker-compose -f docker-compose-staging.yml up	

clean:
	docker-compose -f docker-compose-prod.yml down -v
	docker-compose -f docker-compose-dev.yml down -v
