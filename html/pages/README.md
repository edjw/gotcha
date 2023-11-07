The directory for storing pages. Add them to the pagesMap in the `main` function in `main.go`

```go
pagesMap := map[string]func() templ.Component{
    "/":      pages.Home,
    "/about": pages.About,
}
```
