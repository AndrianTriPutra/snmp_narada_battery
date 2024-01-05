package snmp

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	g "github.com/gosnmp/gosnmp"

	"atp/snmp/narada_battery/pkg/utils/domain"
)

type strucT struct {
	setting Setting
}

type Setting struct {
	Port      uint16
	Community string
	Timeout   time.Duration
}

func NewRepository(setting Setting) RepositoryI {
	return &strucT{
		setting: setting,
	}
}

type RepositoryI interface {
	Narada(ctx context.Context, ip string, pack int) (conclusion domain.Battery, err error)
}

func (m strucT) Narada(ctx context.Context, ip string, pack int) (conclusion domain.Battery, err error) {
	var data domain.DataBattery
	conclusion.Pack_ID = uint8(pack)

	params := &g.GoSNMP{
		Target:    ip,
		Port:      m.setting.Port,
		Community: m.setting.Community,
		Version:   g.Version2c,
		Timeout:   m.setting.Timeout,
	}
	err = params.Connect()
	if err != nil {
		errN := errors.New("failed SNMP connect")
		return conclusion, errN
	}
	defer params.Conn.Close()

	prefix := ".1.3.6.1.4.1.51232.70."

	var suffix = [30]string{
		".1.0", // current
		".2.0", // voltage
		".3.0", // soc

		".9.0",  // vcell-1
		".10.0", // vcell-2
		".11.0", // vcell-3
		".12.0", // vcell-4
		".13.0", // vcell-5
		".14.0", // vcell-6
		".15.0", // vcell-7
		".16.0", // vcell-8
		".17.0", // vcell-9
		".18.0", // vcell-10
		".19.0", // vcell-11
		".20.0", // vcell-12
		".21.0", // vcell-13
		".22.0", // vcell-14
		".23.0", // vcell-15 22
	}

	var oids = []string{}
	for i := 0; i < 30; i++ {
		buffoid := prefix + strconv.Itoa(pack) + suffix[i]
		oids = append(oids, buffoid)
	}

	for i, oid := range oids {
		var myoid []string
		myoid = append(myoid, oid)

		result, err := params.Get(myoid)
		if err != nil {
			errN := errors.New("failed SNMP GET " + oid)
			return conclusion, errN
		}

		// X Balance
		for _, value := range result.Variables {
			switch value.Type {
			case g.OctetString:
				sval := string(value.Value.([]byte))

				fval, _ := strconv.ParseFloat(sval, 64)
				fmt.Printf("[%v] %s-> %.2f\n", i, oid, fval)
				switch i {
				case 0:
					data.Current = float32(fval) / 100.0
				case 1:
					data.Voltage = float32(fval) / 100.0
				case 2:
					data.Soc = float32(fval)
				case 3:
					data.Cells.Cell1 = float32(fval) / 1000.0
				case 4:
					data.Cells.Cell2 = float32(fval) / 1000.0
				case 5:
					data.Cells.Cell3 = float32(fval) / 1000.0
				case 6:
					data.Cells.Cell4 = float32(fval) / 1000.0
				case 7:
					data.Cells.Cell5 = float32(fval) / 1000.0
				case 8:
					data.Cells.Cell6 = float32(fval) / 1000.0
				case 9:
					data.Cells.Cell7 = float32(fval) / 1000.0
				case 10:
					data.Cells.Cell8 = float32(fval) / 1000.0
				case 11:
					data.Cells.Cell9 = float32(fval) / 1000.0
				case 12:
					data.Cells.Cell10 = float32(fval) / 1000.0
				case 13:
					data.Cells.Cell11 = float32(fval) / 1000.0
				case 14:
					data.Cells.Cell12 = float32(fval) / 1000.0
				case 15:
					data.Cells.Cell13 = float32(fval) / 1000.0
				case 16:
					data.Cells.Cell14 = float32(fval) / 1000.0
				case 17:
					data.Cells.Cell15 = float32(fval) / 1000.0

				}
			}
		}

	}

	conclusion.Data = data
	return conclusion, nil
}
