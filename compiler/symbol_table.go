package compiler

type SymbolScope string

const (
	GlobalScope  SymbolScope = "GLOBAL"
	LocalScope   SymbolScope = "LOCAL"
	BuiltinScope SymbolScope = "BUILTIN"
)

type Symbol struct {
	Name  string
	Scope SymbolScope
	Index int
}

type SymbolTable struct {
	Outer          *SymbolTable
	store          map[string]Symbol
	numDefinitions int
}

func NewSymbolTable() *SymbolTable {
	return &SymbolTable{
		store:          make(map[string]Symbol),
		numDefinitions: 0,
	}
}

func NewEnclosedSymbolTable(outer *SymbolTable) *SymbolTable {
	return &SymbolTable{
		Outer:          outer,
		store:          make(map[string]Symbol),
		numDefinitions: 0,
	}
}

func (st *SymbolTable) Define(name string) Symbol {
	s := Symbol{
		Name:  name,
		Index: st.numDefinitions,
	}
	if st.Outer == nil {
		s.Scope = GlobalScope
	} else {
		s.Scope = LocalScope
	}
	st.store[name] = s
	st.numDefinitions++
	return s
}

func (st *SymbolTable) DefineBuiltin(index int, name string) Symbol {
	symbol := Symbol{Name: name, Index: index, Scope: BuiltinScope}
	st.store[name] = symbol
	return symbol
}

func (st *SymbolTable) Resolve(name string) (Symbol, bool) {
	s, ok := st.store[name]
	if !ok && st.Outer != nil {
		s, ok = st.Outer.Resolve(name)
	}
	return s, ok
}
