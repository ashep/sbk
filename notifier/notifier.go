package notifier

type Notifier interface {
	Notify(msg string) error
}
