
# MVP 0.1
## API

### /sign-in/

### /upload/
Purpose:
```http
POST /upload
Content-Type: multipart/form-data
Authorization:


```

Notes:

### /archive/

# Domain modeling

- Content: hash, URL
- Archive: A users private collection of media
- User: A person.
- Post: A collection of media shared from the Users archive.

See http://mermaid.js.org/syntax/entityRelationshipDiagram.html
```mermaid
erDiagram
CONTENT {
  string hash PK
  string url
}
```

# User Stories?

See http://mermaid.js.org/syntax/userJourney.html
```mermaid

```
