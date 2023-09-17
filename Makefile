.PHONY: user
user:
	@echo "running user srv"
	@go run ./app/service/rpc/user &

.PHONY: http
 http:
	@echo "running http api"
	@go run ./app/service/api/http &

.PHONY: msg
msg:
	@echo "running msg srv"
	@go run ./app/service/rpc/msg &

.PHONY: ws
ws:
	@echo "running ws api"
	@go run ./app/service/api/ws &

.PHONY: trans
trans:
	@echo "running mq transfer"
	@go run ./app/service/mq/trans &

