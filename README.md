# go-direct
Simple server for doing "http redirects" 

[![GitHub language count](https://img.shields.io/github/languages/count/thorsager/go-direct)](https://github.com/thorsager/go-direct)
[![GitHub code size in bytes](https://img.shields.io/github/languages/code-size/thorsager/go-direct)](https://github.com/thorsager/go-direct)
[![Go Report Card](https://goreportcard.com/badge/github.com/thorsager/go-direct)](https://goreportcard.com/report/github.com/thorsager/go-direct)
[![Build Status](https://travis-ci.com/thorsager/go-direct.svg?branch=master)](https://travis-ci.com/thorsager/go-direct)
[![Docker Pulls](https://img.shields.io/docker/pulls/thorsager/go-direct)](https://hub.docker.com/r/thorsager/go-direct)


Configuration is done by setting the `REDIRECTS` env var, the content must be
JSON of the following format:
```json
{
  "/path" : "http://target.url"
}
```
Note that the same structure can be added to `TEMPORARY_REDIRECTS`env var, 
which will the do a "temporarily moved" redirect
