package context

type (
	Context struct {
	}

	IContext interface {
	}
)

func New() IContext {
	return &Context{}
}
