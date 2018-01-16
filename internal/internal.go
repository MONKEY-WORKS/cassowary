package internal

type SymbolType int

const (
	Invalid SymbolType = iota
	External
	Slack
	Error
	Dummy
)

type Symbol struct {
	Type SymbolType
}

type Tag struct {
	Marker *Symbol
	Other  *Symbol
}

func FromTag(tag *Tag) *Tag {
	return &Tag{
		Marker: tag.Marker,
		Other:  tag.Other,
	}
}

type Row struct {
	Cells    map[*Symbol]float64
	Constant float64
}

func NewRow(c float64) *Row {
	return &Row{
		Cells:    make(map[*Symbol]float64),
		Constant: c,
	}
}

func (row *Row) SolveForSymbol(symbol *Symbol) {
	if _, ok := row.Cells[symbol]; !ok {
		panic("Symbol not contained by Row")
	}

	coefficient := -1.0 / row.Cells[symbol]
	delete(row.Cells, symbol)
	row.Constant *= coefficient
	for key, value := range row.Cells {
		row.Cells[key] = value * coefficient
	}
}

func (row *Row) Substitute(symbol *Symbol, secRow *Row) {
	coefficient, ok := row.Cells[symbol]

	if !ok {
		return
	}

	delete(row.Cells, symbol)
	row.InsertRow(secRow, coefficient)
}

func (row *Row) InsertRow(secRow *Row, coefficient float64) {
	row.Constant += secRow.Constant * coefficient
	for symbol, v := range secRow.Cells {
		row.InsertSymbol(symbol, v*coefficient)
	}
}

func (row *Row) ReverseSign() {
	row.Constant = -row.Constant
	for key, val := range row.Cells {
		row.Cells[key] = -val
	}
}

func (row *Row) InsertSymbol(symbol *Symbol, coefficient float64) {
	var val float64 = 0
	if oldVal, ok := row.Cells[symbol]; ok {
		val = oldVal
	}

	val += coefficient

	if IsNearZero(val) {
		delete(row.Cells, symbol)
	} else {
		row.Cells[symbol] = val
	}
}

func (row *Row) Add(val float64) float64 {
	row.Constant += val

	return row.Constant
}

func (row *Row) SolveForSymbols(lhs *Symbol, rhs *Symbol) {
	row.InsertSymbol(lhs, -1)
	row.SolveForSymbol(rhs)
}

func (row *Row) CoefficientForSymbol(symbol *Symbol) float64 {
	if val, ok := row.Cells[symbol]; ok {
		return val
	}
	return 0
}

func IsNearZero(value float64) bool {
	const epsilon = 1.0e-8
	if value < 0 {
		return -value < epsilon
	} else {
		return value < epsilon
	}
}

func AnyPivotableSymbol(row *Row) *Symbol {
	for symbol := range row.Cells {
		if symbol.Type == Slack || symbol.Type == Error {
			return symbol
		}
	}
	return &Symbol{Invalid}
}

func CheckIfAllDummiesInRow(row *Row) bool {
	for symbol := range row.Cells {
		if symbol.Type != Dummy {
			return false
		}
	}
	return true
}

func CopyRow(srcRow *Row) *Row {
	result := NewRow(srcRow.Constant)
	for key, val := range srcRow.Cells {
		result.Cells[key] = val
	}
	return result
}
