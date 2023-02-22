package room

type Room struct {
	ID      string
	Members []User
}

type User struct {
	ID   string
	Name string
}
