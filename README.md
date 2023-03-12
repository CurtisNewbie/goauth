# goauth v1.0.0

User and resource management service Implemented In Go (For Internal Use).

- [x] Implement resource management
- [x] Implement java client for goauth
- [x] Implement Go client for goauth
- ~~Reimplement auth-service~~ (maybe not, auth-service is working just fine)

### Generate Resource Scripts for Production Environment

Sometimes we may modify the paths and resources in development environment, we can use the following endpoint to generate scripts that can be executed in production environment without modifing these paths and resources again manually.

```sh
curl -X POST "http://localhost:8081/internal/resource/script/generate" \
    -H 'content-type:application/json' \
    -d '{ "resCodes" : ["basic-user", "manage-users"]}' \
    -o output.sql
```