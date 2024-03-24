
## QA

### 路径问题

打包后可执行文件在bin目录下，因此可执行目录就是bin。bin/app.yml就可以读到。
在IDE中可执行目录其实是Goland的一个目录，因此在启动时设置一个环境变量，避免去Goland目录放这个文件，不好管理。

静态资源同样的道理，使用绝对路径解决问题。

```go
//bin目录同级
frontendPath := filepath.Join(executePath, "../frontend")
```
