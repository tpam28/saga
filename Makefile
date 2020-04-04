#generate saga output
gen:
	@go run main.go --path=example/config/kitchen.yaml --output=example/saga/lib/saga.go

test:
	rm -rf example/build
	mkdir example/build/
	@go build example/saga/main.go -o example/build/saga
	@go build example/vf_consumer/main.go -o example/build/vf_consumer
	@go build example/mk_tiket/main.go -o example/build/mk_tiket
	@go build example/vf_card/main.go -o example/build/vf_card
	@go build example/confirm_tiket/main.go -o example/build/confirm_tiket
	@go build example/confirm_order/main.go -o example/build/confirm_order
	./example/build/saga &
	./example/build/vf_consumer &
	./example/build/mk_tiket &
	./example/build/vf_card &
	./example/build/confirm_tiket &
	./example/build/confirm_order