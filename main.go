package main

import (
	"fmt"

	"github.com/abiosoft/ishell"
)

type Kind int

const (
	Equipment Kind = 0
	Property  Kind = 1
	Link      Kind = 2
)

func (k Kind) String() string {
	names := [...]string{
		"Equipment",
		"Property",
		"Link",
	}
	if k < Equipment || k > Link {
		return "Unknown"
	}
	return names[k]
}

type Node struct {
	ID       string
	Name     string
	Category Kind
	ParentID string
	Revision int
}

func NewNode(name string, category Kind, parentID string, revision int) *Node {
	return &Node{
		ID:       idGen.NewID(category),
		ParentID: parentID,
		Category: category,
		Name:     name,
	}
}

type Store struct {
	DB []*Node
}

func (s *Store) Save(n *Node) {
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
			node := NewNode(name, Equipment, parentID, revision)
			store.Save(node)
		},
	}
}

func addProperty() *ishell.Cmd {
	return &ishell.Cmd{
		Name: "addprop",
		Help: "add property",
		Func: func(ctx *ishell.Context) {
			numArgs := len(ctx.Args)

			name, err := GetArg(ctx, 0)
			if err != nil {
				shell.Println(err)
				return
			}

			var parentID = "Nil"
			// get optional param: parentID
			if numArgs >= 2 {
				parent, err := GetArg(ctx, 1)
				if err != nil {
					goto Execute
				}
				parentID = parent
			}

		Execute:
			node := NewNode(name, Property, parentID, 0)
			store.Save(node)
		},
	}
}

func listNodes() *ishell.Cmd {
	return &ishell.Cmd{
		Name: "nodes",
		Help: "list nodes",
		Func: func(ctx *ishell.Context) {
			for _, node := range store.DB {
				shell.Printf("%s %s %s %s\n", node.ID, node.ParentID, node.Category, node.Name)
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
	case Equipment:
		return fmt.Sprintf("E%d", curID)
	case Property:
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
	shell.SetPrompt("$ ")

	shell.Println("Topology Manager Shell. READY.")

	shell.AddCmd(addEquipment())
	shell.AddCmd(addProperty())
	shell.AddCmd(listNodes())

	shell.Run()
}
