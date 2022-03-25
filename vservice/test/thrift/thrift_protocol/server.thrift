namespace go testService

// thrift -I {self location}/tool/idl/thrift
// -gen go:package_prefix=vtool/idl/thrift/gen-go/ test_server.thrift
include "thrift.thrift"

struct SayHelloReq {
    1: i64 val
}

struct SayHelloData {
    1: string val
}

struct SayHelloRes {
    1: optional SayHelloData data
    2: optional thrift.ErrInfo errInfo
}

service TestService {
    SayHelloRes SayHello(1:SayHelloReq req)
}