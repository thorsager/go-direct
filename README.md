# go-direct
Simple server for doing "http redirects" 

Configuration is done by setting the `REDIRECTS` env var, the content must be
JSON of the following format:
```json
{
  "/path" : "http://target.url"
}
```
Note that the same structure can be added to `TEMPORARY_REDIRECTS`env var, 
which will the do a "temporarily moved" redirect
