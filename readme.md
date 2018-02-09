# Iron 艾恩

> 从自己的一个项目中抽离出的web框架，用来快速构建后端接口项目。框架底层基于 Echo (高性能的web框架)。 Model层推荐采用Gorm(一个优雅的对象关系映射库)当然，你也可以通过框架内部的中间件，引入自己喜欢的 ORM 库。

## 依赖

* github.com/labstack/echo
* github.com/jinzhu/orm
* gopkg.in/yaml.v2

## 若你想在你的项目里面使用它，请...

```
go get -u github.com/alixez/Iron
```

**如果你经常使用glide作为依赖管理工具**

```
glide get -u github.com/alixez/Iron
```