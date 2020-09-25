
## 代码结构

```bash

```

### ginhttp
```
s := ginhttp.NewServer(ginhttp.Addr(":4000))
s.AddBeforeServerStartFunc(bs.InitPprof(), bs.InitExpvar())
s.AddAfterServerStopFunc(bs.CloseLogger())
s.Serve();
```