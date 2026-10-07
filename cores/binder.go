package cores

import (
	"errors"

	"github.com/gofiber/fiber/v3"
	"github.com/shamaton/msgpack/v3"
)

var msgpackMIMETypes = []string{
	"application/x-msgpack",
	"application/msgpack",
	"application/vnd.msgpack",
}

// MsgPackBinder implements fiber.CustomBinder for MessagePack content types.
type MsgPackBinder struct{}

// NewMsgPackBinder creates a new instance of MsgPackBinder.
func NewMsgPackBinder() *MsgPackBinder {
	return &MsgPackBinder{}
}

// Name returns the binder name.
func (b *MsgPackBinder) Name() string {
	return "msgpack"
}

// MIMETypes returns the supported MessagePack MIME types without dynamic allocations.
func (b *MsgPackBinder) MIMETypes() []string {
	return msgpackMIMETypes
}

// SafeUnmarshal decodes MessagePack data into v using an isolated byte copy.
// This prevents memory corruption when structs contain []byte fields that could
// otherwise reference fasthttp request buffers recycled across connections.
func SafeUnmarshal(data []byte, v any) error {
	if v == nil {
		return errors.New("destination cannot be nil")
	}
	if len(data) == 0 {
		return errors.New("empty request body")
	}
	buf := make([]byte, len(data))
	copy(buf, data)
	return msgpack.Unmarshal(buf, v)
}

// Parse decodes the request body as MessagePack into out.
func (b *MsgPackBinder) Parse(c fiber.Ctx, out any) error {
	return SafeUnmarshal(c.Body(), out)
}
