#generate saga output
RABBITMQ_DSN ?= "amqp://evgen:wZCfo9@127.0.0.1:5672/test"

gen:
	@go run main.go --path=example/config/kitchen.yaml --output=example/saga/lib/saga.go


test-build:
	$(shell cd example/saga ; go build -o ../../example/_build/saga)
	$(shell cd example/vf_consumer ; go build -o ../../example/_build/vf_consumer)
	$(shell cd  example/mk_tiket ; go build -o ../../example/_build/mk_tiket)
	$(shell cd  example/vf_card ; go build -o ../../example/_build/vf_card)
	$(shell cd  example/confirm_tiket ; go build -o ../../example/_build/confirm_tiket)
	$(shell cd  example/confirm_order ; go build -o ../../example/_build/confirm_order)
	$(shell chmod -R 777 example/_build/)

.ONESHELL:
test-run:
	cd example/_build
	./saga &
	./vf_consumer  $(args) &
	./mk_tiket  $(args) &
	./vf_card  $(args) &
	./confirm_tiket  $(args) &
	./confirm_order
