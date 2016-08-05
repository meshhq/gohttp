MeshRedis
======
### Use

When first using MeshRedis, you need to supply a URL for it to connect to via its `Connect` method. In order to get a new session, you must then call `NewSession` for a new Redis Session. All interactions w/ redis needs to be through the RedisSession objects. 

After you obtin a RedisSession through the `NewSession` method, you need to close its session by calling the `Close()` on the RedisSession instance.
