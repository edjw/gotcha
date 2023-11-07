A folder to put partials/fragments for HTMX to request

Add them to the partialsMaps in the `main` function in `main.go`

```go
partialsMap := map[string]func() templ.Component{
    "new_headline": partials.NewHeadline,
}
```