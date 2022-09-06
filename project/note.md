## docker-compose

### network

使用 docker-compose 构建的集群会默认创建一个一个 docker network，并将所有的服务注册进去，不同的服务之间可以通过服务名进行访问

### down

docker-compose down 会移除对应的 container

## Http Request

### go request

https://blog.csdn.net/ilini/article/details/110069526?spm=1001.2101.3001.6650.3&utm_medium=distribute.pc_relevant.none-task-blog-2%7Edefault%7ECTRLIST%7ERate-3-110069526-blog-111583382.pc_relevant_aa_2&depth_1-utm_source=distribute.pc_relevant.none-task-blog-2%7Edefault%7ECTRLIST%7ERate-3-110069526-blog-111583382.pc_relevant_aa_2&utm_relevant_index=4

go 使用 http request 发送请求后，得到的 response 需要使用

```
defer response.Body.Close()，这是因为会造成goroutine泄漏，
```
