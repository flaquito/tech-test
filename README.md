# Tech test

The frontend was created with `npm create vite@latest frontend`.

The backend was created as a folder with `main.go` and by then running `go mod init example.com/backend`.

## Set up and Run

After cloning this repository,

- for the frontend:

    ```bash
    cd frontend
    npm install
    npm run dev
    ```

    This will spin up the frontend at `localhost:5173`.
- for the backend:

    ```bash
    cd backend
    go run ./cmd/server
    ```

    This will spin up the backend at `localhost:8080`.

## Design / Architecture decision points

### Frontend

For the frontend, I started out with the template provided by Vite. I moved the transport interacting with post data (as in "feed post") and the components to render into dedicated folders, as well as the core domain object of `Post`. Styling with CSS was applied next to each component.

State and data are placed on the app top level to centralize the data flow.

Tag filtering was implemented client-side to avoid increasing the backend complexity. In production, this should be moved server-side to not mix domain data logic into the frontend and to reduce network load.

I decided against including newly created posts in the feed immediately, but instead wait for the WebSocket notification and only show persisted posts, to avoid having to reconcile posts in draft state with potentially failing backend uploads. This becomes advisable for asynchronous media handling / uploads at scale.

I made both image and text mandatory for creating a post (not the tags), showing an image preview with the creation form, as I expect this to be a realistic user experience.

Further ideas would be extracting the hook / effect logic from the app top level into a dedicated folder (it's probably already large enough to warrant this) and especially improving the display of loading state and errors (messages).

### Backend

For the backend, I tried to learn about idiomatic Go project layouts and went with a flat hierarchy without `internal` to keep the structure simple and readable. The packages represent transport, service and database layers, with domain objects in a dedicated own package. The file-based object storage and SQLite DB have a top-level package as well, to separate data from business logic.

The database table for feed post data storage was created with a seed script and interacted with by the database package / layer using plain SQL. It should not be necessary to run the seed script again, since DB and assets (images) are populated and included with the project under `data`. To run it, execute:

```bash
cd backend
go run ./cmd/seed
```

The image normalization was intentionally kept light for this scope, to focus on a working end-to-end flow rather than media optimization.

The Mutex was used for the WebSocket to ensure safe concurrent access and indicate best practices at scale.

When refining this further, I'd start by improving the error flow by introducing typed errors to improve error propagation and HTTP mapping.

## Example API queries

### Fetch a page of posts

```curl
❯ curl -v "http://localhost:8080/posts?page=1"
* Host localhost:8080 was resolved.
* IPv6: ::1
* IPv4: 127.0.0.1
*   Trying [::1]:8080...
* Connected to localhost (::1) port 8080
> GET /posts?page=1 HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/8.7.1
> Accept: */*
> 
* Request completely sent off
< HTTP/1.1 200 OK
< Content-Type: application/json; charset=utf-8
< Date: Tue, 03 Mar 2026 19:02:45 GMT
< Transfer-Encoding: chunked
< 
[
    {
        "id": 33,
        "imageUrl": "http://localhost:8080/images/33.jpg",
        "text": "Hello world",
        "tags": [
            "test",
            "example"
        ],
        "createdAt": "2026-03-03T20:01:44.939718+01:00"
    },
    {
        "id": 25,
        "imageUrl": "http://localhost:8080/images/25.jpg",
        "text": "A day with flowers and moon.",
        "tags": [
            "flowers",
            "moon"
        ],
        "createdAt": "2026-03-03T11:22:34.646015+01:00"
    },
    {
        "id": 24,
        "imageUrl": "http://localhost:8080/images/24.jpg",
        "text": "A day with flowers and dog.",
        "tags": [
            "flowers",
            "dog"
        ],
        "createdAt": "2026-03-03T11:22:34.588724+01:00"
    },
    ...,
    {
        "id": 17,
        "imageUrl": "http://localhost:8080/images/17.jpg",
        "text": "A quiet moment in nature.",
        "tags": [],
        "createdAt": "2026-03-03T11:22:34.16307+01:00"
    }
* Connection #0 to host localhost left intact
]%
```

### Upload a new post

```curl
❯ curl -v http://localhost:8080/uploads \ 
  -F "file=@./example-image.jpg" \
  -F "text=Hello world" \
  -F "tags=test,example"
* Host localhost:8080 was resolved.
* IPv6: ::1
* IPv4: 127.0.0.1
*   Trying [::1]:8080...
* Connected to localhost (::1) port 8080
> POST /uploads HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/8.7.1
> Accept: */*
> Content-Length: 10437
> Content-Type: multipart/form-data; boundary=------------------------cSuRw3MTWcfh2Zs3JoHIvs
> 
* upload completely sent off: 10437 bytes
< HTTP/1.1 201 Created
< Date: Tue, 03 Mar 2026 19:01:44 GMT
< Content-Length: 0
< 
* Connection #0 to host localhost left intact
```

### Connect to websocket

```curl
❯ wscat -c ws://localhost:8080/ws        
Connected (press CTRL+C to quit)
# Upload post in another tab
< {"id":35,"imageUrl":"http://localhost:8080/images/35.jpg","text":"Hello world","tags":["test","example"],"createdAt":"2026-03-03T20:13:14.392641+01:00"}
```

## Notes

- The `backend/cmd/seed/main` script for seeding the database (with Picsum images) is generated by an LLM. Random tags and texts are created for the posts.
