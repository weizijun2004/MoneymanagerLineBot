package structType

type Group struct {
	GroupName    string
	ExistEvents  []string // exist events in group
	GroupMembers []string // memebers in group
}

type EventsMember struct {
	EventsMember map[string][]string // what members are included in each event
}

// if knows which members in which event, than just catch these members docs to change

type Member struct {
	MemberName string
	Events     map[string]Event // event name -> event
}

type Event struct {
	EventName   string
	MemebersPay map[string]int // member name -> pay
}
