// Autogenerated by Thrift Compiler (0.9.2)
// DO NOT EDIT UNLESS YOU ARE SURE THAT YOU KNOW WHAT YOU ARE DOING

package testService

import (
	"bytes"
	"fmt"
	"git.apache.org/thrift.git/lib/go/thrift"
	"vtool/idl/thrift/gen-go/thriftError"
)

// (needed to ensure safety because of naive import list construction.)
var _ = thrift.ZERO
var _ = fmt.Printf
var _ = bytes.Equal

var _ = thriftError.GoUnusedProtection__
var GoUnusedProtection__ int

type SayHelloReq struct {
	Val int64 `rpc_client:"val,1" json:"val"`
}

func NewSayHelloReq() *SayHelloReq {
	return &SayHelloReq{}
}

func (p *SayHelloReq) GetVal() int64 {
	return p.Val
}
func (p *SayHelloReq) Read(iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(); err != nil {
		return fmt.Errorf("%T read error: %s", p, err)
	}
	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if err != nil {
			return fmt.Errorf("%T field %d read error: %s", p, fieldId, err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 1:
			if err := p.ReadField1(iprot); err != nil {
				return err
			}
		default:
			if err := iprot.Skip(fieldTypeId); err != nil {
				return err
			}
		}
		if err := iprot.ReadFieldEnd(); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(); err != nil {
		return fmt.Errorf("%T read struct end error: %s", p, err)
	}
	return nil
}

func (p *SayHelloReq) ReadField1(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI64(); err != nil {
		return fmt.Errorf("error reading field 1: %s", err)
	} else {
		p.Val = v
	}
	return nil
}

func (p *SayHelloReq) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("SayHelloReq"); err != nil {
		return fmt.Errorf("%T write struct begin error: %s", p, err)
	}
	if err := p.writeField1(oprot); err != nil {
		return err
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return fmt.Errorf("write field stop error: %s", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return fmt.Errorf("write struct stop error: %s", err)
	}
	return nil
}

func (p *SayHelloReq) writeField1(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("val", thrift.I64, 1); err != nil {
		return fmt.Errorf("%T write field begin error 1:val: %s", p, err)
	}
	if err := oprot.WriteI64(int64(p.Val)); err != nil {
		return fmt.Errorf("%T.val (1) field write error: %s", p, err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return fmt.Errorf("%T write field end error 1:val: %s", p, err)
	}
	return err
}

func (p *SayHelloReq) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("SayHelloReq(%+v)", *p)
}

type SayHelloData struct {
	Val string `rpc_client:"val,1" json:"val"`
}

func NewSayHelloData() *SayHelloData {
	return &SayHelloData{}
}

func (p *SayHelloData) GetVal() string {
	return p.Val
}
func (p *SayHelloData) Read(iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(); err != nil {
		return fmt.Errorf("%T read error: %s", p, err)
	}
	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if err != nil {
			return fmt.Errorf("%T field %d read error: %s", p, fieldId, err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 1:
			if err := p.ReadField1(iprot); err != nil {
				return err
			}
		default:
			if err := iprot.Skip(fieldTypeId); err != nil {
				return err
			}
		}
		if err := iprot.ReadFieldEnd(); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(); err != nil {
		return fmt.Errorf("%T read struct end error: %s", p, err)
	}
	return nil
}

func (p *SayHelloData) ReadField1(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return fmt.Errorf("error reading field 1: %s", err)
	} else {
		p.Val = v
	}
	return nil
}

func (p *SayHelloData) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("SayHelloData"); err != nil {
		return fmt.Errorf("%T write struct begin error: %s", p, err)
	}
	if err := p.writeField1(oprot); err != nil {
		return err
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return fmt.Errorf("write field stop error: %s", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return fmt.Errorf("write struct stop error: %s", err)
	}
	return nil
}

func (p *SayHelloData) writeField1(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("val", thrift.STRING, 1); err != nil {
		return fmt.Errorf("%T write field begin error 1:val: %s", p, err)
	}
	if err := oprot.WriteString(string(p.Val)); err != nil {
		return fmt.Errorf("%T.val (1) field write error: %s", p, err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return fmt.Errorf("%T write field end error 1:val: %s", p, err)
	}
	return err
}

func (p *SayHelloData) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("SayHelloData(%+v)", *p)
}

type SayHelloRes struct {
	Data    *SayHelloData        `rpc_client:"data,1" json:"data"`
	ErrInfo *thriftError.ErrInfo `rpc_client:"errInfo,2" json:"errInfo"`
}

func NewSayHelloRes() *SayHelloRes {
	return &SayHelloRes{}
}

var SayHelloRes_Data_DEFAULT *SayHelloData

func (p *SayHelloRes) GetData() *SayHelloData {
	if !p.IsSetData() {
		return SayHelloRes_Data_DEFAULT
	}
	return p.Data
}

var SayHelloRes_ErrInfo_DEFAULT *thriftError.ErrInfo

func (p *SayHelloRes) GetErrInfo() *thriftError.ErrInfo {
	if !p.IsSetErrInfo() {
		return SayHelloRes_ErrInfo_DEFAULT
	}
	return p.ErrInfo
}
func (p *SayHelloRes) IsSetData() bool {
	return p.Data != nil
}

func (p *SayHelloRes) IsSetErrInfo() bool {
	return p.ErrInfo != nil
}

func (p *SayHelloRes) Read(iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(); err != nil {
		return fmt.Errorf("%T read error: %s", p, err)
	}
	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if err != nil {
			return fmt.Errorf("%T field %d read error: %s", p, fieldId, err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 1:
			if err := p.ReadField1(iprot); err != nil {
				return err
			}
		case 2:
			if err := p.ReadField2(iprot); err != nil {
				return err
			}
		default:
			if err := iprot.Skip(fieldTypeId); err != nil {
				return err
			}
		}
		if err := iprot.ReadFieldEnd(); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(); err != nil {
		return fmt.Errorf("%T read struct end error: %s", p, err)
	}
	return nil
}

func (p *SayHelloRes) ReadField1(iprot thrift.TProtocol) error {
	p.Data = &SayHelloData{}
	if err := p.Data.Read(iprot); err != nil {
		return fmt.Errorf("%T error reading struct: %s", p.Data, err)
	}
	return nil
}

func (p *SayHelloRes) ReadField2(iprot thrift.TProtocol) error {
	p.ErrInfo = &thriftError.ErrInfo{}
	if err := p.ErrInfo.Read(iprot); err != nil {
		return fmt.Errorf("%T error reading struct: %s", p.ErrInfo, err)
	}
	return nil
}

func (p *SayHelloRes) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("SayHelloRes"); err != nil {
		return fmt.Errorf("%T write struct begin error: %s", p, err)
	}
	if err := p.writeField1(oprot); err != nil {
		return err
	}
	if err := p.writeField2(oprot); err != nil {
		return err
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return fmt.Errorf("write field stop error: %s", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return fmt.Errorf("write struct stop error: %s", err)
	}
	return nil
}

func (p *SayHelloRes) writeField1(oprot thrift.TProtocol) (err error) {
	if p.IsSetData() {
		if err := oprot.WriteFieldBegin("data", thrift.STRUCT, 1); err != nil {
			return fmt.Errorf("%T write field begin error 1:data: %s", p, err)
		}
		if err := p.Data.Write(oprot); err != nil {
			return fmt.Errorf("%T error writing struct: %s", p.Data, err)
		}
		if err := oprot.WriteFieldEnd(); err != nil {
			return fmt.Errorf("%T write field end error 1:data: %s", p, err)
		}
	}
	return err
}

func (p *SayHelloRes) writeField2(oprot thrift.TProtocol) (err error) {
	if p.IsSetErrInfo() {
		if err := oprot.WriteFieldBegin("errInfo", thrift.STRUCT, 2); err != nil {
			return fmt.Errorf("%T write field begin error 2:errInfo: %s", p, err)
		}
		if err := p.ErrInfo.Write(oprot); err != nil {
			return fmt.Errorf("%T error writing struct: %s", p.ErrInfo, err)
		}
		if err := oprot.WriteFieldEnd(); err != nil {
			return fmt.Errorf("%T write field end error 2:errInfo: %s", p, err)
		}
	}
	return err
}

func (p *SayHelloRes) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("SayHelloRes(%+v)", *p)
}
