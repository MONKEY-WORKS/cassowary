package cassowary

import (
	"container/list"

	"math"

	"github.com/monkey-works/cassowary/internal"
	"github.com/pkg/errors"
)

type editInfo struct {
	tag        *internal.Tag
	constraint *Constraint
	constant   float64
}

type Solver struct {
	constraints    map[*Constraint]*internal.Tag
	rows           map[*internal.Symbol]*internal.Row
	variables      map[*Variable]*internal.Symbol
	edits          map[*Variable]*editInfo
	objective      *internal.Row
	infeasibleRows *list.List
	artificial     *internal.Row
}

func NewSolver() *Solver {
	return &Solver{
		constraints:    make(map[*Constraint]*internal.Tag),
		rows:           make(map[*internal.Symbol]*internal.Row),
		variables:      make(map[*Variable]*internal.Symbol),
		edits:          make(map[*Variable]*editInfo),
		objective:      internal.NewRow(0.0),
		infeasibleRows: list.New(),
		artificial:     internal.NewRow(0.0),
	}
}

func (s *Solver) AddConstraint(constraint *Constraint) error {
	if _, ok := s.constraints[constraint]; ok {
		return errors.New("duplicate")
	}

	tag := &internal.Tag{
		Marker: &internal.Symbol{internal.Invalid},
		Other:  &internal.Symbol{internal.Invalid},
	}

	row := s.createRow(constraint, tag)

	subject := s.chooseSubjectForRow(row, tag)

	if subject.Type == internal.Invalid && internal.CheckIfAllDummiesInRow(row) {
		if !internal.IsNearZero(row.Constant) {
			return errors.New("Unsatisfiable")
		} else {
			subject = tag.Marker
		}
	}

	if subject.Type == internal.Invalid {
		if ok, err := s.addWithArtificalVariableOnRow(row); !ok {
			if err != nil {
				return err
			}
			return errors.New("Unsatisfiable")
		}
	} else {
		row.SolveForSymbol(subject)
		s.substitute(subject, row)
		s.rows[subject] = row
	}

	s.constraints[constraint] = tag

	return s.optimizeObjectiveRow(s.objective)
}

func (s *Solver) AddConstraints(constraints ...*Constraint) error {
	applier := s.AddConstraint
	undoer := s.RemoveConstraint

	return s.bulkEdit(constraints, applier, undoer)
}

func (s *Solver) RemoveConstraints(constraints ...*Constraint) error {
	applier := s.RemoveConstraint
	undoer := s.AddConstraint

	return s.bulkEdit(constraints, applier, undoer)
}

func (s *Solver) bulkEdit(constraints []*Constraint, applier, undoer func(*Constraint) error) error {
	needsCleanup := false

	var result error

	for _, constraint := range constraints {
		result = applier(constraint)
		if result == nil {
			defer func(c *Constraint) {
				if needsCleanup {
					undoer(c)
				}
			}(constraint)
		} else {
			needsCleanup = true
			break
		}
	}

	return result
}

func (s *Solver) RemoveConstraint(constraint *Constraint) error {
	tag, ok := s.constraints[constraint]
	if !ok {
		return errors.New("Unknown constraint")
	}

	tag = internal.FromTag(tag)
	delete(s.constraints, constraint)

	s.removeConstraintEffects(constraint, tag)

	row, ok := s.rows[tag.Marker]
	if ok {
		delete(s.rows, tag.Marker)
	} else {
		leaving := s.leavingSymbolForMarkerSymbol(tag.Marker)
		if leaving == nil {
			panic("Nil symbol")
		}

		row = s.rows[leaving]
		delete(s.rows, leaving)

		row.SolveForSymbols(leaving, tag.Marker)
		s.substitute(tag.Marker, row)
	}

	return s.optimizeObjectiveRow(s.objective)
}

func (s *Solver) leavingSymbolForMarkerSymbol(marker *internal.Symbol) *internal.Symbol {
	r1 := math.MaxFloat64
	r2 := math.MaxFloat64

	var first, second, third *internal.Symbol

	for symbol, row := range s.rows {
		c := row.CoefficientForSymbol(marker)
		if c == 0 {
			continue
		}

		if symbol.Type == internal.External {
			third = symbol
		} else if c < 0 {
			r := -row.Constant / c
			if r < r1 {
				r1 = r
				first = symbol
			}
		} else {
			r := -row.Constant / c
			if r < r2 {
				r2 = r
				second = symbol
			}
		}
	}

	if first != nil {
		return first
	}
	if second != nil {
		return second
	}
	return third
}

func (s *Solver) removeConstraintEffects(c *Constraint, tag *internal.Tag) {
	if tag.Marker.Type == internal.Error {
		s.removeMarkerEffects(tag.Marker, float64(c.priority))
	}
	if tag.Other.Type == internal.Error {
		s.removeMarkerEffects(tag.Other, float64(c.priority))
	}
}

func (s *Solver) removeMarkerEffects(marker *internal.Symbol, strength float64) {
	row, ok := s.rows[marker]

	if ok {
		s.objective.InsertRow(row, -strength)
	} else {
		s.objective.InsertSymbol(marker, -strength)
	}
}

func (s *Solver) createRow(c *Constraint, tag *internal.Tag) *internal.Row {
	expression := FromExpression(c.expression)

	row := internal.NewRow(expression.constant)

	for _, term := range expression.terms {
		if internal.IsNearZero(term.coefficient) {
			continue
		}

		symbol := s.symbolForVariable(term.variable)

		foundRow, ok := s.rows[symbol]

		if ok {
			row.InsertRow(foundRow, term.coefficient)
		} else {
			row.InsertSymbol(symbol, term.coefficient)
		}
	}

	switch c.relation {
	case LessThanOrEqualTo, GreaterThanOrEqualTo:
		coefficient := 1.0
		if c.relation == GreaterThanOrEqualTo {
			coefficient = -1
		}

		slack := &internal.Symbol{internal.Slack}

		tag.Marker = slack
		row.InsertSymbol(slack, coefficient)

		if c.priority < PriorityRequired {
			error := &internal.Symbol{internal.Error}

			tag.Other = error
			row.InsertSymbol(error, -coefficient)
			s.objective.InsertSymbol(error, coefficient)
		}
	case EqualTo:
		if c.priority < PriorityRequired {
			errPlus := &internal.Symbol{internal.Error}
			errMinus := &internal.Symbol{internal.Error}
			tag.Marker = errPlus
			tag.Other = errMinus
			row.InsertSymbol(errPlus, -1.0)
			row.InsertSymbol(errMinus, 1.0)

			s.objective.InsertSymbol(errPlus, float64(c.priority))
			s.objective.InsertSymbol(errMinus, float64(c.priority))
		} else {
			dummy := &internal.Symbol{internal.Dummy}
			tag.Marker = dummy
			row.InsertSymbol(dummy, 1.0)
		}
	}

	if row.Constant < 0.0 {
		row.ReverseSign()
	}

	return row
}

func (s *Solver) chooseSubjectForRow(row *internal.Row, tag *internal.Tag) *internal.Symbol {
	for symbol := range row.Cells {
		if symbol.Type == internal.External {
			return symbol
		}
	}

	if tag.Marker.Type == internal.Slack || tag.Marker.Type == internal.Error {
		if row.CoefficientForSymbol(tag.Marker) < 0.0 {
			return tag.Marker
		}
	}

	if tag.Other.Type == internal.Slack || tag.Other.Type == internal.Error {
		if row.CoefficientForSymbol(tag.Other) < 0.0 {
			return tag.Other
		}
	}

	return &internal.Symbol{internal.Invalid}
}

func (s *Solver) addWithArtificalVariableOnRow(row *internal.Row) (bool, error) {
	artificial := &internal.Symbol{internal.Slack}
	s.rows[artificial] = internal.CopyRow(row)
	artificialRow := internal.CopyRow(row)

	err := s.optimizeObjectiveRow(artificialRow)
	if err != nil {
		return false, err
	}

	success := internal.IsNearZero(artificialRow.Constant)
	artificialRow = internal.NewRow(0)

	if foundRow, ok := s.rows[artificial]; ok {
		delete(s.rows, artificial)

		if len(foundRow.Cells) == 0 {
			return true, nil
		}

		entering := internal.AnyPivotableSymbol(foundRow)
		if entering.Type == internal.Invalid {
			return false, nil
		}

		foundRow.SolveForSymbols(artificial, entering)
		s.substitute(entering, foundRow)
		s.rows[entering] = foundRow
	}

	for _, row := range s.rows {
		delete(row.Cells, artificial)
	}
	delete(s.objective.Cells, artificial)
	return success, nil
}

func (s *Solver) substitute(symbol *internal.Symbol, row *internal.Row) {
	for key, secRow := range s.rows {
		secRow.Substitute(symbol, row)

		if key.Type != internal.External && secRow.Constant < 0.0 {
			s.infeasibleRows.PushBack(key)
		}
	}
	s.objective.Substitute(symbol, row)

	if s.artificial != nil {
		s.artificial.Substitute(symbol, row)
	}
}

func (s *Solver) optimizeObjectiveRow(objective *internal.Row) error {

	for true {
		entering := s.enteringSymbolForObjectiveRow(objective)
		if entering.Type == internal.Invalid {
			return nil
		}

		leaving := s.leavingSymbolForEnteringSymbol(entering)
		if leaving == nil {
			panic("UNDEFINED symbol!")
		}

		row, _ := s.rows[leaving]

		delete(s.rows, leaving)

		row.SolveForSymbols(leaving, entering)

		s.substitute(entering, row)
		s.rows[entering] = row
	}

	return errors.New("Never ever")
}

func (s *Solver) enteringSymbolForObjectiveRow(objective *internal.Row) *internal.Symbol {
	for symbol, val := range objective.Cells {
		if symbol.Type != internal.Dummy && val < 0.0 {
			return symbol
		}
	}
	return &internal.Symbol{internal.Invalid}
}

func (s *Solver) leavingSymbolForEnteringSymbol(entering *internal.Symbol) *internal.Symbol {
	ratio := math.MaxFloat64
	var result *internal.Symbol

	for symbol, row := range s.rows {
		if symbol.Type == internal.External {
			continue
		}

		temp := row.CoefficientForSymbol(entering)
		if temp < 0 {
			tempRatio := -row.Constant / temp
			if tempRatio < ratio {
				ratio = tempRatio
				result = symbol
			}
		}
	}

	return result
}

func (s *Solver) symbolForVariable(v *Variable) *internal.Symbol {
	symbol, ok := s.variables[v]

	if ok {
		return symbol
	}

	symbol = &internal.Symbol{internal.External}
	s.variables[v] = symbol
	return symbol
}

func (s *Solver) AddEditVariable(v *Variable, priority float64) error {
	if _, ok := s.edits[v]; ok {
		return errors.New("DUPLICATE")
	}

	if priority < 0 || Priority(priority) == PriorityRequired {
		return errors.New("Bad Priority")
	}

	constraint := NewConstraint(NewExpression([]*Term{NewTerm(v, 1.0)}, 0.0), EqualTo)
	constraint.priority = Priority(priority)

	s.AddConstraint(constraint)

	info := &editInfo{
		tag:        s.constraints[constraint],
		constraint: constraint,
		constant:   0.0,
	}

	s.edits[v] = info

	return nil
}

func (s *Solver) SuggestValueForVariable(v *Variable, value float64) {
	edit, ok := s.edits[v]

	if !ok {
		return
	}

	s.suggestValueForEditInfoWithoutDualOptimization(edit, value)

	s.dualOptimize()
}

func (s *Solver) suggestValueForEditInfoWithoutDualOptimization(info *editInfo, val float64) {
	delta := val - info.constant
	info.constant = val

	{
		symbol := info.tag.Marker
		row, ok := s.rows[symbol]

		if ok {
			if row.Add(-delta) < 0.0 {
				s.infeasibleRows.PushBack(symbol)
			}
			return
		}

		symbol = info.tag.Other
		row, ok = s.rows[symbol]

		if ok {
			if row.Add(delta) < 0.0 {
				s.infeasibleRows.PushBack(symbol)
			}
			return
		}
	}

	for symbol, row := range s.rows {
		coeff := row.CoefficientForSymbol(info.tag.Marker)
		if coeff != 0.0 && row.Add(delta*coeff) < 0.0 && symbol.Type != internal.External {
			s.infeasibleRows.PushBack(symbol)
		}
	}
}

type Update struct {
	Context    interface{}
	UpdatedVal float64
}

func (s *Solver) FlushUpdates() []*Update {
	result := make([]*Update, 0)

	for variable, symbol := range s.variables {
		row, ok := s.rows[symbol]

		var updatedValue float64 = 0
		if ok {
			updatedValue = row.Constant
		}

		variable.applyUpdate(updatedValue)

		if variable.owner != nil && variable.owner.Context != nil {
			result = append(result, &Update{variable.owner.Context, variable.Value})
		}
	}

	return result
}
func (s *Solver) dualOptimize() {
	for s.infeasibleRows.Len() > 0 {
		e := s.infeasibleRows.Back()
		s.infeasibleRows.Remove(e)

		leaving := e.Value.(*internal.Symbol)

		row, ok := s.rows[leaving]

		if ok && row.Constant < 0.0 {
			entering := s.dualEnteringSymbolForRow(row)

			delete(s.rows, leaving)

			row.SolveForSymbols(leaving, entering)
			s.substitute(entering, row)
			s.rows[entering] = row
		}
	}
}

func (s *Solver) dualEnteringSymbolForRow(row *internal.Row) *internal.Symbol {
	var entering *internal.Symbol

	ratio := math.MaxFloat64

	for symbol, val := range row.Cells {

		if val > 0 && symbol.Type != internal.Dummy {
			coeff := s.objective.CoefficientForSymbol(symbol)
			r := coeff / val
			if r < ratio {
				ratio = r
				entering = symbol
			}
		}
	}

	return entering
}
