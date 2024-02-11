.PHONY: dev build preview test fmt lint

dev: build
	npm run preview

build:
	npm install
	npm run build

test:
	npm install
	npm run test
	npm run test:unit

fmt:
	npm run format

lint:
	npm run lint