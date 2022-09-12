package main

import (
	"fmt"

	plc4go "github.com/apache/plc4x/plc4go/pkg/api"
	"github.com/apache/plc4x/plc4go/pkg/api/drivers"
	"github.com/apache/plc4x/plc4go/pkg/api/model"
)

func main() {
	driverManager := plc4go.NewPlcDriverManager()
	drivers.RegisterModbusTcpDriver(driverManager)

	crc := driverManager.GetConnection("modbus-tcp://192.168.23.30")

	connectionResult := <-crc
	if connectionResult.GetErr() != nil {
		fmt.Printf("------Error connecting to PLC: %s", connectionResult.GetErr().Error())
		return
	}
	connection := connectionResult.GetConnection()

	defer connection.BlockingClose()

	readRequest, err := connection.ReadRequestBuilder().
		AddQuery("field", "holding-register:26:REAL").
		Build()
	if err != nil {
		fmt.Printf("------Error preparing read-request: %s", connectionResult.GetErr().Error())
		return
	}

	rrc := readRequest.Execute()

	rrr := <-rrc
	if rrr.GetErr() != nil {
		fmt.Printf("------Error executing read-request: %s", rrr.GetErr().Error())
		return
	}

	if rrr.GetResponse().GetResponseCode("field") != model.PlcResponseCode_OK {
		fmt.Printf("------Error an non-ok return code: %s", rrr.GetResponse().GetResponseCode("field").GetName())
		return
	}

	value := rrr.GetResponse().GetValue("field")
	fmt.Printf("Got result %f", value.GetFloat32())
}