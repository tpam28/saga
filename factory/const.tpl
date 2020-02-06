//DO NOT EDIT
package main

import (
    "encoding/json"

    "github.com/micro/go-micro/v2/broker"
)

//TODO chacnge
const orkestratorRoutingKey = "milestone.orkestrator"

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
    id      string `json:"id"`
    command string `json:"command"`
}
func MessageByte(id string, command string) []byte {
    m := &Message{
        id:id,
        command:command,
    }
    b, _ :=json.Marshal(m)
    return b
}

{{range .}}
type {{.Name | camelcase}}Transmitter struct{
    b broker.Broker
}
func (t *{{.Name | camelcase}}Transmitter){{.Sl.Pending | camelcase}}(id string) error{
    body := &broker.Message{Body:MessageByte(id, string({{.Sl.Pending | camelcase}}{{.Name  | camelcase}}))}
    return t.b.Publish(orkestratorRoutingKey, body)
}

func (t *{{.Name | camelcase}}Transmitter){{.Sl.Approval | camelcase}}(id string) error{
    body := &broker.Message{Body:MessageByte(id, string({{.Sl.Approval | camelcase}}{{.Name  | camelcase}}))}
    return t.b.Publish(orkestratorRoutingKey,  body)
}
{{if ne .Sl.Rejected ""}}
func (t *{{.Name | camelcase}}Transmitter){{.Sl.Rejected | camelcase}}(id string) error{
    body := &broker.Message{Body:MessageByte(id, string({{.Sl.Rejected | camelcase}}{{.Name  | camelcase}}))}
    return t.b.Publish(orkestratorRoutingKey, body)
}
{{end}}
func New{{.Name | camelcase}}Transmitter(b broker.Broker) *{{.Name | camelcase}}Transmitter{
    return &{{.Name | camelcase}}Transmitter{b:b}
}

type {{.Name | camelcase}}Receiver struct{
    b broker.Broker
}

func New{{.Name | camelcase}}Receiver(b broker.Broker) *{{.Name | camelcase}}Receiver{
    return &{{.Name | camelcase}}Receiver{b:b}
}
{{end}}