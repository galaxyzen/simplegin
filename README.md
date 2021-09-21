# simplegin

仿照 [gin](https://github.com/gin-gonic/gin) 框架的一个简单实现，代码主要参考了 [gin](https://github.com/gin-gonic/gin) 和 [gee](https://github.com/geektutu/7days-golang/tree/master/gee-web)

路由算法基于前缀树实现，支持`*path`通配符和`:path`通配符，`*path`可以匹配任意多个路径（目录），`:path`可以匹配任意一个路径（目录）：

```
/post/:category/:id
可以匹配 /post/golang/2009、/post/cpp/1983 等

/static/*path/index.js
可以匹配 /static/index.js、/static/home/index.js、/static/home/api/index.js 等

/*path/post/:category/2009
可以匹配 /post/golang/2009、/bookmark/private/post/java/2009 等
```

具体实现和 leetcode第10题：[正则表达式匹配 ](https://leetcode-cn.com/problems/regular-expression-matching/)类似

另外支持分组控制、中间件等。
