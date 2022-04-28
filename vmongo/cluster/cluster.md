## 集群部署(分片)

先启动配置节点，再启动分片节点，最后启动路由节点

####shard0
/usr/local/mongoDB/bin/mongod --port 27020 --dbpath=/Users/zhangwei/Desktop/mongodb/shard/shard0/data0 --logpath=/Users/zhangwei/Desktop/mongodb/shard/shard0/data0/shard0.log --logappend --fork -f /Users/zhangwei/Desktop/mongodb/shard/shard.config --shardsvr --replSet shard0
/usr/local/mongoDB/bin/mongod --port 27021 --dbpath=/Users/zhangwei/Desktop/mongodb/shard/shard0/data1 --logpath=/Users/zhangwei/Desktop/mongodb/shard/shard0/data1/shard1.log --logappend --fork -f /Users/zhangwei/Desktop/mongodb/shard/shard.config --shardsvr --replSet shard0

/usr/local/mongoDB/bin/mongo 127.0.0.1:27020
rs.initiate({_id: 'shard0', members: [{_id: 0, host: 'localhost:27020'}, {_id: 1, host: 'localhost:27021'}]})

####shard1
/usr/local/mongoDB/bin/mongod --port 27030 --dbpath=/Users/zhangwei/Desktop/mongodb/shard/shard1/data0 --logpath=/Users/zhangwei/Desktop/mongodb/shard/shard1/data0/shard0.log --logappend --fork -f /Users/zhangwei/Desktop/mongodb/shard/shard.config --shardsvr --replSet shard1
/usr/local/mongoDB/bin/mongod --port 27031 --dbpath=/Users/zhangwei/Desktop/mongodb/shard/shard1/data1 --logpath=/Users/zhangwei/Desktop/mongodb/shard/shard1/data1/shard0.log --logappend --fork -f /Users/zhangwei/Desktop/mongodb/shard/shard.config --shardsvr --replSet shard1

/usr/local/mongoDB/bin/mongo 127.0.0.1:27030
rs.initiate({_id: 'shard1', members: [{_id: 0, host: 'localhost:27030'}, {_id: 1, host: 'localhost:27031'}]})


####shard2
/usr/local/mongoDB/bin/mongod --port 27040 --dbpath=/Users/zhangwei/Desktop/mongodb/shard/shard2/data0 --logpath=/Users/zhangwei/Desktop/mongodb/shard/shard2/data0/shard0.log --logappend --fork -f /Users/zhangwei/Desktop/mongodb/shard/shard.config --shardsvr --replSet shard2
/usr/local/mongoDB/bin/mongod --port 27041 --dbpath=/Users/zhangwei/Desktop/mongodb/shard/shard2/data1 --logpath=/Users/zhangwei/Desktop/mongodb/shard/shard2/data1/shard0.log --logappend --fork -f /Users/zhangwei/Desktop/mongodb/shard/shard.config --shardsvr --replSet shard2

/usr/local/mongoDB/bin/mongo 127.0.0.1:27040
rs.initiate({_id: 'shard2', members: [{_id: 0, host: 'localhost:27040'}, {_id: 1, host: 'localhost:27041'}]})


####config
/usr/local/mongoDB/bin/mongod --port 27100 --dbpath=/Users/zhangwei/Desktop/mongodb/config/config0 --logpath=/Users/zhangwei/Desktop/mongodb/config/config0/config0.log --logappend --fork  -f /Users/zhangwei/Desktop/mongodb/config/config.config --configsvr --replSet=config
/usr/local/mongoDB/bin/mongod --port 27101 --dbpath=/Users/zhangwei/Desktop/mongodb/config/config1 --logpath=/Users/zhangwei/Desktop/mongodb/config/config1/config1.log --logappend --fork -f /Users/zhangwei/Desktop/mongodb/config/config.config  --configsvr --replSet=config
/usr/local/mongoDB/bin/mongod --port 27102 --dbpath=/Users/zhangwei/Desktop/mongodb/config/config2 --logpath=/Users/zhangwei/Desktop/mongodb/config/config2/config2.log --logappend --fork -f /Users/zhangwei/Desktop/mongodb/config/config.config --configsvr --replSet=config

/usr/local/mongoDB/bin/mongo 127.0.0.1:27100
rs.initiate({_id: 'config', members: [{_id: 0, host: 'localhost:27100'}, {_id: 1, host: 'localhost:27101'},  {_id: 2, host: 'localhost:27102'}]})


####route
/usr/local/mongoDB/bin/mongos --port 40000 --configdb config/localhost:27100,localhost:27101,localhost:27102 --fork --logpath=/Users/zhangwei/Desktop/mongodb/route/route0/route0.log --logappend -f /Users/zhangwei/Desktop/mongodb/route/route.config
/usr/local/mongoDB/bin/mongos --port 40001 --configdb config/localhost:27100,localhost:27101,localhost:27102 --fork --logpath=/Users/zhangwei/Desktop/mongodb/route/route0/route0.log --logappend -f /Users/zhangwei/Desktop/mongodb/route/route.config

#### set shard
db.runCommand({ addshard: 'shard0/localhost:27020,localhost:27021'})
db.runCommand({ addshard: 'shard1/localhost:27030,localhost:27031'})

#### 查询所有数据
db.col.find().pretty()

####使库支持分片
db.runCommand({ enablesharding: 'test'})
####建索引
####db.col.ensureIndex({name:"hashed"})
####分片
sh.shardCollection("db.col", {name:"hashed"})


####查看分片情况
db.col.stats();
####查看数据分布
db.col.getShardDistribution()
