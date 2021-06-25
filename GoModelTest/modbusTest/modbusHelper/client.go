package modbusHelper

import (
	"encoding/binary"
	"fmt"
)

type ClientHandler interface {
	Packager
	Transporter
}

type ProtocolDataUnit struct {
	FunctionCode byte
	Data         []byte
}

type Client interface {
	// Bit access
	ReadHoldingRegisters(address, quantity uint16) (results []byte, err error)

	ReadDiscreteInputs(address, quantity uint16) (results []byte, err error)
}

type Packager interface {
	Encode(pdu *ProtocolDataUnit) (adu []byte, err error)
	Decode(adu []byte) (pdu *ProtocolDataUnit, err error)
	Verify(aduRequest []byte, aduResponse []byte) (err error)
}

// Transporter specifies the transport layer.
type Transporter interface {
	Send(aduRequest []byte) (aduResponse []byte, err error)
}

type client struct {
	packager    Packager
	transporter Transporter
}

func (mb *client) ReadHoldingRegisters(address, quantity uint16) (results []byte, err error) {
	if quantity < 1 || quantity > 125 {
		err = fmt.Errorf("modbus: quantity '%v' must be between '%v' and '%v',", quantity, 1, 125)
		return
	}
	request := ProtocolDataUnit{
		FunctionCode: FuncCodeReadHoldingRegisters,
		Data:         dataBlock(address, quantity),
	}
	response, err := mb.send(&request)
	if err != nil {
		return
	}
	count := int(response.Data[0])
	length := len(response.Data) - 1
	if count != length {
		err = fmt.Errorf("modbus: response data size '%v' does not match count '%v'", length, count)
		return
	}
	results = response.Data[1:]
	return
}

func dataBlock(value ...uint16) []byte {
	data := make([]byte, 2*len(value))
	for i, v := range value {
		binary.BigEndian.PutUint16(data[i*2:], v)
	}
	return data
}

func (mb *client) send(request *ProtocolDataUnit) (response *ProtocolDataUnit, err error) {
	aduRequest, err := mb.packager.Encode(request)
	if err != nil {
		return
	}
	aduResponse, err := mb.transporter.Send(aduRequest)
	if err != nil {
		return
	}
	fmt.Println(aduResponse)
	if err = mb.packager.Verify(aduRequest, aduResponse); err != nil {
		return
	}
	response, err = mb.packager.Decode(aduResponse)
	if err != nil {
		return
	}
	// Check correct function code returned (exception)
	if response.FunctionCode != request.FunctionCode {
		err = responseError(response)
		return
	}
	if response.Data == nil || len(response.Data) == 0 {
		// Empty response
		err = fmt.Errorf("modbus: response data is empty")
		return
	}
	return
}

func responseError(response *ProtocolDataUnit) error {
	mbError := &ModbusError{FunctionCode: response.FunctionCode}
	if response.Data != nil && len(response.Data) > 0 {
		mbError.ExceptionCode = response.Data[0]
	}
	return mbError
}

// NewClient creates a new modbus client with given backend handler.
func NewClient(handler ClientHandler) Client {
	return &client{packager: handler, transporter: handler}
}

func (mb *client) ReadDiscreteInputs(address, quantity uint16) (results []byte, err error) {
	if quantity < 1 || quantity > 2000 {
		err = fmt.Errorf("modbus: quantity '%v' must be between '%v' and '%v',", quantity, 1, 2000)
		return
	}
	request := ProtocolDataUnit{
		FunctionCode: FuncCodeReadDiscreteInputs,
		Data:         dataBlock(address, quantity),
	}
	response, err := mb.send(&request)
	if err != nil {
		return
	}
	count := int(response.Data[0])
	length := len(response.Data) - 1
	if count != length {
		err = fmt.Errorf("modbus: response data size '%v' does not match count '%v'", length, count)
		return
	}
	results = response.Data[1:]
	return
}
