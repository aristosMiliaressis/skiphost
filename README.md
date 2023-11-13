skiphost
==

sends two requests ont to the base url and one to a non existent path than checks if any of the following apply to both responses

- 502, 503, 504 status code
- redirect cross origin
- redirect cross scheme

**output**
```json
{ "looks_good": false, "reason": "502" }
```