GINKGO := $$GOPATH/bin/ginkgo
MONGO_URL := mongodb://mongo_user:mongo_secret@0.0.0.0:27017/kudos

setup: run_services
	@MONGO_URL=${MONGO_URL} go run ./cmd/db/setup.go

run_services:
	@docker-compose up --build -d

run_tests: run_services
	@${GINKGO} pkg/**

run_server:
	@MONGO_URL=${MONGO_URL} PORT=4444 go run cmd/main.go

run_client:
	@/bin/bash -c "cd $$GOPATH/src/github.com/klebervirgilio/vue-crud-app-with-golang/pkg/http/web/app && yarn serve"