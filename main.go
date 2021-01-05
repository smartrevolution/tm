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

type Event struct {
	ID       string
	Category Kind
	ParentID string
	Revision int
	Payload  string
}

func NewEvent(category Kind, parentID string, revision int, payload string) *Event {
	return &Event{
		ID:       idGen.NewID(category),
		ParentID: parentID,
		Category: category,
		Payload:  payload,
	}
}

type Store struct {
	DB []*Event
}

func (s *Store) Save(n *Event) {
	s.DB = append(s.DB, n)
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

			// get mandatory param: payload
			payload, err := GetArg(ctx, 0)
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
			event := NewEvent(AddEquipment, parentID, revision, fmt.Sprintf("{name: %s}", payload))
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
			event := NewEvent(AddProperty, parentID, 0, fmt.Sprintf("{%s: %s}", key, value))
			store.Save(event)
		},
	}
}

func listEvents() *ishell.Cmd {
	return &ishell.Cmd{
		Name: "events",
		Help: "list events",
		Func: func(ctx *ishell.Context) {
			for _, event := range store.DB {
				shell.Printf("%s %s %s %s\n", event.Category, event.ID, event.ParentID, event.Payload)
			}
		},
	}
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
	shell = ishell.New()
	store = NewStore()
	idGen = IdGen{}
)

func main() {
	shell.Println("Topology Manager Shell. READY.")

	shell.AddCmd(addEquipment())
	shell.AddCmd(addProperty())
	shell.AddCmd(listEvents())

	shell.Run()
}
