package structType

type EventExist struct {
	ExistEventArr []string
}
type Event struct {
	EventName  string
	MemeberPay map[string]int
}

type User struct {
	MemberName  string
	EventAttend map[string]Event
}
