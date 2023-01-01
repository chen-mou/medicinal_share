# 医疗管家网站

## 主要功能



1. 让用户可以在线和医生进行交流:
    - 使用websocket实现服务器主动推送消息至客户端
    - 使用redis作为消息代理或者使用消息队列作为消息代理（等做分布式时再做这一步）
    - 医生的匹配算法
    - 医生可以在线写药品清单或者推荐如何去就医


2. 构建一个论坛用于让用户在论坛上交流病情:
    - 将论坛根据病情分为交流圈
    - 根据病情描述搜索相应的解决方法
    - 可以求助

3. 商城用于网上药品的购买:
    - 可以根据处方自动下单
    - 解决超卖问题
    
4. 线上医院
    - 用户的病厉
    - 用户预约挂号
