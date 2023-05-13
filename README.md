# goauth

User and resource management service Implemented In Go (For Internal Use).

- [x] Implement resource management
- [x] Implement java client for goauth
- [x] Implement Go client for goauth
- ~~Reimplement auth-service~~ (maybe not, auth-service is working just fine)

## Client Lib Integration

### goauth-client-go

Run following command in the source root of your project (where go.mod is at).

```
# For v1.0.4
go get github.com/curtisnewbie/goauth/client/goauth-client-go@f41e4ca6eaba

# For v1.0.3

go get github.com/curtisnewbie/goauth/client/goauth-client-go@609b199d9b07

# For v1.0.2

go get github.com/curtisnewbie/goauth/client/goauth-client-go@927902695f12

# For v1.0.1

go get github.com/curtisnewbie/goauth/client/goauth-client-go@e7460b23dbaa
```

### goauth-client-java

A spring and feign based implementation that must be installed locally using Maven.

