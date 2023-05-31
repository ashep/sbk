package notifier

type Noop struct {
}

func NewNoop() Notifier {
	return &Noop{}
}

func (n *Noop) Notify(string) error {
	return nil
}
