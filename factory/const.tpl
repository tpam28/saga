//Automatically generated file; DO NOT EDIT
package main

import (
    "encoding/json"
    "errors"

    "github.com/micro/go-micro/v2/broker"
)

//TODO chacnge
const orchestratorRoutingKey = "milestone.orchestrator"
var ErrMethodNotAvailable = errors.New("method not available")

type EventTransmitter struct{
    t Transmitter
    id string
    broker.Event
}

func (e *EventTransmitter)Approval() error {
    return e.t.Approval(e.id)
}

func (e *EventTransmitter)Rejected() error {
    return e.t.Rejected(e.id)
}

type Transmitter interface{
    Approval(id string) error
    Rejected(id string) error
}


type steps string
const(
{{range .}}    {{.Name}} steps = "{{.Name}}"
{{end}}
)

{{range .}}
type {{.Name  | camelcase}} string
const(
    {{.Sl.Pending | camelcase}}{{.Name  | camelcase}} {{.Name  | camelcase}} = "{{.Sl.Pending}}"
    {{.Sl.Approval | camelcase}}{{.Name  | camelcase}} {{.Name  | camelcase}} = "{{.Sl.Approval}}"
    {{if ne .Sl.Rejected ""}}{{.Sl.Rejected | camelcase}}{{.Name  | camelcase}} {{.Name  | camelcase}} = "{{.Sl.Rejected}}"{{end}}
)
{{end}}

{{range .}}
func (t {{.Name  | camelcase}}) Is () bool{
    if t == {{.Sl.Pending | camelcase}}{{.Name  | camelcase}} || t == {{.Sl.Approval | camelcase}}{{.Name  | camelcase}} {{if ne .Sl.Rejected ""}} || t == {{.Sl.Rejected | camelcase}}{{.Name  | camelcase}}{{end}}{
        return true
    }
    return false
}
{{end}}

type Message struct{
    ID          string `json:"id"`
    Command     string `json:"command"`
    StepName    string `json:"step_name"`
    //If it need we can add payload to message.
    Payload     []byte `json:'payload"`
}

func MessageByte(id string, command string) []byte {
    m := &Message{
        ID:id,
        Command:command,
    }
    b, _ :=json.Marshal(m)
    return b
}

{{range .}}
type {{.Name | camelcase}}Transmitter struct{
    b broker.Broker
}
func (t *{{.Name | camelcase}}Transmitter)Pending(id string) error {
    body := &broker.Message{Body:MessageByte(id, string({{.Sl.Pending | camelcase}}{{.Name  | camelcase}}))}
    return t.b.Publish(orchestratorRoutingKey, body)
}

func (t *{{.Name | camelcase}}Transmitter)Approval(id string) error {
    body := &broker.Message{Body:MessageByte(id, string({{.Sl.Approval | camelcase}}{{.Name  | camelcase}}))}
    return t.b.Publish(orchestratorRoutingKey,  body)
}

func (t *{{.Name | camelcase}}Transmitter)Rejected(id string) error {
{{if ne .Sl.Rejected ""}}body := &broker.Message{Body:MessageByte(id, string({{.Sl.Rejected | camelcase}}{{.Name  | camelcase}}))}
    return t.b.Publish(orchestratorRoutingKey, body)
{{else}}    return ErrMethodNotAvailable
{{end}}}

func New{{.Name | camelcase}}Transmitter(b broker.Broker) *{{.Name | camelcase}}Transmitter{
    return &{{.Name | camelcase}}Transmitter{b:b}
}

type {{.Name | camelcase}}Receiver struct{
    b broker.Broker
    t Transmitter
}

func (r *{{.Name | camelcase}}Receiver) Pending(f func(EventTransmitter) error) (broker.Subscriber, error){
    return r.b.Subscribe("{{.Name}}_approval", func(event broker.Event) error {
        m := Message{}
        err := json.Unmarshal(event.Message().Body, &m)
        if err != nil {
            panic(err)
        }
        return f(EventTransmitter{
            t: r.t,
            id:m.ID,
            Event: event,
        })
    })
}

func (r *{{.Name | camelcase}}Receiver) Rejected(f func(EventTransmitter)error) (broker.Subscriber, error){
    return r.b.Subscribe("{{.Name}}_rejected", func(event broker.Event) error {
        m := Message{}
        err := json.Unmarshal(event.Message().Body, &m)
        if err != nil {
            panic(err)
        }
        return f(EventTransmitter{
            t: r.t,
            id:m.ID,
            Event: event,
        })
    })

}

func New{{.Name | camelcase}}Receiver(b broker.Broker) *{{.Name | camelcase}}Receiver{
    return &{{.Name | camelcase}}Receiver{b:b}
}
{{end}}

type Orchestrator struct {
    b broker.Broker
    //TODO add callback for bad transaction for example: use this if we reject transaction witch has rejected.
    //TODO use micro logger interface
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
        {{range . }}case {{.Name}}:
            return o.{{.Name}}Route(e.Message(), {{.Name | camelcase}}(m.Command))
        {{end}}
    }
    return nil
}

{{range . }}
func (o *Orchestrator){{.Name}}Route(m *broker.Message, typeOf {{.Name  | camelcase}}) error {
    if  !typeOf.Is() {
        return errors.New("invalid type of")
    }
    switch typeOf{
        case {{.Sl.Pending | camelcase}}{{.Name  | camelcase}}:
            return o.b.Publish("{{.Name}}.pending", m, nil)
        case {{.Sl.Approval | camelcase}}{{.Name  | camelcase}}:
            return o.b.Publish("{{.Name}}.approval", m, nil)
        {{if ne .Sl.Rejected ""}}
        case {{.Sl.Rejected | camelcase}}{{.Name  | camelcase}}:
            return o.b.Publish("{{.Name}}.rejected", m, nil)
            {{end}}
        default:
            panic(typeOf)
    }
}
{{end}}

//TODO will add handler if it need
//type OrchestratorHandler func(event broker.Event)
//func AddHandler(h OrchestratorHandler) {
//
//}
