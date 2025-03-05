## Summary
The project contains two executables, a client & a server.

The client monitors a directory for changes, and posts them to an HTTP server which writes them to a directory.

An example of running both, from this directory:
```bash
go run server/cmd/main.go /Users/zacharybriggs/sync/to
go run client/cmd/main.go /Users/zacharybriggs/sync/from
```

## Notes
Basic principle of client application:
1. Poll the file system for changes
2. Submit them to an in-memory cache, compare cache to new changes
3. Changes detected against the cache are posted to the server

The server is _very_ simple, it:
1. Accepts the post request, and reads back into the shared FileChange{} struct
2. Calls os package methods to create/delete the file

The current setup has some caveats:
- On startup, the client will currently immediately post the contents of the dir to the server
  - We could make an implementation of persistence that is stored on disk, to get state on startup
  - We could ignore the first read of the directory, and only monitor further changes
- The post mechanism isn't very sophisticated, it's not been tested for very large files.
- The in-memory cache currently stores file contents, this isn't scalable. Probably better to read the file contents at time of posting