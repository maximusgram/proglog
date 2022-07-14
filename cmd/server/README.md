## Test your API

```bash
$ curl -X POST localhost:8080 -d \
    '{"record": {"value" : "SGVsbG8sIHdvcmxk"}}'

> {"offset":0}
```

Go’s encoding/json package encodes []byte as a base64-encoding string. The
record’s value is a []byte, so that’s why our requests have the base64 encoded
forms of "Hello, world".

You can read the records back by running the following
commands and verifying that you get the associated records back from the
server:

```bash
$ curl -X GET localhost:8080 -d '{"offset": 0}'
```