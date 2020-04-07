
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
const(
    Up direction = iota
    Down
)

func (d direction)Is() bool{
    return d == Up || d == Down
}

type EventTransmitter struct{
    t Transmitter
    m *Message
    broker.Event
}

func (e *EventTransmitter)Retry() int{
    return e.m.Retry
}

func (e *EventTransmitter)SetPayload(b []byte){
    e.m.Payload = b
}

func (e *EventTransmitter)Payload() []byte {
    return e.m.Payload
}

func (e *EventTransmitter)ID()string{
    return e.m.ID
}

func (e *EventTransmitter)Approve() error {
    return e.t.Approve(e.m)
}

func (e *EventTransmitter)Reject() error {
    return e.t.Reject(e.m)
}

type Transmitter interface{
    Approve(m *Message) error
    Reject(m *Message) error
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
    {{if ne .Sl.Rejected ""}}{{.Sl.Rejected | camelcase}}{{.Name  | camelcase}} {{.Name  | camelcase}} = "{{.Sl.Rejected}}"
    {{else}} Failed{{.Sl.Rejected | camelcase}}{{.Name  | camelcase}} {{.Name  | camelcase}} = "Failed"{{end}}
)
{{end}}

{{range .}}
func (t {{.Name  | camelcase}}) Is () bool{
    if t == {{.Sl.Pending | camelcase}}{{.Name  | camelcase}} || t == {{.Sl.Approval | camelcase}}{{.Name  | camelcase}} {{if ne .Sl.Rejected ""}} || t == {{.Sl.Rejected | camelcase}}{{.Name  | camelcase}}{{else}} || t == Failed{{.Name  | camelcase}}{{end}}{
        return true
    }
    return false
}
{{end}}

type Message struct{
    ID          string                  `json:"id"`
    Command     string                  `json:"command"`
    StepName    string                  `json:"step_name"`
    Direction   direction               `json:"direction"`
    //The current number of the retry
    Retry       int                     `json:"retry"`
    //If it need we can add payload to message.
    Payload     []byte                  `json:"payload"`
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

func (t *{{.Name | camelcase}}Transmitter)Approve(m *Message) error {
    m.Command = string({{.Sl.Approval | camelcase}}{{.Name  | camelcase}})
    m.StepName = "{{.Name}}"
    m.Retry = 0
    b,_ :=json.Marshal(m)
    body := &broker.Message{Body:b}
    return t.b.Publish(orchestratorRoutingKey,  body)
}

func (t *{{.Name | camelcase}}Transmitter)Reject(m *Message) error {
{{if ne .T "retriable"}}
    m.Command = string({{.Sl.Rejected | camelcase}}{{.Name  | camelcase}})
    if m.Direction != Down{
        m.Command = string({{.Sl.Approval | camelcase}}{{.Name  | camelcase}})
    }
    m.Direction = Down
{{else}}
    m.Command = string(Failed{{.Sl.Rejected | camelcase}}{{.Name  | camelcase}}){{end}}
    m.StepName = "{{.Name}}"
    if m.Retry>t.MaxRetries{
        return ErrToManyRetries
    }

    m.Retry++
    b,_ :=json.Marshal(m)
    body := &broker.Message{Body:b}
    return t.b.Publish(orchestratorRoutingKey, body)
}
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

    switch steps(m.StepName){
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
    {{if ne .T "retriable"}}
    switch direction{
        case Up:
            switch typeOf{
                case {{.Sl.Approval | camelcase}}{{.Name  | camelcase}}:
                {{if .Next}}    return o.b.Publish("{{.Next.Name}}.pending", m){{else}}    return nil{{end}}
                default:
                    panic(typeOf)
            }
        case Down:
            switch typeOf{
                case {{.Sl.Approval | camelcase}}{{.Name  | camelcase}}:
                {{if .Prev}}    return o.b.Publish("{{.Prev.Name}}.rejected", m){{else}}    return nil{{end}}
                case {{if ne .Sl.Rejected ""}}{{.Sl.Rejected | camelcase}}{{.Name  | camelcase}}{{else}}Failed{{.Name  | camelcase}}{{end}}:
                    return o.b.Publish("{{.Name}}.rejected", m)
                default:
                    panic(typeOf)
            }
    }
    {{else}}
        switch typeOf{
        case {{.Sl.Approval | camelcase}}{{.Name  | camelcase}}:
        {{if .Next}}    return o.b.Publish("{{.Next.Name}}.pending", m){{else}}    return nil{{end}}
            case Failed{{.Name  | camelcase}}:
            return o.b.Publish("{{.Name}}.rejected", m)
        default:
        panic(typeOf)
        }
    {{end}}
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
