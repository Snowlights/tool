
// include "thrift.thrift"

namespace go thriftBase

struct ErrInfo {
  1: required i32 code
  2: required string msg
}

struct Context {
    1: optional map<string, string> spanCtx
}