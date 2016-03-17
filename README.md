# gcore
Simple golang web framework

```go
var HomeController = func() func(w http.ResponseWriter, r *http.Request) {
	tpl := &controller.Controller{
		Template: "frontend/home/home.html",
		Layout:   "frontend.html",
	}

    var tplVars struct {}

    return func(w http.ResponseWriter, r *http.Request) {
        tpl.Render(w, r, tplVars)
    }
}()

gcore := gcore.New("<server_addr>")
gcore.AddRoute("/", HomeController)
```

### Using decorators
```go
gcore.RegisterDecorator(decorator.NewLogger())
```

or

```go
gcore.AddRoute("/", HomeController, decorator.NewLogger())
```

### Available decorators
- Admin
- Auth
- Context
- Logger
- Mongo
- Recover