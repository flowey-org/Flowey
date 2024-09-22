# Usage

1. Start the server:

```bash session
$ (backend) go run . localhost:5000
```

2. Run the client:

```bash session
$ (backend/client) node ws.js localhost:5000
> ws.send("Hello, world!")
< Hello, world!
>
(To exit, press Ctrl+C again or Ctrl+D or type .exit)
>
```
