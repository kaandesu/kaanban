package main

type Step struct {
	status      int
	title       string
	description string
}

func NewTask(status int, title, description string) Step {
	return Step{status: status, title: title, description: description}
}

func (t *Step) Next() {
	if t.status == TotalStages {
		t.status = 0
	} else {
		t.status++
	}
}

// implement the list.Item interface
func (t Step) FilterValue() string {
	return t.title
}

func (t Step) Title() string {
	return t.title
}

func (t Step) Description() string {
	return t.description
}
