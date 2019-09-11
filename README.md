# go-direct
Simple server for doing "http redirects" 

Configuration is done by setting the "REDIRECTS" env var, the content must be
JSON of the following format:
```json
{
  "/path" : "http://target.url"
}
```
