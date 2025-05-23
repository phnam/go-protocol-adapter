// Code generated by Thrift Compiler (0.21.0). DO NOT EDIT.

package thriftapi

import (
	"bytes"
	"context"
	"database/sql/driver"
	"errors"
	"fmt"
	thrift "github.com/apache/thrift/lib/go/thrift"
	"log/slog"
	"regexp"
	"strings"
	"time"
)

// (needed to ensure safety because of naive import list construction.)
var _ = bytes.Equal
var _ = context.Background
var _ = errors.New
var _ = fmt.Printf
var _ = slog.Log
var _ = time.Now
var _ = thrift.ZERO
// (needed by validator.)
var _ = strings.Contains
var _ = regexp.MatchString

type Status int64
const (
	Status_OK           Status = 200
	Status_INVALID      Status = 400
	Status_UNAUTHORIZED Status = 401
	Status_FORBIDDEN    Status = 403
	Status_NOT_FOUND    Status = 404
	Status_EXISTED      Status = 409
	Status_ERROR        Status = 500
	Status_REDIRECTED   Status = 302
)

func (p Status) String() string {
	switch p {
	case Status_OK: return "OK"
	case Status_INVALID: return "INVALID"
	case Status_UNAUTHORIZED: return "UNAUTHORIZED"
	case Status_FORBIDDEN: return "FORBIDDEN"
	case Status_NOT_FOUND: return "NOT_FOUND"
	case Status_EXISTED: return "EXISTED"
	case Status_ERROR: return "ERROR"
	case Status_REDIRECTED: return "REDIRECTED"
	}
	return "<UNSET>"
}

func StatusFromString(s string) (Status, error) {
	switch s {
	case "OK": return Status_OK, nil
	case "INVALID": return Status_INVALID, nil
	case "UNAUTHORIZED": return Status_UNAUTHORIZED, nil
	case "FORBIDDEN": return Status_FORBIDDEN, nil
	case "NOT_FOUND": return Status_NOT_FOUND, nil
	case "EXISTED": return Status_EXISTED, nil
	case "ERROR": return Status_ERROR, nil
	case "REDIRECTED": return Status_REDIRECTED, nil
	}
	return Status(0), fmt.Errorf("not a valid Status string")
}


func StatusPtr(v Status) *Status { return &v }

func (p Status) MarshalText() ([]byte, error) {
	return []byte(p.String()), nil
}

func (p *Status) UnmarshalText(text []byte) error {
	q, err := StatusFromString(string(text))
	if err != nil {
		return err
	}
	*p = q
	return nil
}

func (p *Status) Scan(value interface{}) error {
	v, ok := value.(int64)
	if !ok {
		return errors.New("Scan value is not int64")
	}
	*p = Status(v)
	return nil
}

func (p *Status) Value() (driver.Value, error) {
	if p == nil {
		return nil, nil
	}
	return int64(*p), nil
}

// Attributes:
//  - Path
//  - Method
//  - Content
//  - Params
//  - Headers
// 
type APIRequest struct {
	Path string `thrift:"path,1" db:"path" json:"path"`
	Method string `thrift:"method,2" db:"method" json:"method"`
	Content string `thrift:"content,3" db:"content" json:"content"`
	Params map[string]string `thrift:"params,4" db:"params" json:"params"`
	Headers map[string]string `thrift:"headers,5" db:"headers" json:"headers"`
}

func NewAPIRequest() *APIRequest {
	return &APIRequest{}
}



func (p *APIRequest) GetPath() string {
	return p.Path
}



func (p *APIRequest) GetMethod() string {
	return p.Method
}



func (p *APIRequest) GetContent() string {
	return p.Content
}



func (p *APIRequest) GetParams() map[string]string {
	return p.Params
}



func (p *APIRequest) GetHeaders() map[string]string {
	return p.Headers
}

func (p *APIRequest) Read(ctx context.Context, iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(ctx); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
	}


	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin(ctx)
		if err != nil {
			return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 1:
			if fieldTypeId == thrift.STRING {
				if err := p.ReadField1(ctx, iprot); err != nil {
					return err
				}
			} else {
				if err := iprot.Skip(ctx, fieldTypeId); err != nil {
					return err
				}
			}
		case 2:
			if fieldTypeId == thrift.STRING {
				if err := p.ReadField2(ctx, iprot); err != nil {
					return err
				}
			} else {
				if err := iprot.Skip(ctx, fieldTypeId); err != nil {
					return err
				}
			}
		case 3:
			if fieldTypeId == thrift.STRING {
				if err := p.ReadField3(ctx, iprot); err != nil {
					return err
				}
			} else {
				if err := iprot.Skip(ctx, fieldTypeId); err != nil {
					return err
				}
			}
		case 4:
			if fieldTypeId == thrift.MAP {
				if err := p.ReadField4(ctx, iprot); err != nil {
					return err
				}
			} else {
				if err := iprot.Skip(ctx, fieldTypeId); err != nil {
					return err
				}
			}
		case 5:
			if fieldTypeId == thrift.MAP {
				if err := p.ReadField5(ctx, iprot); err != nil {
					return err
				}
			} else {
				if err := iprot.Skip(ctx, fieldTypeId); err != nil {
					return err
				}
			}
		default:
			if err := iprot.Skip(ctx, fieldTypeId); err != nil {
				return err
			}
		}
		if err := iprot.ReadFieldEnd(ctx); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(ctx); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
	}
	return nil
}

func (p *APIRequest) ReadField1(ctx context.Context, iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(ctx); err != nil {
		return thrift.PrependError("error reading field 1: ", err)
	} else {
		p.Path = v
	}
	return nil
}

func (p *APIRequest) ReadField2(ctx context.Context, iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(ctx); err != nil {
		return thrift.PrependError("error reading field 2: ", err)
	} else {
		p.Method = v
	}
	return nil
}

func (p *APIRequest) ReadField3(ctx context.Context, iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(ctx); err != nil {
		return thrift.PrependError("error reading field 3: ", err)
	} else {
		p.Content = v
	}
	return nil
}

func (p *APIRequest) ReadField4(ctx context.Context, iprot thrift.TProtocol) error {
	_, _, size, err := iprot.ReadMapBegin(ctx)
	if err != nil {
		return thrift.PrependError("error reading map begin: ", err)
	}
	tMap := make(map[string]string, size)
	p.Params = tMap
	for i := 0; i < size; i++ {
		var _key0 string
		if v, err := iprot.ReadString(ctx); err != nil {
			return thrift.PrependError("error reading field 0: ", err)
		} else {
			_key0 = v
		}
		var _val1 string
		if v, err := iprot.ReadString(ctx); err != nil {
			return thrift.PrependError("error reading field 0: ", err)
		} else {
			_val1 = v
		}
		p.Params[_key0] = _val1
	}
	if err := iprot.ReadMapEnd(ctx); err != nil {
		return thrift.PrependError("error reading map end: ", err)
	}
	return nil
}

func (p *APIRequest) ReadField5(ctx context.Context, iprot thrift.TProtocol) error {
	_, _, size, err := iprot.ReadMapBegin(ctx)
	if err != nil {
		return thrift.PrependError("error reading map begin: ", err)
	}
	tMap := make(map[string]string, size)
	p.Headers = tMap
	for i := 0; i < size; i++ {
		var _key2 string
		if v, err := iprot.ReadString(ctx); err != nil {
			return thrift.PrependError("error reading field 0: ", err)
		} else {
			_key2 = v
		}
		var _val3 string
		if v, err := iprot.ReadString(ctx); err != nil {
			return thrift.PrependError("error reading field 0: ", err)
		} else {
			_val3 = v
		}
		p.Headers[_key2] = _val3
	}
	if err := iprot.ReadMapEnd(ctx); err != nil {
		return thrift.PrependError("error reading map end: ", err)
	}
	return nil
}

func (p *APIRequest) Write(ctx context.Context, oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin(ctx, "APIRequest"); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
	}
	if p != nil {
		if err := p.writeField1(ctx, oprot); err != nil { return err }
		if err := p.writeField2(ctx, oprot); err != nil { return err }
		if err := p.writeField3(ctx, oprot); err != nil { return err }
		if err := p.writeField4(ctx, oprot); err != nil { return err }
		if err := p.writeField5(ctx, oprot); err != nil { return err }
	}
	if err := oprot.WriteFieldStop(ctx); err != nil {
		return thrift.PrependError("write field stop error: ", err)
	}
	if err := oprot.WriteStructEnd(ctx); err != nil {
		return thrift.PrependError("write struct stop error: ", err)
	}
	return nil
}

func (p *APIRequest) writeField1(ctx context.Context, oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin(ctx, "path", thrift.STRING, 1); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:path: ", p), err)
	}
	if err := oprot.WriteString(ctx, string(p.Path)); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.path (1) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(ctx); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 1:path: ", p), err)
	}
	return err
}

func (p *APIRequest) writeField2(ctx context.Context, oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin(ctx, "method", thrift.STRING, 2); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 2:method: ", p), err)
	}
	if err := oprot.WriteString(ctx, string(p.Method)); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.method (2) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(ctx); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 2:method: ", p), err)
	}
	return err
}

func (p *APIRequest) writeField3(ctx context.Context, oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin(ctx, "content", thrift.STRING, 3); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 3:content: ", p), err)
	}
	if err := oprot.WriteString(ctx, string(p.Content)); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.content (3) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(ctx); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 3:content: ", p), err)
	}
	return err
}

func (p *APIRequest) writeField4(ctx context.Context, oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin(ctx, "params", thrift.MAP, 4); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 4:params: ", p), err)
	}
	if err := oprot.WriteMapBegin(ctx, thrift.STRING, thrift.STRING, len(p.Params)); err != nil {
		return thrift.PrependError("error writing map begin: ", err)
	}
	for k, v := range p.Params {
		if err := oprot.WriteString(ctx, string(k)); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T. (0) field write error: ", p), err)
		}
		if err := oprot.WriteString(ctx, string(v)); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T. (0) field write error: ", p), err)
		}
	}
	if err := oprot.WriteMapEnd(ctx); err != nil {
		return thrift.PrependError("error writing map end: ", err)
	}
	if err := oprot.WriteFieldEnd(ctx); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 4:params: ", p), err)
	}
	return err
}

func (p *APIRequest) writeField5(ctx context.Context, oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin(ctx, "headers", thrift.MAP, 5); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 5:headers: ", p), err)
	}
	if err := oprot.WriteMapBegin(ctx, thrift.STRING, thrift.STRING, len(p.Headers)); err != nil {
		return thrift.PrependError("error writing map begin: ", err)
	}
	for k, v := range p.Headers {
		if err := oprot.WriteString(ctx, string(k)); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T. (0) field write error: ", p), err)
		}
		if err := oprot.WriteString(ctx, string(v)); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T. (0) field write error: ", p), err)
		}
	}
	if err := oprot.WriteMapEnd(ctx); err != nil {
		return thrift.PrependError("error writing map end: ", err)
	}
	if err := oprot.WriteFieldEnd(ctx); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 5:headers: ", p), err)
	}
	return err
}

func (p *APIRequest) Equals(other *APIRequest) bool {
	if p == other {
		return true
	} else if p == nil || other == nil {
		return false
	}
	if p.Path != other.Path { return false }
	if p.Method != other.Method { return false }
	if p.Content != other.Content { return false }
	if len(p.Params) != len(other.Params) { return false }
	for k, _tgt := range p.Params {
		_src4 := other.Params[k]
		if _tgt != _src4 { return false }
	}
	if len(p.Headers) != len(other.Headers) { return false }
	for k, _tgt := range p.Headers {
		_src5 := other.Headers[k]
		if _tgt != _src5 { return false }
	}
	return true
}

func (p *APIRequest) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("APIRequest(%+v)", *p)
}

func (p *APIRequest) LogValue() slog.Value {
	if p == nil {
		return slog.AnyValue(nil)
	}
	v := thrift.SlogTStructWrapper{
		Type: "*api.APIRequest",
		Value: p,
	}
	return slog.AnyValue(v)
}

var _ slog.LogValuer = (*APIRequest)(nil)

func (p *APIRequest) Validate() error {
	return nil
}

// Attributes:
//  - Status
//  - Message
//  - Headers
//  - Content
//  - Total
//  - ErrorCode
// 
type APIResponse struct {
	Status  Status `thrift:"status,1" db:"status" json:"status"`
	Message string `thrift:"message,2" db:"message" json:"message"`
	Headers map[string]string `thrift:"headers,3" db:"headers" json:"headers"`
	Content string `thrift:"content,4" db:"content" json:"content"`
	Total int64 `thrift:"total,5" db:"total" json:"total"`
	ErrorCode string `thrift:"errorCode,6" db:"errorCode" json:"errorCode"`
}

func NewAPIResponse() *APIResponse {
	return &APIResponse{}
}



func (p *APIResponse) GetStatus() Status {
	return p.Status
}



func (p *APIResponse) GetMessage() string {
	return p.Message
}



func (p *APIResponse) GetHeaders() map[string]string {
	return p.Headers
}



func (p *APIResponse) GetContent() string {
	return p.Content
}



func (p *APIResponse) GetTotal() int64 {
	return p.Total
}



func (p *APIResponse) GetErrorCode() string {
	return p.ErrorCode
}

func (p *APIResponse) Read(ctx context.Context, iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(ctx); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
	}


	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin(ctx)
		if err != nil {
			return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 1:
			if fieldTypeId == thrift.I32 {
				if err := p.ReadField1(ctx, iprot); err != nil {
					return err
				}
			} else {
				if err := iprot.Skip(ctx, fieldTypeId); err != nil {
					return err
				}
			}
		case 2:
			if fieldTypeId == thrift.STRING {
				if err := p.ReadField2(ctx, iprot); err != nil {
					return err
				}
			} else {
				if err := iprot.Skip(ctx, fieldTypeId); err != nil {
					return err
				}
			}
		case 3:
			if fieldTypeId == thrift.MAP {
				if err := p.ReadField3(ctx, iprot); err != nil {
					return err
				}
			} else {
				if err := iprot.Skip(ctx, fieldTypeId); err != nil {
					return err
				}
			}
		case 4:
			if fieldTypeId == thrift.STRING {
				if err := p.ReadField4(ctx, iprot); err != nil {
					return err
				}
			} else {
				if err := iprot.Skip(ctx, fieldTypeId); err != nil {
					return err
				}
			}
		case 5:
			if fieldTypeId == thrift.I64 {
				if err := p.ReadField5(ctx, iprot); err != nil {
					return err
				}
			} else {
				if err := iprot.Skip(ctx, fieldTypeId); err != nil {
					return err
				}
			}
		case 6:
			if fieldTypeId == thrift.STRING {
				if err := p.ReadField6(ctx, iprot); err != nil {
					return err
				}
			} else {
				if err := iprot.Skip(ctx, fieldTypeId); err != nil {
					return err
				}
			}
		default:
			if err := iprot.Skip(ctx, fieldTypeId); err != nil {
				return err
			}
		}
		if err := iprot.ReadFieldEnd(ctx); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(ctx); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
	}
	return nil
}

func (p *APIResponse) ReadField1(ctx context.Context, iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI32(ctx); err != nil {
		return thrift.PrependError("error reading field 1: ", err)
	} else {
		temp := Status(v)
		p.Status = temp
	}
	return nil
}

func (p *APIResponse) ReadField2(ctx context.Context, iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(ctx); err != nil {
		return thrift.PrependError("error reading field 2: ", err)
	} else {
		p.Message = v
	}
	return nil
}

func (p *APIResponse) ReadField3(ctx context.Context, iprot thrift.TProtocol) error {
	_, _, size, err := iprot.ReadMapBegin(ctx)
	if err != nil {
		return thrift.PrependError("error reading map begin: ", err)
	}
	tMap := make(map[string]string, size)
	p.Headers = tMap
	for i := 0; i < size; i++ {
		var _key6 string
		if v, err := iprot.ReadString(ctx); err != nil {
			return thrift.PrependError("error reading field 0: ", err)
		} else {
			_key6 = v
		}
		var _val7 string
		if v, err := iprot.ReadString(ctx); err != nil {
			return thrift.PrependError("error reading field 0: ", err)
		} else {
			_val7 = v
		}
		p.Headers[_key6] = _val7
	}
	if err := iprot.ReadMapEnd(ctx); err != nil {
		return thrift.PrependError("error reading map end: ", err)
	}
	return nil
}

func (p *APIResponse) ReadField4(ctx context.Context, iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(ctx); err != nil {
		return thrift.PrependError("error reading field 4: ", err)
	} else {
		p.Content = v
	}
	return nil
}

func (p *APIResponse) ReadField5(ctx context.Context, iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI64(ctx); err != nil {
		return thrift.PrependError("error reading field 5: ", err)
	} else {
		p.Total = v
	}
	return nil
}

func (p *APIResponse) ReadField6(ctx context.Context, iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(ctx); err != nil {
		return thrift.PrependError("error reading field 6: ", err)
	} else {
		p.ErrorCode = v
	}
	return nil
}

func (p *APIResponse) Write(ctx context.Context, oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin(ctx, "APIResponse"); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
	}
	if p != nil {
		if err := p.writeField1(ctx, oprot); err != nil { return err }
		if err := p.writeField2(ctx, oprot); err != nil { return err }
		if err := p.writeField3(ctx, oprot); err != nil { return err }
		if err := p.writeField4(ctx, oprot); err != nil { return err }
		if err := p.writeField5(ctx, oprot); err != nil { return err }
		if err := p.writeField6(ctx, oprot); err != nil { return err }
	}
	if err := oprot.WriteFieldStop(ctx); err != nil {
		return thrift.PrependError("write field stop error: ", err)
	}
	if err := oprot.WriteStructEnd(ctx); err != nil {
		return thrift.PrependError("write struct stop error: ", err)
	}
	return nil
}

func (p *APIResponse) writeField1(ctx context.Context, oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin(ctx, "status", thrift.I32, 1); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:status: ", p), err)
	}
	if err := oprot.WriteI32(ctx, int32(p.Status)); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.status (1) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(ctx); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 1:status: ", p), err)
	}
	return err
}

func (p *APIResponse) writeField2(ctx context.Context, oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin(ctx, "message", thrift.STRING, 2); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 2:message: ", p), err)
	}
	if err := oprot.WriteString(ctx, string(p.Message)); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.message (2) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(ctx); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 2:message: ", p), err)
	}
	return err
}

func (p *APIResponse) writeField3(ctx context.Context, oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin(ctx, "headers", thrift.MAP, 3); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 3:headers: ", p), err)
	}
	if err := oprot.WriteMapBegin(ctx, thrift.STRING, thrift.STRING, len(p.Headers)); err != nil {
		return thrift.PrependError("error writing map begin: ", err)
	}
	for k, v := range p.Headers {
		if err := oprot.WriteString(ctx, string(k)); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T. (0) field write error: ", p), err)
		}
		if err := oprot.WriteString(ctx, string(v)); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T. (0) field write error: ", p), err)
		}
	}
	if err := oprot.WriteMapEnd(ctx); err != nil {
		return thrift.PrependError("error writing map end: ", err)
	}
	if err := oprot.WriteFieldEnd(ctx); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 3:headers: ", p), err)
	}
	return err
}

func (p *APIResponse) writeField4(ctx context.Context, oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin(ctx, "content", thrift.STRING, 4); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 4:content: ", p), err)
	}
	if err := oprot.WriteString(ctx, string(p.Content)); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.content (4) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(ctx); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 4:content: ", p), err)
	}
	return err
}

func (p *APIResponse) writeField5(ctx context.Context, oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin(ctx, "total", thrift.I64, 5); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 5:total: ", p), err)
	}
	if err := oprot.WriteI64(ctx, int64(p.Total)); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.total (5) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(ctx); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 5:total: ", p), err)
	}
	return err
}

func (p *APIResponse) writeField6(ctx context.Context, oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin(ctx, "errorCode", thrift.STRING, 6); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 6:errorCode: ", p), err)
	}
	if err := oprot.WriteString(ctx, string(p.ErrorCode)); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.errorCode (6) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(ctx); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 6:errorCode: ", p), err)
	}
	return err
}

func (p *APIResponse) Equals(other *APIResponse) bool {
	if p == other {
		return true
	} else if p == nil || other == nil {
		return false
	}
	if p.Status != other.Status { return false }
	if p.Message != other.Message { return false }
	if len(p.Headers) != len(other.Headers) { return false }
	for k, _tgt := range p.Headers {
		_src8 := other.Headers[k]
		if _tgt != _src8 { return false }
	}
	if p.Content != other.Content { return false }
	if p.Total != other.Total { return false }
	if p.ErrorCode != other.ErrorCode { return false }
	return true
}

func (p *APIResponse) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("APIResponse(%+v)", *p)
}

func (p *APIResponse) LogValue() slog.Value {
	if p == nil {
		return slog.AnyValue(nil)
	}
	v := thrift.SlogTStructWrapper{
		Type: "*api.APIResponse",
		Value: p,
	}
	return slog.AnyValue(v)
}

var _ slog.LogValuer = (*APIResponse)(nil)

func (p *APIResponse) Validate() error {
	return nil
}

type APIService interface {
	// Parameters:
	//  - Request
	// 
	Call(ctx context.Context, request *APIRequest) (_r *APIResponse, _err error)
}

type APIServiceClient struct {
	c thrift.TClient
	meta thrift.ResponseMeta
}

func NewAPIServiceClientFactory(t thrift.TTransport, f thrift.TProtocolFactory) *APIServiceClient {
	return &APIServiceClient{
		c: thrift.NewTStandardClient(f.GetProtocol(t), f.GetProtocol(t)),
	}
}

func NewAPIServiceClientProtocol(t thrift.TTransport, iprot thrift.TProtocol, oprot thrift.TProtocol) *APIServiceClient {
	return &APIServiceClient{
		c: thrift.NewTStandardClient(iprot, oprot),
	}
}

func NewAPIServiceClient(c thrift.TClient) *APIServiceClient {
	return &APIServiceClient{
		c: c,
	}
}

func (p *APIServiceClient) Client_() thrift.TClient {
	return p.c
}

func (p *APIServiceClient) LastResponseMeta_() thrift.ResponseMeta {
	return p.meta
}

func (p *APIServiceClient) SetLastResponseMeta_(meta thrift.ResponseMeta) {
	p.meta = meta
}

// Parameters:
//  - Request
// 
func (p *APIServiceClient) Call(ctx context.Context, request *APIRequest) (_r *APIResponse, _err error) {
	var _args9 APIServiceCallArgs
	_args9.Request = request
	var _result11 APIServiceCallResult
	var _meta10 thrift.ResponseMeta
	_meta10, _err = p.Client_().Call(ctx, "call", &_args9, &_result11)
	p.SetLastResponseMeta_(_meta10)
	if _err != nil {
		return
	}
	if _ret12 := _result11.GetSuccess(); _ret12 != nil {
		return _ret12, nil
	}
	return nil, thrift.NewTApplicationException(thrift.MISSING_RESULT, "call failed: unknown result")
}

type APIServiceProcessor struct {
	processorMap map[string]thrift.TProcessorFunction
	handler      APIService
}

func (p *APIServiceProcessor) AddToProcessorMap(key string, processor thrift.TProcessorFunction) {
	p.processorMap[key] = processor
}

func (p *APIServiceProcessor) GetProcessorFunction(key string) (processor thrift.TProcessorFunction, ok bool) {
	processor, ok = p.processorMap[key]
	return processor, ok
}

func (p *APIServiceProcessor) ProcessorMap() map[string]thrift.TProcessorFunction {
	return p.processorMap
}

func NewAPIServiceProcessor(handler APIService) *APIServiceProcessor {

	self13 := &APIServiceProcessor{handler: handler, processorMap:make(map[string]thrift.TProcessorFunction)}
	self13.processorMap["call"] = &aPIServiceProcessorCall{handler: handler}
	return self13
}

func (p *APIServiceProcessor) Process(ctx context.Context, iprot, oprot thrift.TProtocol) (success bool, err thrift.TException) {
	name, _, seqId, err2 := iprot.ReadMessageBegin(ctx)
	if err2 != nil { return false, thrift.WrapTException(err2) }
	if processor, ok := p.GetProcessorFunction(name); ok {
		return processor.Process(ctx, seqId, iprot, oprot)
	}
	iprot.Skip(ctx, thrift.STRUCT)
	iprot.ReadMessageEnd(ctx)
	x14 := thrift.NewTApplicationException(thrift.UNKNOWN_METHOD, "Unknown function " + name)
	oprot.WriteMessageBegin(ctx, name, thrift.EXCEPTION, seqId)
	x14.Write(ctx, oprot)
	oprot.WriteMessageEnd(ctx)
	oprot.Flush(ctx)
	return false, x14
}

type aPIServiceProcessorCall struct {
	handler APIService
}

func (p *aPIServiceProcessorCall) Process(ctx context.Context, seqId int32, iprot, oprot thrift.TProtocol) (success bool, err thrift.TException) {
	var _write_err15 error
	args := APIServiceCallArgs{}
	if err2 := args.Read(ctx, iprot); err2 != nil {
		iprot.ReadMessageEnd(ctx)
		x := thrift.NewTApplicationException(thrift.PROTOCOL_ERROR, err2.Error())
		oprot.WriteMessageBegin(ctx, "call", thrift.EXCEPTION, seqId)
		x.Write(ctx, oprot)
		oprot.WriteMessageEnd(ctx)
		oprot.Flush(ctx)
		return false, thrift.WrapTException(err2)
	}
	iprot.ReadMessageEnd(ctx)

	tickerCancel := func() {}
	// Start a goroutine to do server side connectivity check.
	if thrift.ServerConnectivityCheckInterval > 0 {
		var cancel context.CancelCauseFunc
		ctx, cancel = context.WithCancelCause(ctx)
		defer cancel(nil)
		var tickerCtx context.Context
		tickerCtx, tickerCancel = context.WithCancel(context.Background())
		defer tickerCancel()
		go func(ctx context.Context, cancel context.CancelCauseFunc) {
			ticker := time.NewTicker(thrift.ServerConnectivityCheckInterval)
			defer ticker.Stop()
			for {
				select {
				case <-ctx.Done():
					return
				case <-ticker.C:
					if !iprot.Transport().IsOpen() {
						cancel(thrift.ErrAbandonRequest)
						return
					}
				}
			}
		}(tickerCtx, cancel)
	}

	result := APIServiceCallResult{}
	if retval, err2 := p.handler.Call(ctx, args.Request); err2 != nil {
		tickerCancel()
		err = thrift.WrapTException(err2)
		if errors.Is(err2, thrift.ErrAbandonRequest) {
			return false, thrift.WrapTException(err2)
		}
		if errors.Is(err2, context.Canceled) {
			if err := context.Cause(ctx); errors.Is(err, thrift.ErrAbandonRequest) {
				return false, thrift.WrapTException(err)
			}
		}
		_exc16 := thrift.NewTApplicationException(thrift.INTERNAL_ERROR, "Internal error processing call: " + err2.Error())
		if err2 := oprot.WriteMessageBegin(ctx, "call", thrift.EXCEPTION, seqId); err2 != nil {
			_write_err15 = thrift.WrapTException(err2)
		}
		if err2 := _exc16.Write(ctx, oprot); _write_err15 == nil && err2 != nil {
			_write_err15 = thrift.WrapTException(err2)
		}
		if err2 := oprot.WriteMessageEnd(ctx); _write_err15 == nil && err2 != nil {
			_write_err15 = thrift.WrapTException(err2)
		}
		if err2 := oprot.Flush(ctx); _write_err15 == nil && err2 != nil {
			_write_err15 = thrift.WrapTException(err2)
		}
		if _write_err15 != nil {
			return false, thrift.WrapTException(_write_err15)
		}
		return true, err
	} else {
		result.Success = retval
	}
	tickerCancel()
	if err2 := oprot.WriteMessageBegin(ctx, "call", thrift.REPLY, seqId); err2 != nil {
		_write_err15 = thrift.WrapTException(err2)
	}
	if err2 := result.Write(ctx, oprot); _write_err15 == nil && err2 != nil {
		_write_err15 = thrift.WrapTException(err2)
	}
	if err2 := oprot.WriteMessageEnd(ctx); _write_err15 == nil && err2 != nil {
		_write_err15 = thrift.WrapTException(err2)
	}
	if err2 := oprot.Flush(ctx); _write_err15 == nil && err2 != nil {
		_write_err15 = thrift.WrapTException(err2)
	}
	if _write_err15 != nil {
		return false, thrift.WrapTException(_write_err15)
	}
	return true, err
}


// HELPER FUNCTIONS AND STRUCTURES

// Attributes:
//  - Request
// 
type APIServiceCallArgs struct {
	Request *APIRequest `thrift:"request,1" db:"request" json:"request"`
}

func NewAPIServiceCallArgs() *APIServiceCallArgs {
	return &APIServiceCallArgs{}
}

var APIServiceCallArgs_Request_DEFAULT *APIRequest

func (p *APIServiceCallArgs) GetRequest() *APIRequest {
	if !p.IsSetRequest() {
		return APIServiceCallArgs_Request_DEFAULT
	}
	return p.Request
}

func (p *APIServiceCallArgs) IsSetRequest() bool {
	return p.Request != nil
}

func (p *APIServiceCallArgs) Read(ctx context.Context, iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(ctx); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
	}


	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin(ctx)
		if err != nil {
			return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 1:
			if fieldTypeId == thrift.STRUCT {
				if err := p.ReadField1(ctx, iprot); err != nil {
					return err
				}
			} else {
				if err := iprot.Skip(ctx, fieldTypeId); err != nil {
					return err
				}
			}
		default:
			if err := iprot.Skip(ctx, fieldTypeId); err != nil {
				return err
			}
		}
		if err := iprot.ReadFieldEnd(ctx); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(ctx); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
	}
	return nil
}

func (p *APIServiceCallArgs) ReadField1(ctx context.Context, iprot thrift.TProtocol) error {
	p.Request = &APIRequest{}
	if err := p.Request.Read(ctx, iprot); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T error reading struct: ", p.Request), err)
	}
	return nil
}

func (p *APIServiceCallArgs) Write(ctx context.Context, oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin(ctx, "call_args"); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
	}
	if p != nil {
		if err := p.writeField1(ctx, oprot); err != nil { return err }
	}
	if err := oprot.WriteFieldStop(ctx); err != nil {
		return thrift.PrependError("write field stop error: ", err)
	}
	if err := oprot.WriteStructEnd(ctx); err != nil {
		return thrift.PrependError("write struct stop error: ", err)
	}
	return nil
}

func (p *APIServiceCallArgs) writeField1(ctx context.Context, oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin(ctx, "request", thrift.STRUCT, 1); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:request: ", p), err)
	}
	if err := p.Request.Write(ctx, oprot); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T error writing struct: ", p.Request), err)
	}
	if err := oprot.WriteFieldEnd(ctx); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 1:request: ", p), err)
	}
	return err
}

func (p *APIServiceCallArgs) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("APIServiceCallArgs(%+v)", *p)
}

func (p *APIServiceCallArgs) LogValue() slog.Value {
	if p == nil {
		return slog.AnyValue(nil)
	}
	v := thrift.SlogTStructWrapper{
		Type: "*api.APIServiceCallArgs",
		Value: p,
	}
	return slog.AnyValue(v)
}

var _ slog.LogValuer = (*APIServiceCallArgs)(nil)

// Attributes:
//  - Success
// 
type APIServiceCallResult struct {
	Success *APIResponse `thrift:"success,0" db:"success" json:"success,omitempty"`
}

func NewAPIServiceCallResult() *APIServiceCallResult {
	return &APIServiceCallResult{}
}

var APIServiceCallResult_Success_DEFAULT *APIResponse

func (p *APIServiceCallResult) GetSuccess() *APIResponse {
	if !p.IsSetSuccess() {
		return APIServiceCallResult_Success_DEFAULT
	}
	return p.Success
}

func (p *APIServiceCallResult) IsSetSuccess() bool {
	return p.Success != nil
}

func (p *APIServiceCallResult) Read(ctx context.Context, iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(ctx); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
	}


	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin(ctx)
		if err != nil {
			return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 0:
			if fieldTypeId == thrift.STRUCT {
				if err := p.ReadField0(ctx, iprot); err != nil {
					return err
				}
			} else {
				if err := iprot.Skip(ctx, fieldTypeId); err != nil {
					return err
				}
			}
		default:
			if err := iprot.Skip(ctx, fieldTypeId); err != nil {
				return err
			}
		}
		if err := iprot.ReadFieldEnd(ctx); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(ctx); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
	}
	return nil
}

func (p *APIServiceCallResult) ReadField0(ctx context.Context, iprot thrift.TProtocol) error {
	p.Success = &APIResponse{}
	if err := p.Success.Read(ctx, iprot); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T error reading struct: ", p.Success), err)
	}
	return nil
}

func (p *APIServiceCallResult) Write(ctx context.Context, oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin(ctx, "call_result"); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
	}
	if p != nil {
		if err := p.writeField0(ctx, oprot); err != nil { return err }
	}
	if err := oprot.WriteFieldStop(ctx); err != nil {
		return thrift.PrependError("write field stop error: ", err)
	}
	if err := oprot.WriteStructEnd(ctx); err != nil {
		return thrift.PrependError("write struct stop error: ", err)
	}
	return nil
}

func (p *APIServiceCallResult) writeField0(ctx context.Context, oprot thrift.TProtocol) (err error) {
	if p.IsSetSuccess() {
		if err := oprot.WriteFieldBegin(ctx, "success", thrift.STRUCT, 0); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field begin error 0:success: ", p), err)
		}
		if err := p.Success.Write(ctx, oprot); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T error writing struct: ", p.Success), err)
		}
		if err := oprot.WriteFieldEnd(ctx); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field end error 0:success: ", p), err)
		}
	}
	return err
}

func (p *APIServiceCallResult) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("APIServiceCallResult(%+v)", *p)
}

func (p *APIServiceCallResult) LogValue() slog.Value {
	if p == nil {
		return slog.AnyValue(nil)
	}
	v := thrift.SlogTStructWrapper{
		Type: "*api.APIServiceCallResult",
		Value: p,
	}
	return slog.AnyValue(v)
}

var _ slog.LogValuer = (*APIServiceCallResult)(nil)


