#generate saga output
gen:
	@go run main.go --path=example/config/kitchen.yaml --output=example/saga/lib/saga.go

test:
	@go run example/saga/main.go &
	@go run example/vf_consumer/main.go &
	@go run example/mk_tiket/main.go &
	@go run example/vf_card/main.go &
	@go run example/confirm_tiket/main.go &
	@go run example/confirm_order/main.go