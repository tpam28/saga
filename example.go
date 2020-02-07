//DO NOT EDIT
package main

import (
    "encoding/json"
    "errors"

    "github.com/micro/go-micro/v2/broker"
)

type Event struct{
    broker.Event
}
//TODO chacnge
const orchestratorRoutingKey = "milestone.orchestrator"
var ErrMethodNotAvailable = errors.New("method not available")

type EventTransmitter struct{
    t Transmitter
    id string
}

func (e *EventTransmitter)Pending() error {
    return e.t.Pending(e.id)
}

func (e *EventTransmitter)Approval() error {
    return e.t.Approval(e.id)
}

func (e *EventTransmitter)Rejected() error {
    return e.t.Rejected(e.id)
}

type Transmitter interface{
    Pending(id string) error
    Approval(id string) error
    Rejected(id string) error
}


type steps string
const(

    verify_consumer steps = "verify_consumer"

    create_ticket steps = "create_ticket"

    verify_card steps = "verify_card"

    confirm_ticket steps = "confirm_ticket"

    confirm_order steps = "confirm_order"

)


type VerifyConsumer string
const(
    StartcheckVerifyConsumer VerifyConsumer = "startcheck"
    CheckedVerifyConsumer VerifyConsumer = "checked"
    FailedVerifyConsumer VerifyConsumer = "failed"
)

type CreateTicket string
const(
    StartcheckCreateTicket CreateTicket = "startcheck"
    CheckedCreateTicket CreateTicket = "checked"
    FailedCreateTicket CreateTicket = "failed"
)

type VerifyCard string
const(
    StartcheckVerifyCard VerifyCard = "startcheck"
    CheckedVerifyCard VerifyCard = "checked"
    FailedVerifyCard VerifyCard = "failed"
)

type ConfirmTicket string
const(
    StartConfirmTicket ConfirmTicket = "start"
    ConfirmConfirmTicket ConfirmTicket = "confirm"
    
)

type ConfirmOrder string
const(
    StartConfirmOrder ConfirmOrder = "start"
    ConfirmConfirmOrder ConfirmOrder = "confirm"
    
)



func (t VerifyConsumer) Is () bool{
    if t == StartcheckVerifyConsumer || t == CheckedVerifyConsumer  || t == FailedVerifyConsumer{
        return true
    }
    return false
}

func (t CreateTicket) Is () bool{
    if t == StartcheckCreateTicket || t == CheckedCreateTicket  || t == FailedCreateTicket{
        return true
    }
    return false
}

func (t VerifyCard) Is () bool{
    if t == StartcheckVerifyCard || t == CheckedVerifyCard  || t == FailedVerifyCard{
        return true
    }
    return false
}

func (t ConfirmTicket) Is () bool{
    if t == StartConfirmTicket || t == ConfirmConfirmTicket {
        return true
    }
    return false
}

func (t ConfirmOrder) Is () bool{
    if t == StartConfirmOrder || t == ConfirmConfirmOrder {
        return true
    }
    return false
}


type Message struct{
    ID          string `json:"id"`
    Command     string `json:"command"`
    StepName    string `json:"step_name"`
}

func MessageByte(id string, command string) []byte {
    m := &Message{
        ID:id,
        Command:command,
    }
    b, _ :=json.Marshal(m)
    return b
}


type VerifyConsumerTransmitter struct{
    b broker.Broker
}
func (t *VerifyConsumerTransmitter)Pending(id string) error {
    body := &broker.Message{Body:MessageByte(id, string(StartcheckVerifyConsumer))}
    return t.b.Publish(orchestratorRoutingKey, body)
}

func (t *VerifyConsumerTransmitter)Approval(id string) error {
    body := &broker.Message{Body:MessageByte(id, string(CheckedVerifyConsumer))}
    return t.b.Publish(orchestratorRoutingKey,  body)
}

func (t *VerifyConsumerTransmitter)Rejected(id string) error {
body := &broker.Message{Body:MessageByte(id, string(FailedVerifyConsumer))}
    return t.b.Publish(orchestratorRoutingKey, body)
}

func NewVerifyConsumerTransmitter(b broker.Broker) *VerifyConsumerTransmitter{
    return &VerifyConsumerTransmitter{b:b}
}

type VerifyConsumerReceiver struct{
    b broker.Broker
}

func (r *VerifyConsumerReceiver(f func(broker.Event) {
}
func NewVerifyConsumerReceiver(b broker.Broker) *VerifyConsumerReceiver{
    return &VerifyConsumerReceiver{b:b}
}

type CreateTicketTransmitter struct{
    b broker.Broker
}
func (t *CreateTicketTransmitter)Pending(id string) error {
    body := &broker.Message{Body:MessageByte(id, string(StartcheckCreateTicket))}
    return t.b.Publish(orchestratorRoutingKey, body)
}

func (t *CreateTicketTransmitter)Approval(id string) error {
    body := &broker.Message{Body:MessageByte(id, string(CheckedCreateTicket))}
    return t.b.Publish(orchestratorRoutingKey,  body)
}

func (t *CreateTicketTransmitter)Rejected(id string) error {
body := &broker.Message{Body:MessageByte(id, string(FailedCreateTicket))}
    return t.b.Publish(orchestratorRoutingKey, body)
}

func NewCreateTicketTransmitter(b broker.Broker) *CreateTicketTransmitter{
    return &CreateTicketTransmitter{b:b}
}

type CreateTicketReceiver struct{
    b broker.Broker
}

func (r *CreateTicketReceiver(f func(broker.Event) {
}
func NewCreateTicketReceiver(b broker.Broker) *CreateTicketReceiver{
    return &CreateTicketReceiver{b:b}
}

type VerifyCardTransmitter struct{
    b broker.Broker
}
func (t *VerifyCardTransmitter)Pending(id string) error {
    body := &broker.Message{Body:MessageByte(id, string(StartcheckVerifyCard))}
    return t.b.Publish(orchestratorRoutingKey, body)
}

func (t *VerifyCardTransmitter)Approval(id string) error {
    body := &broker.Message{Body:MessageByte(id, string(CheckedVerifyCard))}
    return t.b.Publish(orchestratorRoutingKey,  body)
}

func (t *VerifyCardTransmitter)Rejected(id string) error {
body := &broker.Message{Body:MessageByte(id, string(FailedVerifyCard))}
    return t.b.Publish(orchestratorRoutingKey, body)
}

func NewVerifyCardTransmitter(b broker.Broker) *VerifyCardTransmitter{
    return &VerifyCardTransmitter{b:b}
}

type VerifyCardReceiver struct{
    b broker.Broker
}

func (r *VerifyCardReceiver(f func(broker.Event) {
}
func NewVerifyCardReceiver(b broker.Broker) *VerifyCardReceiver{
    return &VerifyCardReceiver{b:b}
}

type ConfirmTicketTransmitter struct{
    b broker.Broker
}
func (t *ConfirmTicketTransmitter)Pending(id string) error {
    body := &broker.Message{Body:MessageByte(id, string(StartConfirmTicket))}
    return t.b.Publish(orchestratorRoutingKey, body)
}

func (t *ConfirmTicketTransmitter)Approval(id string) error {
    body := &broker.Message{Body:MessageByte(id, string(ConfirmConfirmTicket))}
    return t.b.Publish(orchestratorRoutingKey,  body)
}

func (t *ConfirmTicketTransmitter)Rejected(id string) error {
    return ErrMethodNotAvailable
}

func NewConfirmTicketTransmitter(b broker.Broker) *ConfirmTicketTransmitter{
    return &ConfirmTicketTransmitter{b:b}
}

type ConfirmTicketReceiver struct{
    b broker.Broker
}

func (r *ConfirmTicketReceiver(f func(broker.Event) {
}
func NewConfirmTicketReceiver(b broker.Broker) *ConfirmTicketReceiver{
    return &ConfirmTicketReceiver{b:b}
}

type ConfirmOrderTransmitter struct{
    b broker.Broker
}
func (t *ConfirmOrderTransmitter)Pending(id string) error {
    body := &broker.Message{Body:MessageByte(id, string(StartConfirmOrder))}
    return t.b.Publish(orchestratorRoutingKey, body)
}

func (t *ConfirmOrderTransmitter)Approval(id string) error {
    body := &broker.Message{Body:MessageByte(id, string(ConfirmConfirmOrder))}
    return t.b.Publish(orchestratorRoutingKey,  body)
}

func (t *ConfirmOrderTransmitter)Rejected(id string) error {
    return ErrMethodNotAvailable
}

func NewConfirmOrderTransmitter(b broker.Broker) *ConfirmOrderTransmitter{
    return &ConfirmOrderTransmitter{b:b}
}

type ConfirmOrderReceiver struct{
    b broker.Broker
}

func (r *ConfirmOrderReceiver(f func(broker.Event) {
}
func NewConfirmOrderReceiver(b broker.Broker) *ConfirmOrderReceiver{
    return &ConfirmOrderReceiver{b:b}
}


type Orchestrator struct {
    b broker.Broker
}

func (o *Orchestrator) Do(options ...broker.SubscribeOption) (broker.Subscriber, error) {
    return o.b.Subscribe(orchestratorRoutingKey, o.handler, options...)
}

func (o *Orchestrator) handler(e broker.Event) error {
    m := Message{}
    err := json.Unmarshal(e.Message().Body,&m)
    //it mustn't ever happens
    if err != nil{
        panic(err)
    }

    switch steps(m.StepName){
        case verify_consumer:
            return o.verify_consumerRoute(e.Message(), VerifyConsumer(m.Command))
        case create_ticket:
            return o.create_ticketRoute(e.Message(), CreateTicket(m.Command))
        case verify_card:
            return o.verify_cardRoute(e.Message(), VerifyCard(m.Command))
        case confirm_ticket:
            return o.confirm_ticketRoute(e.Message(), ConfirmTicket(m.Command))
        case confirm_order:
            return o.confirm_orderRoute(e.Message(), ConfirmOrder(m.Command))
        
    }
    return nil
}


func (o *Orchestrator)verify_consumerRoute(m *broker.Message, typeOf VerifyConsumer) error {
    if  !typeOf.Is() {
        return errors.New("invalid type of")
    }
    switch typeOf{
        case StartcheckVerifyConsumer:
            return o.b.Publish("verify_consumer.pending", m, nil)
        case CheckedVerifyConsumer:
            return o.b.Publish("verify_consumer.approval", m, nil)
        
        case FailedVerifyConsumer:
            return o.b.Publish("verify_consumer.rejected", m, nil)
            
        default:
            panic(typeOf)
    }
}

func (o *Orchestrator)create_ticketRoute(m *broker.Message, typeOf CreateTicket) error {
    if  !typeOf.Is() {
        return errors.New("invalid type of")
    }
    switch typeOf{
        case StartcheckCreateTicket:
            return o.b.Publish("create_ticket.pending", m, nil)
        case CheckedCreateTicket:
            return o.b.Publish("create_ticket.approval", m, nil)
        
        case FailedCreateTicket:
            return o.b.Publish("create_ticket.rejected", m, nil)
            
        default:
            panic(typeOf)
    }
}

func (o *Orchestrator)verify_cardRoute(m *broker.Message, typeOf VerifyCard) error {
    if  !typeOf.Is() {
        return errors.New("invalid type of")
    }
    switch typeOf{
        case StartcheckVerifyCard:
            return o.b.Publish("verify_card.pending", m, nil)
        case CheckedVerifyCard:
            return o.b.Publish("verify_card.approval", m, nil)
        
        case FailedVerifyCard:
            return o.b.Publish("verify_card.rejected", m, nil)
            
        default:
            panic(typeOf)
    }
}

func (o *Orchestrator)confirm_ticketRoute(m *broker.Message, typeOf ConfirmTicket) error {
    if  !typeOf.Is() {
        return errors.New("invalid type of")
    }
    switch typeOf{
        case StartConfirmTicket:
            return o.b.Publish("confirm_ticket.pending", m, nil)
        case ConfirmConfirmTicket:
            return o.b.Publish("confirm_ticket.approval", m, nil)
        
        default:
            panic(typeOf)
    }
}

func (o *Orchestrator)confirm_orderRoute(m *broker.Message, typeOf ConfirmOrder) error {
    if  !typeOf.Is() {
        return errors.New("invalid type of")
    }
    switch typeOf{
        case StartConfirmOrder:
            return o.b.Publish("confirm_order.pending", m, nil)
        case ConfirmConfirmOrder:
            return o.b.Publish("confirm_order.approval", m, nil)
        
        default:
            panic(typeOf)
    }
}


//TODO will add handler if it need
//type OrchestratorHandler func(event broker.Event)
//func AddHandler(h OrchestratorHandler) {
//
//}
