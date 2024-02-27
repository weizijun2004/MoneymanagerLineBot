package structType

type Event struct {
	EventName  string
	MemeberPay map[string]int
}

type User struct {
	MemberName  string
	EventAttend map[string]Event
}

type EventExist struct {
	ExistEventNum int
	ExistEventArr []string
}
