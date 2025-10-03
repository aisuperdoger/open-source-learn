https://mp.weixin.qq.com/s?src=11&timestamp=1759410183&ver=6272&signature=ehxOoWqJMZ*VN6euK0qXJdGgCwSoDSZ1eBIpMPYIRfdvC*uwWGLb4gRq8s8ry1w*NGl3miJRm5DiEMOCRlEEc1FHZvRuCv13tUr1s54uD6FO-smYhvuzREn6lDDzPXE0&new=1
底层实现总结：
- SQL语句映射为Statement结构体。
    - 将Select、Omit、where等内容存储在成员变量中。
    - Scheme结构体来建立记录和对象之间的映射。 Scheme对象存储了表名称、列信息等表结构信息
- gorm库中的方法大致可以分为两类：过程方法和结尾方法。
    - 过程方法：过程方法一般只有构建Statement对象的功能，过程方法不会将SQL发送给数据库执行。常见的过程方法有Where、Select、Omit、Model等。
    - 结尾方法：结尾方法会构建完整Statement对象，并将SQL语句发送到数据库执行。常用的结尾方法有Update、Find、Delete、Create等，实际上也就是CURD功能。
    - 方法都是返回*DB对象，因此可以链式调用。

知识点：
- 创建表
    - AutoMigrate以蛇形命名法转换表名，并将表名自动加上s。通过设置NamingStrategy可以不加s，或者TableName方法可以自定义表名。
    - 结构体中嵌入gorm.Model，会自动拥有ID、CreatedAt、UpdatedAt、DeletedAt字段，然后自动获取软删除能力。
        - 只要你的结构体包含 gorm.DeletedAt 类型的字段（通常叫 DeletedAt），GORM 就会自动开启软删除功能。
    - 属性比较多的结构体，可以使用gorm embed功能，将属性拆分到多个结构体中，并使用gorm的embed标签标识即可。embed内嵌的结构体转换成表字段时，还可以设置统一的前缀。
- 查、改
    - First、Last、Take这三个方法会设置Limit 1。Find则是查询满足条件的所有记录
        - Find则是查询满足条件的所有记录，然后在代码层面映射第一条记录到结果对象中（只传一个结构体变量（非切片）），性能差。
    - Pluck 的作用：只查某一列的值，并直接存到一个切片中。使用Find还要从结构体中获取列的值。
    - Update函数在更新记录时只会更新非零字段。一旦遇到这种情况，应该使用Select + Update语句进行特定字段更新或者使用map[string]interface{} kv对结构进行更新。
    - 判断主键和唯一键冲突：需要告诉go-zero拼接一个能解决冲突的sql。
        - gorm的任务只是替我们写SQL语句，实际上只是在SQL语句中添加了ON DUPLICATE KEY UPDATE xxxxx 表示键冲突更新策略。
    - 全局删除和全局更新操作都是被禁止的
- 其他
    - gorm库提供了DryRun（试运行）配置，即只会打印SQL语句而不会真正的执行。
    - 在建立数据库连接时，DSN建议添加上loc=Local选项。因为默认为美国时间，即UTC，而CST为UTC+8。不加的话存储到数据库中去就少了八个小时。
        - 从系统的/etc/localtime文件中解析并读取当前系统的时区信息。
    - gorm中的事务：实际上gorm是每执行一个命令函数就向数据库发送一条SQL，事务的功能是通过数据库保证的，gorm库并没有对应的逻辑实现事务。

