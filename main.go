package main

import (
	"fmt"

	"github.com/abiosoft/ishell"
)

type Kind int

const (
	AddEquipment  Kind = 0
	AddProperty   Kind = 1
	LinkEquipment Kind = 2
)

func (k Kind) String() string {
	names := [...]string{
		"AddEquipment",
		"AddProperty",
		"LinkEquipment",
	}
	if k < AddEquipment || k > LinkEquipment {
		return "Unknown"
	}
	return names[k]
}

type KeyValue map[string]string

type Event struct {
	ID        string
	Category  Kind
	ParentID  string
	Revision  int
	Payload   KeyValue
	Timestamp int
}

func NewEvent(category Kind, parentID string, revision int, payload KeyValue) *Event {
	event := Event{
		ID:        idGen.NewID(category),
		ParentID:  parentID,
		Category:  category,
		Payload:   payload,
		Timestamp: timestamp,
	}
	timestamp += 1
	return &event
}

type Store struct {
	DB []*Event
}

func (s *Store) Save(n *Event) {
	s.DB = append(s.DB, n)
	objects = s.Build()
}

type Object struct {
	ID         string
	Payload    KeyValue
	Properties KeyValue
	Children   []*Object
}

func (s *Store) Build() []*Object {
	var rootEvents []*Event
	for _, evt := range s.DB {
		if evt.ParentID == "" {
			rootEvents = append(rootEvents, evt)
		}
	}

	var objects []*Object
	for _, evt := range rootEvents {
		var obj = Object{
			ID:         evt.ID,
			Payload:    make(KeyValue),
			Properties: make(KeyValue),
		}

		// add the Payload
		for k, v := range evt.Payload {
			obj.Payload[k] = v
		}

		// execute addprop events
		for _, evt := range s.DB {
			if evt.Category == AddProperty && evt.ParentID == obj.ID {
				for k, v := range evt.Payload {
					obj.Properties[k] = v
				}
			}
		}

		objects = append(objects, &obj)
	}
	return objects
}

func NewStore() *Store {
	return &Store{}
}

func GetArg(ctx *ishell.Context, index int) (string, error) {
	if len(ctx.Args) == 0 {
		return "", fmt.Errorf("no arguments")
	}
	if len(ctx.Args)-1 < index {
		return "", fmt.Errorf("out of bounds")
	}
	return ctx.Args[index], nil
}

func addEquipment() *ishell.Cmd {
	return &ishell.Cmd{
		Name: "add",
		Help: "add <equipment-name> [equipment-parent-id]",
		Func: func(ctx *ishell.Context) {
			numArgs := len(ctx.Args)

			// get mandatory param: name
			name, err := GetArg(ctx, 0)
			if err != nil {
				shell.Println(err)
				return
			}

			var parentID string = "Nil"
			var revision int

			// get optional param: parentID
			if numArgs >= 2 {
				parent, err := GetArg(ctx, 1)
				if err != nil {
					goto Execute
				}
				parentID = parent
			}

		Execute:
			payload := NewKeyValue("Name", name)
			event := NewEvent(AddEquipment, parentID, revision, payload)
			store.Save(event)
		},
	}
}

func addProperty() *ishell.Cmd {
	return &ishell.Cmd{
		Name: "addprop",
		Help: "add property",
		Func: func(ctx *ishell.Context) {
			numArgs := len(ctx.Args)

			// get mandatory param: payload-key
			key, err := GetArg(ctx, 0)
			if err != nil {
				shell.Println(err)
				return
			}

			// get mandatory param: payload-value
			value, err := GetArg(ctx, 1)
			if err != nil {
				shell.Println(err)
				return
			}

			var parentID = "Nil"
			// get optional param: parentID
			if numArgs >= 3 {
				parent, err := GetArg(ctx, 2)
				if err != nil {
					goto Execute
				}
				parentID = parent
			}

		Execute:
			payload := NewKeyValue(key, value)
			event := NewEvent(AddProperty, parentID, 0, payload)
			store.Save(event)
		},
	}
}

func list() *ishell.Cmd {
	return &ishell.Cmd{
		Name: "list",
		Help: "list events and resulting objects",
		Func: func(ctx *ishell.Context) {
			for _, evt := range store.DB {
				shell.Printf("%d %s %s %s %s\n", evt.Timestamp, evt.Category, evt.ID, evt.ParentID, evt.Payload)
			}
			for _, obj := range objects {
				shell.Printf("%+v \n", obj)
			}
		},
	}
}

func NewKeyValue(key string, value string) map[string]string {
	return map[string]string{
		key: value,
	}

}

func loadEvents() int {
	events := []Event{
		Event{"E1", AddEquipment, "", 0, NewKeyValue("Name", "Laptop"), 1},
		Event{"P2", AddProperty, "E1", 0, NewKeyValue("Manufacturer", "Apple"), 2},
		Event{"E3", AddEquipment, "E1", 0, NewKeyValue("Name", "Keyboard"), 3},
		Event{"P4", AddProperty, "E2", 0, NewKeyValue("Manufacturer", "Apple"), 4},
		Event{"E5", AddEquipment, "", 0, NewKeyValue("Name", "Mouse"), 3},
		Event{"P6", AddProperty, "E5", 0, NewKeyValue("Manufacturer", "Razer"), 4},
	}

	for _, event := range events {
		event := event
		store.Save(&event)
	}

	return len(events)
}

type IdGen struct {
	nextID int
}

func (g *IdGen) NewID(category Kind) string {
	curID := g.nextID
	g.nextID += 1

	switch category {
	case AddEquipment:
		return fmt.Sprintf("E%d", curID)
	case AddProperty:
		return fmt.Sprintf("P%d", curID)
	default:
		return fmt.Sprintf("X%d", curID)
	}
}

var (
	shell     = ishell.New()
	store     = NewStore()
	objects   []*Object
	idGen     = IdGen{}
	timestamp int
)

func main() {
	numEvents := loadEvents()
	shell.Println("Topology Manager Shell. READY.")
	shell.Printf("%d events loaded. READY.\n", numEvents)

	shell.AddCmd(addEquipment())
	shell.AddCmd(addProperty())
	shell.AddCmd(list())

	shell.Run()
}
