# itemfilterRPC

基于redis的物品过滤rpc服务,提供上下文去重和时间间隔去重功能

## 注意事项

1. 本服务仅针对小时级数据去重,天级的请离线算好后与本项目结果配合使用.
2. 去重使用布隆过滤器,因此存在碰撞概率,不能保证一定不重复
3. 使用插件[RedisBloom](https://github.com/RedisBloom/RedisBloom)

## 接口类型

+ grpc
+ http-RESTful(需要启动是使用`--Use_Gateway`开启,同时会启动swagger用于调试)

## 概念

+ `EntitySource`,简写`ES`,指实体资源,也就是要去重的对象
+ `EntitySourceType`,简写`EST`指实体资源类型,即要去重的对象类型,比如我们要去重的对象有漫画,音乐,电影,他们往往有自己单独的id,因此分别进行去重
+ `EntitySourceID`,简写`ESID`,指实体资源id,即要去重的对象.
+ `BlackListID`,简写`BID`,指黑名单.
+ `RangeID`,简写`RID`,指范围的id,往往一个范围下的资源类型和资源范围都是相对固定的,规定每个范围只能包含同一类实体
+ `PickerID`,简写`PID`,指挑选者id,挑选者就是获取资源的对象.
+ `PickerInterval`,简写`PI`,会根据服务设置记录用户一个时间范围内的资源使用情况,用于长期去重,可以选的时间范围包括
    + `day`,当天(以utc时间记)
    + `week`,本周
    + `month`,本月
    + `year`,本年
    + `infi`,从记录开始
+ `ContextID`,简写`CID`,指上下文id,上下文可以用于限制一次去重的范围.一次上下文中的资源不会重复,上下文在每次有行为时都会刷新过期时间,这意味着它可以记录活跃行为序列中使用的资源,可以用于短期去重.
+ `Candidate`,简写`C`,指请求中用于操作的实体
+ `Survivor`,简写`S`,指请求操作后剩下的符合要求的实体
+ `PickerCounter`,简写`PC`,记录使用用户数的计数器,每天会记录一份,可以设置过期
+ `PickerCounterDate`,简写`PCD`,`PickerCounter`当天的日期
+ `Duration`,简写`D`,指时间间隔,我们可以指定间隔以确保从当前时间点到这个间隔内不会重复.

+ 操作,本服务提供2种类型的操作器
    1. 检查器`check`,用于检查物品是否在集合中
    2. 过滤器`filter`用于将候选物品中符合要求的保留,不符合要求的删除

## 基本思路

本项目基于如下假设:

1. 范围下的实体资源规模很大
2. 黑名单规模很小,数量级最大为千级,且绝对不可以出现
3. 范围和黑名单独立
4. ID均为字符串类型
5. 需要过滤的物品受范围,黑名单,会话,挑选者3个因素制约变化

本过滤器的过滤流程为

> 设置:

1. 预先设置范围,包括范围使用概率型过滤器设置存在项,
2. 预先设置黑名单,只能使用`set`结构设置黑名单
3. 设置挑选者在时间范围下已经使用过的实体资源(根据设置保存每天/周/月/年/总体各一份,写入不会刷新过期时间)
4. 设置不同上下文中已经使用过的实体(过期时长由请求端带入,如果缺省则默认使用设置的默认值作为刷新时间,写入会刷新过期时间)

> 过滤: 请求中的实体id如果在结果中对应的值为True则表示可用,False则表示不可用

1. 使用set首先去掉请求中的重复项
2. 如果有带挑选者ID将挑选者add到当天的hyperloglog中(配置如果设置)
3. (如果请求中有带黑名单信息,或范围有配置引用的黑名单)判断在不在这两个黑名单中,如果在则该项返回False
4. (去过请求中有带范围信息)
   1. 判断在指定范围中是否可用
   2. 如果有带挑选者ID,在对应的当天hyperloglog中(配置如果设置)添加当前挑选者
5. 如果有带挑选者ID:
    1. 如果有指定会话但没有指定时间间隔,则判断是否在会话中已经被使用
    2. 如果有指定时间间隔但没有指定会话,则判断是否在指定的时间间隔中已经被使用
    3. 如果时间和会话都有指定则两边都会判断
6. 使用hyperloglog统计当天去重请求的挑选者数量

> 统计流量: (配置如果设置)如果配置中启动该特性,则本工具还可以顺便各个范围下以及总体下的去重挑选者数

