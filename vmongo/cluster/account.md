## 用户
super:
db.createUser({
user:"root",
pwd:"123456",
roles:[
{
role:"userAdminAnyDatabase",
db:"admin"
}
]
})

normal:
db.createUser({
user:"user1",
pwd:"pwd1",
roles:[
{
role:"readWrite",
db:"test"
}
]
})