# gitub.com/gkzy/gow/lib/rpc

### server 

```go

// InitRPCServer 
func InitRPCServer(){
    g,err:=rpc.NewServer(10001)
    if err!=nil{
        panic(err)
    }
    handler(g.Server)
    g.Run()
}

// handler register struct
func handler(g *rpc.Server){

}

```

### client

```
client,err:=rpc.NewClient("192.168.0.100",10001)
if err!=nil{
        panic(err)
}
...
```