//Automatically generated file; DO NOT EDIT
package lib

import (
	"encoding/json"
	"errors"

	"github.com/micro/go-micro/v2/broker"
	"github.com/micro/go-micro/v2/logger"
)

//TODO chacnge
const orchestratorRoutingKey = "milestone.orchestrator"

var ErrToManyRetries = errors.New("the number of attempts is too large")

type direction int

const (
	Up direction = iota
	Down
)

func (d direction) Is() bool {
	return d == Up || d == Down
}

type EventTransmitter struct {
	t Transmitter
	m *Message
	broker.Event
}

func (e *EventTransmitter) Retry() int {
	return e.m.Retry
}

func (e *EventTransmitter) SetPayload(b []byte) {
	e.m.Payload = b
}

func (e *EventTransmitter) Payload() []byte {
	return e.m.Payload
}

func (e *EventTransmitter) ID() string {
	return e.m.ID
}

func (e *EventTransmitter) Approve() error {
	return e.t.Approve(e.m)
}

func (e *EventTransmitter) Reject() error {
	return e.t.Reject(e.m)
}

type Transmitter interface {
	Approve(m *Message) error
	Reject(m *Message) error
}

type steps string

const (
	verify_consumer steps = "verify_consumer"
	create_ticket   steps = "create_ticket"
	verify_card     steps = "verify_card"
	confirm_ticket  steps = "confirm_ticket"
	confirm_order   steps = "confirm_order"
)

type VerifyConsumer string

const (
	BeginVerifyVerifyConsumer VerifyConsumer = "begin_verify"
	CheckedVerifyConsumer     VerifyConsumer = "checked"
	FailedVerifyConsumer      VerifyConsumer = "failed"
)

type CreateTicket string

const (
	BeginCheckCreateTicket CreateTicket = "begin_check"
	VerifedCreateTicket    CreateTicket = "verifed"
	FailedCreateTicket     CreateTicket = "failed"
)

type VerifyCard string

const (
	BeginVerifyVerifyCard VerifyCard = "begin_verify"
	VerifedVerifyCard     VerifyCard = "verifed"
	FailedVerifyCard      VerifyCard = "failed"
)

type ConfirmTicket string

const (
	StartConfirmTicket   ConfirmTicket = "start"
	ConfirmConfirmTicket ConfirmTicket = "confirm"
	FailedConfirmTicket  ConfirmTicket = "Failed"
)

type ConfirmOrder string

const (
	StartConfirmOrder   ConfirmOrder = "start"
	ConfirmConfirmOrder ConfirmOrder = "confirm"
	FailedConfirmOrder  ConfirmOrder = "Failed"
)

func (t VerifyConsumer) Is() bool {
	if t == BeginVerifyVerifyConsumer || t == CheckedVerifyConsumer || t == FailedVerifyConsumer {
		return true
	}
	return false
}

func (t CreateTicket) Is() bool {
	if t == BeginCheckCreateTicket || t == VerifedCreateTicket || t == FailedCreateTicket {
		return true
	}
	return false
}

func (t VerifyCard) Is() bool {
	if t == BeginVerifyVerifyCard || t == VerifedVerifyCard || t == FailedVerifyCard {
		return true
	}
	return false
}

func (t ConfirmTicket) Is() bool {
	if t == StartConfirmTicket || t == ConfirmConfirmTicket || t == FailedConfirmTicket {
		return true
	}
	return false
}

func (t ConfirmOrder) Is() bool {
	if t == StartConfirmOrder || t == ConfirmConfirmOrder || t == FailedConfirmOrder {
		return true
	}
	return false
}

type Message struct {
	ID        string    `json:"id"`
	Command   string    `json:"command"`
	StepName  string    `json:"step_name"`
	Direction direction `json:"direction"`
	//The current number of the retry
	Retry int `json:"retry"`
	//If it need we can add payload to message.
	Payload []byte `json:"payload"`
}

func NewMessage(id string) *Message {
	m := &Message{
		ID:        id,
		Direction: Up,
	}
	return m
}

type VerifyConsumerTransmitter struct {
	b          broker.Broker
	MaxRetries int
}

func (t *VerifyConsumerTransmitter) Approve(m *Message) error {
	m.Command = string(CheckedVerifyConsumer)
	m.StepName = "verify_consumer"
	m.Retry = 0
	b, _ := json.Marshal(m)
	body := &broker.Message{Body: b}
	return t.b.Publish(orchestratorRoutingKey, body)
}

func (t *VerifyConsumerTransmitter) Reject(m *Message) error {

	m.Command = string(FailedVerifyConsumer)
	if m.Direction != Down {
		m.Command = string(CheckedVerifyConsumer)
	}
	m.Direction = Down

	m.StepName = "verify_consumer"
	if m.Retry > t.MaxRetries {
		return ErrToManyRetries
	}

	m.Retry++
	b, _ := json.Marshal(m)
	body := &broker.Message{Body: b}
	return t.b.Publish(orchestratorRoutingKey, body)
}
func NewVerifyConsumerTransmitter(b broker.Broker) *VerifyConsumerTransmitter {
	return &VerifyConsumerTransmitter{b: b, MaxRetries: 10}
}

type VerifyConsumerReceiver struct {
	b broker.Broker
	t Transmitter
}

func (r *VerifyConsumerReceiver) Pending(f func(*EventTransmitter) error) (broker.Subscriber, error) {
	return r.b.Subscribe("verify_consumer.pending", func(event broker.Event) error {
		m := Message{}
		err := json.Unmarshal(event.Message().Body, &m)
		if err != nil {
			panic(err)
		}
		return f(&EventTransmitter{
			t:     r.t,
			m:     &m,
			Event: event,
		})
	})
}

func (r *VerifyConsumerReceiver) Rejected(f func(*EventTransmitter) error) (broker.Subscriber, error) {
	return r.b.Subscribe("verify_consumer.rejected", func(event broker.Event) error {
		m := Message{}
		err := json.Unmarshal(event.Message().Body, &m)
		if err != nil {
			panic(err)
		}
		return f(&EventTransmitter{
			t:     r.t,
			m:     &m,
			Event: event,
		})
	})

}

func NewVerifyConsumerReceiver(b broker.Broker) *VerifyConsumerReceiver {
	return &VerifyConsumerReceiver{
		b: b,
		t: NewVerifyConsumerTransmitter(b),
	}
}

type CreateTicketTransmitter struct {
	b          broker.Broker
	MaxRetries int
}

func (t *CreateTicketTransmitter) Approve(m *Message) error {
	m.Command = string(VerifedCreateTicket)
	m.StepName = "create_ticket"
	m.Retry = 0
	b, _ := json.Marshal(m)
	body := &broker.Message{Body: b}
	return t.b.Publish(orchestratorRoutingKey, body)
}

func (t *CreateTicketTransmitter) Reject(m *Message) error {

	m.Command = string(FailedCreateTicket)
	if m.Direction != Down {
		m.Command = string(VerifedCreateTicket)
	}
	m.Direction = Down

	m.StepName = "create_ticket"
	if m.Retry > t.MaxRetries {
		return ErrToManyRetries
	}

	m.Retry++
	b, _ := json.Marshal(m)
	body := &broker.Message{Body: b}
	return t.b.Publish(orchestratorRoutingKey, body)
}
func NewCreateTicketTransmitter(b broker.Broker) *CreateTicketTransmitter {
	return &CreateTicketTransmitter{b: b, MaxRetries: 10}
}

type CreateTicketReceiver struct {
	b broker.Broker
	t Transmitter
}

func (r *CreateTicketReceiver) Pending(f func(*EventTransmitter) error) (broker.Subscriber, error) {
	return r.b.Subscribe("create_ticket.pending", func(event broker.Event) error {
		m := Message{}
		err := json.Unmarshal(event.Message().Body, &m)
		if err != nil {
			panic(err)
		}
		return f(&EventTransmitter{
			t:     r.t,
			m:     &m,
			Event: event,
		})
	})
}

func (r *CreateTicketReceiver) Rejected(f func(*EventTransmitter) error) (broker.Subscriber, error) {
	return r.b.Subscribe("create_ticket.rejected", func(event broker.Event) error {
		m := Message{}
		err := json.Unmarshal(event.Message().Body, &m)
		if err != nil {
			panic(err)
		}
		return f(&EventTransmitter{
			t:     r.t,
			m:     &m,
			Event: event,
		})
	})

}

func NewCreateTicketReceiver(b broker.Broker) *CreateTicketReceiver {
	return &CreateTicketReceiver{
		b: b,
		t: NewCreateTicketTransmitter(b),
	}
}

type VerifyCardTransmitter struct {
	b          broker.Broker
	MaxRetries int
}

func (t *VerifyCardTransmitter) Approve(m *Message) error {
	m.Command = string(VerifedVerifyCard)
	m.StepName = "verify_card"
	m.Retry = 0
	b, _ := json.Marshal(m)
	body := &broker.Message{Body: b}
	return t.b.Publish(orchestratorRoutingKey, body)
}

func (t *VerifyCardTransmitter) Reject(m *Message) error {

	m.Command = string(FailedVerifyCard)
	if m.Direction != Down {
		m.Command = string(VerifedVerifyCard)
	}
	m.Direction = Down

	m.StepName = "verify_card"
	if m.Retry > t.MaxRetries {
		return ErrToManyRetries
	}

	m.Retry++
	b, _ := json.Marshal(m)
	body := &broker.Message{Body: b}
	return t.b.Publish(orchestratorRoutingKey, body)
}
func NewVerifyCardTransmitter(b broker.Broker) *VerifyCardTransmitter {
	return &VerifyCardTransmitter{b: b, MaxRetries: 10}
}

type VerifyCardReceiver struct {
	b broker.Broker
	t Transmitter
}

func (r *VerifyCardReceiver) Pending(f func(*EventTransmitter) error) (broker.Subscriber, error) {
	return r.b.Subscribe("verify_card.pending", func(event broker.Event) error {
		m := Message{}
		err := json.Unmarshal(event.Message().Body, &m)
		if err != nil {
			panic(err)
		}
		return f(&EventTransmitter{
			t:     r.t,
			m:     &m,
			Event: event,
		})
	})
}

func (r *VerifyCardReceiver) Rejected(f func(*EventTransmitter) error) (broker.Subscriber, error) {
	return r.b.Subscribe("verify_card.rejected", func(event broker.Event) error {
		m := Message{}
		err := json.Unmarshal(event.Message().Body, &m)
		if err != nil {
			panic(err)
		}
		return f(&EventTransmitter{
			t:     r.t,
			m:     &m,
			Event: event,
		})
	})

}

func NewVerifyCardReceiver(b broker.Broker) *VerifyCardReceiver {
	return &VerifyCardReceiver{
		b: b,
		t: NewVerifyCardTransmitter(b),
	}
}

type ConfirmTicketTransmitter struct {
	b          broker.Broker
	MaxRetries int
}

func (t *ConfirmTicketTransmitter) Approve(m *Message) error {
	m.Command = string(ConfirmConfirmTicket)
	m.StepName = "confirm_ticket"
	m.Retry = 0
	b, _ := json.Marshal(m)
	body := &broker.Message{Body: b}
	return t.b.Publish(orchestratorRoutingKey, body)
}

func (t *ConfirmTicketTransmitter) Reject(m *Message) error {

	m.Command = string(FailedConfirmTicket)
	m.StepName = "confirm_ticket"
	if m.Retry > t.MaxRetries {
		return ErrToManyRetries
	}

	m.Retry++
	b, _ := json.Marshal(m)
	body := &broker.Message{Body: b}
	return t.b.Publish(orchestratorRoutingKey, body)
}
func NewConfirmTicketTransmitter(b broker.Broker) *ConfirmTicketTransmitter {
	return &ConfirmTicketTransmitter{b: b, MaxRetries: 10}
}

type ConfirmTicketReceiver struct {
	b broker.Broker
	t Transmitter
}

func (r *ConfirmTicketReceiver) Pending(f func(*EventTransmitter) error) (broker.Subscriber, error) {
	return r.b.Subscribe("confirm_ticket.pending", func(event broker.Event) error {
		m := Message{}
		err := json.Unmarshal(event.Message().Body, &m)
		if err != nil {
			panic(err)
		}
		return f(&EventTransmitter{
			t:     r.t,
			m:     &m,
			Event: event,
		})
	})
}

func (r *ConfirmTicketReceiver) Rejected(f func(*EventTransmitter) error) (broker.Subscriber, error) {
	return r.b.Subscribe("confirm_ticket.rejected", func(event broker.Event) error {
		m := Message{}
		err := json.Unmarshal(event.Message().Body, &m)
		if err != nil {
			panic(err)
		}
		return f(&EventTransmitter{
			t:     r.t,
			m:     &m,
			Event: event,
		})
	})

}

func NewConfirmTicketReceiver(b broker.Broker) *ConfirmTicketReceiver {
	return &ConfirmTicketReceiver{
		b: b,
		t: NewConfirmTicketTransmitter(b),
	}
}

type ConfirmOrderTransmitter struct {
	b          broker.Broker
	MaxRetries int
}

func (t *ConfirmOrderTransmitter) Approve(m *Message) error {
	m.Command = string(ConfirmConfirmOrder)
	m.StepName = "confirm_order"
	m.Retry = 0
	b, _ := json.Marshal(m)
	body := &broker.Message{Body: b}
	return t.b.Publish(orchestratorRoutingKey, body)
}

func (t *ConfirmOrderTransmitter) Reject(m *Message) error {

	m.Command = string(FailedConfirmOrder)
	m.StepName = "confirm_order"
	if m.Retry > t.MaxRetries {
		return ErrToManyRetries
	}

	m.Retry++
	b, _ := json.Marshal(m)
	body := &broker.Message{Body: b}
	return t.b.Publish(orchestratorRoutingKey, body)
}
func NewConfirmOrderTransmitter(b broker.Broker) *ConfirmOrderTransmitter {
	return &ConfirmOrderTransmitter{b: b, MaxRetries: 10}
}

type ConfirmOrderReceiver struct {
	b broker.Broker
	t Transmitter
}

func (r *ConfirmOrderReceiver) Pending(f func(*EventTransmitter) error) (broker.Subscriber, error) {
	return r.b.Subscribe("confirm_order.pending", func(event broker.Event) error {
		m := Message{}
		err := json.Unmarshal(event.Message().Body, &m)
		if err != nil {
			panic(err)
		}
		return f(&EventTransmitter{
			t:     r.t,
			m:     &m,
			Event: event,
		})
	})
}

func (r *ConfirmOrderReceiver) Rejected(f func(*EventTransmitter) error) (broker.Subscriber, error) {
	return r.b.Subscribe("confirm_order.rejected", func(event broker.Event) error {
		m := Message{}
		err := json.Unmarshal(event.Message().Body, &m)
		if err != nil {
			panic(err)
		}
		return f(&EventTransmitter{
			t:     r.t,
			m:     &m,
			Event: event,
		})
	})

}

func NewConfirmOrderReceiver(b broker.Broker) *ConfirmOrderReceiver {
	return &ConfirmOrderReceiver{
		b: b,
		t: NewConfirmOrderTransmitter(b),
	}
}

type Orchestrator struct {
	b   broker.Broker
	log logger.Logger
	//TODO add callback for bad transaction for example: use this if we reject transaction witch has rejected.
	//TODO add storage
}

func (o *Orchestrator) Do(options ...broker.SubscribeOption) (broker.Subscriber, error) {
	return o.b.Subscribe(orchestratorRoutingKey, o.handler, options...)
}

func (o *Orchestrator) handler(e broker.Event) error {
	m := Message{}
	err := json.Unmarshal(e.Message().Body, &m)
	//it mustn't ever happens
	if err != nil {
		panic(err)
	}

	switch steps(m.StepName) {
	case verify_consumer:
		return o.verify_consumerRoute(e.Message(), VerifyConsumer(m.Command), m.Direction)
	case create_ticket:
		return o.create_ticketRoute(e.Message(), CreateTicket(m.Command), m.Direction)
	case verify_card:
		return o.verify_cardRoute(e.Message(), VerifyCard(m.Command), m.Direction)
	case confirm_ticket:
		return o.confirm_ticketRoute(e.Message(), ConfirmTicket(m.Command), m.Direction)
	case confirm_order:
		return o.confirm_orderRoute(e.Message(), ConfirmOrder(m.Command), m.Direction)

	}
	return nil
}

func (o *Orchestrator) verify_consumerRoute(m *broker.Message, typeOf VerifyConsumer, direction direction) error {
	if !typeOf.Is() {
		return errors.New("invalid typeOf")
	}

	switch direction {
	case Up:
		switch typeOf {
		case CheckedVerifyConsumer:
			return o.b.Publish("create_ticket.pending", m)
		default:
			panic(typeOf)
		}
	case Down:
		switch typeOf {
		case CheckedVerifyConsumer:
			return nil
		case FailedVerifyConsumer:
			return o.b.Publish("verify_consumer.rejected", m)
		default:
			panic(typeOf)
		}
	}

	return nil
}

func (o *Orchestrator) create_ticketRoute(m *broker.Message, typeOf CreateTicket, direction direction) error {
	if !typeOf.Is() {
		return errors.New("invalid typeOf")
	}

	switch direction {
	case Up:
		switch typeOf {
		case VerifedCreateTicket:
			return o.b.Publish("verify_card.pending", m)
		default:
			panic(typeOf)
		}
	case Down:
		switch typeOf {
		case VerifedCreateTicket:
			return o.b.Publish("verify_consumer.rejected", m)
		case FailedCreateTicket:
			return o.b.Publish("create_ticket.rejected", m)
		default:
			panic(typeOf)
		}
	}

	return nil
}

func (o *Orchestrator) verify_cardRoute(m *broker.Message, typeOf VerifyCard, direction direction) error {
	if !typeOf.Is() {
		return errors.New("invalid typeOf")
	}

	switch direction {
	case Up:
		switch typeOf {
		case VerifedVerifyCard:
			return o.b.Publish("confirm_ticket.pending", m)
		default:
			panic(typeOf)
		}
	case Down:
		switch typeOf {
		case VerifedVerifyCard:
			return o.b.Publish("create_ticket.rejected", m)
		case FailedVerifyCard:
			return o.b.Publish("verify_card.rejected", m)
		default:
			panic(typeOf)
		}
	}

	return nil
}

func (o *Orchestrator) confirm_ticketRoute(m *broker.Message, typeOf ConfirmTicket, direction direction) error {
	if !typeOf.Is() {
		return errors.New("invalid typeOf")
	}

	switch typeOf {
	case ConfirmConfirmTicket:
		return o.b.Publish("confirm_order.pending", m)
	case FailedConfirmTicket:
		return o.b.Publish("confirm_ticket.rejected", m)
	default:
		panic(typeOf)
	}

	return nil
}

func (o *Orchestrator) confirm_orderRoute(m *broker.Message, typeOf ConfirmOrder, direction direction) error {
	if !typeOf.Is() {
		return errors.New("invalid typeOf")
	}

	switch typeOf {
	case ConfirmConfirmOrder:
		return nil
	case FailedConfirmOrder:
		return o.b.Publish("confirm_order.rejected", m)
	default:
		panic(typeOf)
	}

	return nil
}

func NewOrchestrator(b broker.Broker, log logger.Logger) *Orchestrator {
	if log == nil {
		log = logger.NewLogger()
	}
	return &Orchestrator{
		b:   b,
		log: log,
	}
}
