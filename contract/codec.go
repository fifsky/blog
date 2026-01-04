package contract

import (
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"regexp"
	"slices"

	"app/response"

	"github.com/gorilla/schema"
	"google.golang.org/protobuf/encoding/protojson"

	"buf.build/go/protovalidate"
	"github.com/goapt/grpc-http/contract"
	"google.golang.org/protobuf/proto"
)

var _ contract.Codec = (*Codec)(nil)

type Codec struct {
	formDecoder *schema.Decoder
}

func NewCodec() *Codec {
	decoder := schema.NewDecoder()
	decoder.SetAliasTag("json")
	decoder.IgnoreUnknownKeys(true)
	return &Codec{
		formDecoder: decoder,
	}
}

var protodecoder = protojson.UnmarshalOptions{}

func (c *Codec) Encode(w http.ResponseWriter, _ *http.Request, v any) {
	if _, ok := v.(error); ok {
		response.Fail(w, 400, v)
		return
	}
	response.Success(w, v)
}

func (c *Codec) Decode(r *http.Request, v any) error {
	// 如果有URL参数，则解析并赋值给 v
	if r.URL != nil && r.URL.RawQuery != "" {
		q := r.URL.Query()
		if err := c.formDecoder.Decode(v, q); err != nil {
			return fmt.Errorf("decode query unmarshal error: %w", err)
		}
	}

	contentType, _, _ := mime.ParseMediaType(r.Header.Get("Content-Type"))
	// 如果是POST，并且是"x-www-form-urlencoded"，则解析form
	if r.Method == http.MethodPost && contentType == "application/x-www-form-urlencoded" {
		if err := r.ParseForm(); err != nil {
			return fmt.Errorf("parse form error: %w", err)
		}
		if err := c.formDecoder.Decode(v, r.PostForm); err != nil {
			return fmt.Errorf("decode form unmarshal error: %w", err)
		}
	}

	// 如果有路径参数，则解析并赋值给 v
	if r.Pattern != "" {
		names := extractVarNames(r.Pattern)
		if len(names) > 0 {
			vals := url.Values{}
			for _, n := range names {
				vals.Set(n, r.PathValue(n))
			}
			if err := c.formDecoder.Decode(v, vals); err != nil {
				return fmt.Errorf("decode path params error: %w", err)
			}
		}
	}

	if !hasBody(r) {
		return nil
	}

	// 如果不是 JSON 请求体，则跳过，暂不支持其他类型
	if contentType != "application/json" {
		return nil
	}

	if msg, ok := v.(proto.Message); ok {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			return fmt.Errorf("read body error: %w", err)
		}
		if err := protodecoder.Unmarshal(body, msg); err != nil {
			return fmt.Errorf("decode protojson error: %w", err)
		}
	} else {
		if err := json.NewDecoder(r.Body).Decode(v); err != nil {
			return fmt.Errorf("decode body error: %w", err)
		}
	}
	return nil
}

func hasBody(r *http.Request) bool {
	if r.Body == nil || r.Body == http.NoBody {
		return false
	}
	if r.ContentLength > 0 {
		return true
	}
	if slices.Contains(r.TransferEncoding, "chunked") {
		return true
	}
	return false
}

var pathVarReg = regexp.MustCompile(`\{([^{}]+)\}`)

func extractVarNames(pattern string) []string {
	matches := pathVarReg.FindAllStringSubmatch(pattern, -1)
	out := make([]string, 0, len(matches))
	for _, m := range matches {
		out = append(out, m[1])
	}
	return out
}

func (c *Codec) Validate(v any) error {
	if msg, ok := v.(proto.Message); ok {
		if err := protovalidate.Validate(msg); err != nil {
			return err
		}
	}
	// else {
	// 兼容指针与非指针结构体
	// t := reflect.TypeOf(v)
	// if t.Kind() == reflect.Pointer {
	// 	t = t.Elem()
	// }
	// if t.Kind() == reflect.Struct {
	// 	// 传入实际对象（指针或结构体均可），内部 validator 能正确处理
	// 	if err := validate.Validate(v); err != nil {
	// 		return err
	// 	}
	// }
	// }
	return nil
}
