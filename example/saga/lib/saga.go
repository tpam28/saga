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

var ErrMethodNotAvailable = errors.New("method not available")
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
	t  Transmitter
	id string
	m  *Message
	broker.Event
}

func (e *EventTransmitter) ID() string {
	return e.id
}

func (e *EventTransmitter) Approval() error {
	return e.t.Approval(e.m)
}

func (e *EventTransmitter) Rejected() error {
	return e.t.Rejected(e.m)
}

type Transmitter interface {
	Approval(m *Message) error
	Rejected(m *Message) error
}

type states string

const (
	verify_consumer states = "verify_consumer"
	create_ticket   states = "create_ticket"
	verify_card     states = "verify_card"
	confirm_ticket  states = "confirm_ticket"
	confirm_order   states = "confirm_order"
)

type VerifyConsumer string

const (
	StartcheckVerifyConsumer VerifyConsumer = "startcheck"
	CheckedVerifyConsumer    VerifyConsumer = "checked"
	FailedVerifyConsumer     VerifyConsumer = "failed"
)

type CreateTicket string

const (
	StartcheckCreateTicket CreateTicket = "startcheck"
	CheckedCreateTicket    CreateTicket = "checked"
	FailedCreateTicket     CreateTicket = "failed"
)

type VerifyCard string

const (
	StartcheckVerifyCard VerifyCard = "startcheck"
	CheckedVerifyCard    VerifyCard = "checked"
	FailedVerifyCard     VerifyCard = "failed"
)

type ConfirmTicket string

const (
	StartConfirmTicket       ConfirmTicket = "start"
	ConfirmConfirmTicket     ConfirmTicket = "confirm"
	errRejectedConfirmTicket ConfirmTicket = "errRejected"
)

type ConfirmOrder string

const (
	StartConfirmOrder       ConfirmOrder = "start"
	ConfirmConfirmOrder     ConfirmOrder = "confirm"
	errRejectedConfirmOrder ConfirmOrder = "errRejected"
)

func (t VerifyConsumer) Is() bool {
	if t == StartcheckVerifyConsumer || t == CheckedVerifyConsumer || t == FailedVerifyConsumer {
		return true
	}
	return false
}

func (t CreateTicket) Is() bool {
	if t == StartcheckCreateTicket || t == CheckedCreateTicket || t == FailedCreateTicket {
		return true
	}
	return false
}

func (t VerifyCard) Is() bool {
	if t == StartcheckVerifyCard || t == CheckedVerifyCard || t == FailedVerifyCard {
		return true
	}
	return false
}

func (t ConfirmTicket) Is() bool {
	if t == StartConfirmTicket || t == ConfirmConfirmTicket {
		return true
	}
	return false
}

func (t ConfirmOrder) Is() bool {
	if t == StartConfirmOrder || t == ConfirmConfirmOrder {
		return true
	}
	return false
}

type Message struct {
	ID        string    `json:"id"`
	Command   string    `json:"command"`
	StepName  string    `json:"step_name"`
	Direction direction `json:"direction"`
	Retry     int       `json:"retry"`
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

func (t *VerifyConsumerTransmitter) Pending(m *Message) error {
	m.Command = string(StartcheckVerifyConsumer)
	m.StepName = "verify_consumer"
	b, _ := json.Marshal(m)
	body := &broker.Message{Body: b}
	return t.b.Publish(orchestratorRoutingKey, body)
}

func (t *VerifyConsumerTransmitter) Approval(m *Message) error {
	m.Command = string(CheckedVerifyConsumer)
	m.StepName = "verify_consumer"
	b, _ := json.Marshal(m)
	body := &broker.Message{Body: b}
	return t.b.Publish(orchestratorRoutingKey, body)
}

func (t *VerifyConsumerTransmitter) Rejected(m *Message) error {
	m.Command = string(FailedVerifyConsumer)
	m.StepName = "verify_consumer"
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
			id:    m.ID,
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
			id:    m.ID,
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

func (t *CreateTicketTransmitter) Pending(m *Message) error {
	m.Command = string(StartcheckCreateTicket)
	m.StepName = "create_ticket"
	b, _ := json.Marshal(m)
	body := &broker.Message{Body: b}
	return t.b.Publish(orchestratorRoutingKey, body)
}

func (t *CreateTicketTransmitter) Approval(m *Message) error {
	m.Command = string(CheckedCreateTicket)
	m.StepName = "create_ticket"
	b, _ := json.Marshal(m)
	body := &broker.Message{Body: b}
	return t.b.Publish(orchestratorRoutingKey, body)
}

func (t *CreateTicketTransmitter) Rejected(m *Message) error {
	m.Command = string(FailedCreateTicket)
	m.StepName = "create_ticket"
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
			id:    m.ID,
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
			id:    m.ID,
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

func (t *VerifyCardTransmitter) Pending(m *Message) error {
	m.Command = string(StartcheckVerifyCard)
	m.StepName = "verify_card"
	b, _ := json.Marshal(m)
	body := &broker.Message{Body: b}
	return t.b.Publish(orchestratorRoutingKey, body)
}

func (t *VerifyCardTransmitter) Approval(m *Message) error {
	m.Command = string(CheckedVerifyCard)
	m.StepName = "verify_card"
	b, _ := json.Marshal(m)
	body := &broker.Message{Body: b}
	return t.b.Publish(orchestratorRoutingKey, body)
}

func (t *VerifyCardTransmitter) Rejected(m *Message) error {
	m.Command = string(FailedVerifyCard)
	m.StepName = "verify_card"
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
			id:    m.ID,
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
			id:    m.ID,
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

func (t *ConfirmTicketTransmitter) Pending(m *Message) error {
	m.Command = string(StartConfirmTicket)
	m.StepName = "confirm_ticket"
	b, _ := json.Marshal(m)
	body := &broker.Message{Body: b}
	return t.b.Publish(orchestratorRoutingKey, body)
}

func (t *ConfirmTicketTransmitter) Approval(m *Message) error {
	m.Command = string(ConfirmConfirmTicket)
	m.StepName = "confirm_ticket"
	b, _ := json.Marshal(m)
	body := &broker.Message{Body: b}
	return t.b.Publish(orchestratorRoutingKey, body)
}

func (t *ConfirmTicketTransmitter) Rejected(m *Message) error {

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
			id:    m.ID,
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
			id:    m.ID,
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

func (t *ConfirmOrderTransmitter) Pending(m *Message) error {
	m.Command = string(StartConfirmOrder)
	m.StepName = "confirm_order"
	b, _ := json.Marshal(m)
	body := &broker.Message{Body: b}
	return t.b.Publish(orchestratorRoutingKey, body)
}

func (t *ConfirmOrderTransmitter) Approval(m *Message) error {
	m.Command = string(ConfirmConfirmOrder)
	m.StepName = "confirm_order"
	b, _ := json.Marshal(m)
	body := &broker.Message{Body: b}
	return t.b.Publish(orchestratorRoutingKey, body)
}

func (t *ConfirmOrderTransmitter) Rejected(m *Message) error {

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
			id:    m.ID,
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
			id:    m.ID,
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
		case StartcheckVerifyConsumer:
			o.log.Log(logger.WarnLevel, StartcheckVerifyConsumer+" is not defined for orchestrator")
			return nil
		case CheckedVerifyConsumer:
			return o.b.Publish("create_ticket.pending", m)

		case FailedVerifyConsumer:
			return o.b.Publish("verify_consumer.rejected", m)
		default:
			panic(typeOf)
		}
	case Down:
		switch typeOf {
		case StartcheckVerifyConsumer:
			o.log.Log(logger.WarnLevel, StartcheckVerifyConsumer+" is not defined for orchestrator")
			return nil
		case CheckedVerifyConsumer:
			return nil
		case FailedVerifyConsumer:
			o.log.Log(logger.ErrorLevel, "it's happened rejecting transaction which rejected")
			return errors.New("tx accident")
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
		case StartcheckCreateTicket:
			o.log.Log(logger.WarnLevel, StartcheckCreateTicket+" is not defined for orchestrator")
			return nil
		case CheckedCreateTicket:
			return o.b.Publish("verify_card.pending", m)

		case FailedCreateTicket:
			return o.b.Publish("create_ticket.rejected", m)
		default:
			panic(typeOf)
		}
	case Down:
		switch typeOf {
		case StartcheckCreateTicket:
			o.log.Log(logger.WarnLevel, StartcheckCreateTicket+" is not defined for orchestrator")
			return nil
		case CheckedCreateTicket:
			return o.b.Publish("verify_consumer.pending", m)
		case FailedCreateTicket:
			o.log.Log(logger.ErrorLevel, "it's happened rejecting transaction which rejected")
			return errors.New("tx accident")
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
		case StartcheckVerifyCard:
			o.log.Log(logger.WarnLevel, StartcheckVerifyCard+" is not defined for orchestrator")
			return nil
		case CheckedVerifyCard:
			return o.b.Publish("confirm_ticket.pending", m)

		case FailedVerifyCard:
			return o.b.Publish("verify_card.rejected", m)
		default:
			panic(typeOf)
		}
	case Down:
		switch typeOf {
		case StartcheckVerifyCard:
			o.log.Log(logger.WarnLevel, StartcheckVerifyCard+" is not defined for orchestrator")
			return nil
		case CheckedVerifyCard:
			return o.b.Publish("create_ticket.pending", m)
		case FailedVerifyCard:
			o.log.Log(logger.ErrorLevel, "it's happened rejecting transaction which rejected")
			return errors.New("tx accident")
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
	switch direction {
	case Up:
		switch typeOf {
		case StartConfirmTicket:
			o.log.Log(logger.WarnLevel, StartConfirmTicket+" is not defined for orchestrator")
			return nil
		case ConfirmConfirmTicket:
			return o.b.Publish("confirm_order.pending", m)

		default:
			panic(typeOf)
		}
	case Down:
		switch typeOf {
		case StartConfirmTicket:
			o.log.Log(logger.WarnLevel, StartConfirmTicket+" is not defined for orchestrator")
			return nil
		case ConfirmConfirmTicket:
			return o.b.Publish("verify_card.pending", m)
		case errRejectedConfirmTicket:
			o.log.Log(logger.ErrorLevel, "it's happened rejecting transaction which rejected")
			return errors.New("tx accident")
		default:
			panic(typeOf)
		}
	}
	return nil
}

func (o *Orchestrator) confirm_orderRoute(m *broker.Message, typeOf ConfirmOrder, direction direction) error {
	if !typeOf.Is() {
		return errors.New("invalid typeOf")
	}
	switch direction {
	case Up:
		switch typeOf {
		case StartConfirmOrder:
			o.log.Log(logger.WarnLevel, StartConfirmOrder+" is not defined for orchestrator")
			return nil
		case ConfirmConfirmOrder:
			return nil

		default:
			panic(typeOf)
		}
	case Down:
		switch typeOf {
		case StartConfirmOrder:
			o.log.Log(logger.WarnLevel, StartConfirmOrder+" is not defined for orchestrator")
			return nil
		case ConfirmConfirmOrder:
			return o.b.Publish("confirm_ticket.pending", m)
		case errRejectedConfirmOrder:
			o.log.Log(logger.ErrorLevel, "it's happened rejecting transaction which rejected")
			return errors.New("tx accident")
		default:
			panic(typeOf)
		}
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

//TODO will add handler if it need
//type OrchestratorHandler func(event broker.Event)
//func AddHandler(h OrchestratorHandler) {
//
//}
