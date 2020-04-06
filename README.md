# Monzo++

## API

### Errors

Across each API, if an error occurs, whether it be an server error or just a validation error, it will be returned in a standardised JSON format (as below).

```json
{
	"error": "error message"
}
```

However, if the environment variable `HTTP_ERROR` is set to `DEBUG`, a detailed stack trace will be returned as well - which looks like this:

```json
{
	"error": "error message",
	"stackTrace": "probably quite a long trace"
}
```

This is ideal for debugging and development but probably not for production. When in production mode either don't set the `HTTP_ERROR` variable or set it to something other than `DEBUG`.
