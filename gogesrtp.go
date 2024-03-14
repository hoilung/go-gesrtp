package gogesrtp

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"time"
)

type GoGESRTP struct {
	IPAddress string
	Port      int
	Timeout   int
	conn      net.Conn
}

func NewGeSrtp(ipaddress string, port int, timeout int) GoGESRTP {
	client = GoGESRTP{IPAddress: ipaddress, Port: port, Timeout: timeout, conn: nil}
	return client
}

var cmd_init = [56]byte{}
var cmd_read = []byte{
	2, 0, 0, 0, 0, 0, 0, 0, 0, 1,
	0, 0, 0, 0, 0, 0, 0, 1, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	6, 192, 0, 0, 0, 0, 16, 14, 0, 0,
	1, 1, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0,
}
var cmd_write = []byte{
	2, 0, 0, 0, 0, 0, 0, 0, 0, 2,
	0, 0, 0, 0, 0, 0, 0, 2, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	9, 128, 0, 0, 0, 0, 16, 14, 0, 0,
	1, 1, 2, 0, 0, 0, 0, 0, 1, 1,
	7, 0, 0, 0, 0, 0,
}
var client GoGESRTP

func (GoGESRTP) Open() bool {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", client.IPAddress, client.Port), time.Second*time.Duration(client.Timeout))
	if err != nil {
		return false
	}
	client.conn = conn
	_, err = client.conn.Write(cmd_init[:])
	if err != nil {
		return false
	}
	var resp = make([]byte, 56)
	_, err = client.conn.Read(resp)
	if err != nil {
		return false
	}
	return resp[0] == 0x01
}

func (GoGESRTP) Close() {
	if client.conn != nil {
		client.conn.Close()
	}
	client.conn = nil
}

func (GoGESRTP) ReadBoolean(address string) (bool, error) {
	item, err := ParseDataItemForm(address, 1, true)
	if err != nil {
		return false, err
	}
	resp, err := client.ReadBooleanArray(item.Address, 1)
	if err != nil {
		return false, err
	}
	return resp[0], nil
}

func (GoGESRTP) ReadBooleanArray(address string, valLen uint16) ([]bool, error) {
	item, err := ParseDataItemForm(address, valLen, true)
	if err != nil {
		return nil, err
	}
	resp, err := client.Read(item.DataType, item.StartAddress, valLen)
	if err != nil {
		return nil, err
	}
	array := ConvertByteToBoolArray(resp)
	return array[item.BitAddress : item.BitAddress+valLen], nil
}

func ConvertByteToBoolArray(values []byte) []bool {

	var boolArray = make([]bool, len(values)*8)
	for b := range boolArray {
		boolArray[b] = values[b/8]&(1<<uint(b)) != 0
	}
	return boolArray
}

func (GoGESRTP) Read(datatype byte, address uint16, byteLen uint16) ([]byte, error) {
	resp, err := sendCommand(buildReadCommand(datatype, address, byteLen))
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (GoGESRTP) Write(dataType byte, startAddress uint16, length uint16, value []byte) bool {
	resp, err := sendCommand(buildWriteCommand(dataType, startAddress, length, value))
	if err != nil {
		return false
	}
	return resp[0] == 0x03
}

func sendCommand(command []byte) ([]byte, error) {
	_, err := client.conn.Write(command)
	if err != nil {
		return nil, err
	}
	//fmt.Println(n)
	handData := make([]byte, 56)
	_, err = client.conn.Read(handData)
	if err != nil {
		return nil, err
	}
	dataLen := binary.LittleEndian.Uint16(handData[4:6])
	data := make([]byte, dataLen)
	if dataLen > 0 {
		_, _ = client.conn.Read(data)
	}
	if handData[0] != 0x03 {
		return nil, errors.New("request failure")
	}
	if handData[31] == 212 {
		if binary.LittleEndian.Uint16(handData[42:44]) > 0 {
			return nil, errors.New("31,212 error")
		}
		return handData[44:50], nil
	}
	if handData[31] == 148 {
		return data, nil
	}
	return nil, errors.New("request failure")
}

func buildReadCommand(dataType byte, address uint16, byteLength uint16) []byte {
	valNum := byteLength
	if dataType == R || dataType == AI || dataType == AQ {
		valNum /= 2
	}
	cmd := make([]byte, 56)
	copy(cmd, cmd_read)
	cmd[42] = 0x04 //READ_SYS_MEMORY
	cmd[43] = dataType
	binary.LittleEndian.PutUint16(cmd[44:46], uint16(address))
	binary.LittleEndian.PutUint16(cmd[46:48], uint16(valNum))
	return cmd
}

func buildWriteCommand(dataType byte, startAddress uint16, length uint16, value []byte) []byte {
	cmd := make([]byte, 56+len(value))
	copy(cmd, cmd_write)
	//2,3  id
	binary.LittleEndian.PutUint16(cmd[4:6], uint16(length))
	cmd[51] = dataType
	binary.LittleEndian.PutUint16(cmd[52:54], uint16(startAddress-1))
	binary.LittleEndian.PutUint16(cmd[54:56], uint16(length))
	copy(cmd[56:], value)
	return cmd
}
