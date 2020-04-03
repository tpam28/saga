
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
const(
    Up direction = iota
    Down
)

func (d direction)Is() bool{
    return d == Up || d == Down
}

type EventTransmitter struct{
    t Transmitter
    id string
    m *Message
    broker.Event
}

func (e *EventTransmitter)ID()string{
    return e.id
}

func (e *EventTransmitter)Approval() error {
    return e.t.Approval(e.m)
}

func (e *EventTransmitter)Rejected() error {
    return e.t.Rejected(e.m)
}

type Transmitter interface{
    Approval(m *Message) error
    Rejected(m *Message) error
}


type states string
const(
{{range .}}    {{.Name}} states = "{{.Name}}"
{{end}}
)

{{range .}}
type {{.Name  | camelcase}} string
const(
    {{.Sl.Pending | camelcase}}{{.Name  | camelcase}} {{.Name  | camelcase}} = "{{.Sl.Pending}}"
    {{.Sl.Approval | camelcase}}{{.Name  | camelcase}} {{.Name  | camelcase}} = "{{.Sl.Approval}}"
    {{if ne .Sl.Rejected ""}}{{.Sl.Rejected | camelcase}}{{.Name  | camelcase}} {{.Name  | camelcase}} = "{{.Sl.Rejected}}"{{else}}errRejected{{.Name  | camelcase}} {{.Name  | camelcase}} = "errRejected"{{end}}
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
    ID          string    `json:"id"`
    Command     string    `json:"command"`
    StateName    string    `json:"state_name"`
    Direction   direction `json:"direction"`
    Retry       int       `json:"retry"`
    //If it need we can add payload to message.
    Payload     []byte `json:"payload"`
}

func NewMessage(id string) *Message {
    m := &Message{
        ID:id,
        Direction:Up,
    }
    return m
}

{{range .}}
type {{.Name | camelcase}}Transmitter struct{
    b broker.Broker
    MaxRetries int
}
func (t *{{.Name | camelcase}}Transmitter)Pending(m *Message) error {
    m.Command = string({{.Sl.Pending | camelcase}}{{.Name  | camelcase}})
    m.StateName = "{{.Name}}"
    b,_ :=json.Marshal(m)
    body := &broker.Message{Body:b}
    return t.b.Publish(orchestratorRoutingKey, body)
}

func (t *{{.Name | camelcase}}Transmitter)Approval(m *Message) error {
    m.Command = string({{.Sl.Approval | camelcase}}{{.Name  | camelcase}})
    m.StateName = "{{.Name}}"
    b,_ :=json.Marshal(m)
    body := &broker.Message{Body:b}
    return t.b.Publish(orchestratorRoutingKey,  body)
}

func (t *{{.Name | camelcase}}Transmitter)Rejected(m *Message) error {
{{if ne .Sl.Rejected ""}}    m.Command = string({{.Sl.Rejected | camelcase}}{{.Name  | camelcase}})
    m.StateName = "{{.Name}}"
    b,_ :=json.Marshal(m)
    body := &broker.Message{Body:b}
    return t.b.Publish(orchestratorRoutingKey, body)
{{else}}
    if m.Retry>t.MaxRetries{
        return ErrToManyRetries
    }
    m.Retry++
    b,_ :=json.Marshal(m)
    body := &broker.Message{Body:b}
    return t.b.Publish(orchestratorRoutingKey, body)
{{end}}}

func New{{.Name | camelcase}}Transmitter(b broker.Broker) *{{.Name | camelcase}}Transmitter{
    return &{{.Name | camelcase}}Transmitter{b:b, MaxRetries:10}
}

type {{.Name | camelcase}}Receiver struct{
    b broker.Broker
    t Transmitter
}

func (r *{{.Name | camelcase}}Receiver) Pending(f func(*EventTransmitter) error) (broker.Subscriber, error){
    return r.b.Subscribe("{{.Name}}.pending", func(event broker.Event) error {
        m := Message{}
        err := json.Unmarshal(event.Message().Body, &m)
        if err != nil {
            panic(err)
        }
        return f(&EventTransmitter{
            t: r.t,
            id:m.ID,
            m: &m,
            Event: event,
        })
    })
}

func (r *{{.Name | camelcase}}Receiver) Rejected(f func(*EventTransmitter)error) (broker.Subscriber, error){
    return r.b.Subscribe("{{.Name}}.rejected", func(event broker.Event) error {
        m := Message{}
        err := json.Unmarshal(event.Message().Body, &m)
        if err != nil {
            panic(err)
        }
        return f(&EventTransmitter{
            t: r.t,
            id:m.ID,
            m: &m,
            Event: event,
        })
    })

}

func New{{.Name | camelcase}}Receiver(b broker.Broker) *{{.Name | camelcase}}Receiver{
    return &{{.Name | camelcase}}Receiver{
        b:b,
        t:New{{.Name | camelcase}}Transmitter(b),
        }
}
{{end}}

type Orchestrator struct {
    b broker.Broker
    log  logger.Logger
    //TODO add callback for bad transaction for example: use this if we reject transaction witch has rejected.
    //TODO add storage
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

    switch states(m.StateName){
        {{range . }}case {{.Name}}:
            return o.{{.Name}}Route(e.Message(), {{.Name | camelcase}}(m.Command), m.Direction)
        {{end}}
    }
    return nil
}

{{range . }}
func (o *Orchestrator){{.Name}}Route(m *broker.Message, typeOf {{.Name  | camelcase}}, direction direction) error {
    if  !typeOf.Is() {
        return errors.New("invalid typeOf")
    }
    switch direction{
        case Up:
            switch typeOf{
                case {{.Sl.Pending | camelcase}}{{.Name  | camelcase}}:
                    o.log.Log(logger.WarnLevel,{{.Sl.Pending | camelcase}}{{.Name  | camelcase}}+" is not defined for orchestrator")
                    return nil
                case {{.Sl.Approval | camelcase}}{{.Name  | camelcase}}:
                {{if .Next}}    return o.b.Publish("{{.Next.Name}}.pending", m){{else}}    return nil{{end}}
                {{if ne .Sl.Rejected ""}}
                case {{.Sl.Rejected | camelcase}}{{.Name  | camelcase}}:
                    return o.b.Publish("{{.Name}}.rejected", m){{end}}
                default:
                    panic(typeOf)
            }
        case Down:
            switch typeOf{
                case {{.Sl.Pending | camelcase}}{{.Name  | camelcase}}:
                    o.log.Log(logger.WarnLevel,{{.Sl.Pending | camelcase}}{{.Name  | camelcase}}+" is not defined for orchestrator")
                    return nil
                case {{.Sl.Approval | camelcase}}{{.Name  | camelcase}}:
                {{if .Prev}}    return o.b.Publish("{{.Prev.Name}}.pending", m){{else}}    return nil{{end}}
                case {{if ne .Sl.Rejected ""}}{{.Sl.Rejected | camelcase}}{{.Name  | camelcase}}{{else}}errRejected{{.Name  | camelcase}}{{end}}:
                    o.log.Log(logger.ErrorLevel,"it's happened rejecting transaction which rejected")
                    return errors.New("tx accident")
                default:
                    panic(typeOf)
            }
    }
    return nil
}
{{end}}


func NewOrchestrator(b broker.Broker, log logger.Logger) *Orchestrator{
    if log == nil{
        log = logger.NewLogger()
    }
    return &Orchestrator{
        b: b,
        log:log,
    }
}

//TODO will add handler if it need
//type OrchestratorHandler func(event broker.Event)
//func AddHandler(h OrchestratorHandler) {
//
//}
