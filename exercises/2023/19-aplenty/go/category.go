package exercises

type Category int

const (
	ExtremlyCool Category = iota + 1
	Musical
	Aerodynamic
	Shiny
	Default
)

func getCategory(s rune) Category {
	switch s {
	case 'x':
		return ExtremlyCool
	case 'm':
		return Musical
	case 'a':
		return Aerodynamic
	case 's':
		return Shiny
	default:
		return Default
	}
}
