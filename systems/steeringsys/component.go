package steeringsys

type (
	Component struct {
		behaviors []IBehavior
	}
)

func NewComponent(behaviors []IBehavior) *Component {
	return &Component{behaviors}
}
