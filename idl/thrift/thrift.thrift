
// include "thrift.thrift"

namespace go thriftError

struct ErrInfo {
  1: required i32 code
  2: required string msg
}