package types

type Commit struct {
	Data       []string
	ApplyDoneC chan<- struct{}
}
