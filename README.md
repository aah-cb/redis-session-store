# redis-session-store
Redis Session Store for aah framework


**init.go**

_ "github.com/aah-cb/redis-session-store"
 
**security.conf**
```roboconf
   # Session store is to choose where session value should be persisted.
    store {
      # Currently aah framework supports `cookie` and `file` as store type.
      # Also framework provide extensible `session.Storer` interface to
      # add custom session store.
      # Default value is `cookie`.
      type = "redis"

      # Filepath is used for file store to store session file in the file system.
      # This is only applicable for `type = "file"`, make sure application has
      # Read/Write access to the directory. Provide absolute path.
     #  filepath = "sessions"
      redis {
          # the redis network option, "tcp"
          network = "tcp"
          # the redis address option, "127.0.0.1:6379"
          addr = "127.0.0.1:6379"
          # Password string .If no password then no 'AUTH'. Default ""
          password = ""
          #  If Database is empty "" then no 'SELECT'. Default ""
          database = ""
          #  Prefix "myprefix-for-this-website". Default ""
          prefix = "session_"
          # MaxIdle 0 no limit
          max_idle = 10
          # MaxActive 0 no limit
          max_active = 30
      }
    }
```roboconf
