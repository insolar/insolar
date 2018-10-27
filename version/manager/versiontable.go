package manager



type {{.MyType}}List struct {
	q []{{.MyType}}
}

func New{{.MyType}}List() *{{.MyType}}List {
	return &{{.MyType}}List{
		q: {{.MyType}}{},
}
}

