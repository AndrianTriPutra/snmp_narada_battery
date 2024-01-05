package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"atp/snmp/narada_battery/pkg/repo/snmp"
	"atp/snmp/narada_battery/pkg/utils/domain"
)

func main() {
	flag.Usage = func() {
		log.Printf("Usage: go run . ip pack_id")
		flag.PrintDefaults()
	}

	flag.Parse()

	if len(flag.Args()) != 2 {
		flag.Usage()
		os.Exit(1)
	}

	ip := flag.Args()[0]
	pack := flag.Args()[1]
	pack_id, err := strconv.Atoi(pack)
	if err != nil {
		log.Fatalf("failed parse pack -> %s", err.Error())
	}

	setting := snmp.Setting{
		Port:      161,
		Community: "public",
		Timeout:   300 * time.Millisecond,
	}
	sNMp := snmp.NewRepository(setting)
	ctx := context.Background()

	loc, _ := time.LoadLocation("Asia/Jakarta")
	ts := time.Now().In(loc).Format(time.RFC3339)
	payload := domain.Payload{
		Device_ID: "dev_001",
		Timestamp: ts,
	}

	data, err := sNMp.Narada(ctx, ip, pack_id)
	if err != nil {
		newErr := fmt.Sprintf("failed sNMp.Narada ip [%s] pack [%v]-> %s", ip, pack_id, err.Error())
		log.Fatal(newErr)
	}
	payload.Battery = append(payload.Battery, data)

	js, err := json.MarshalIndent(payload, " ", " ")
	if err != nil {
		log.Fatalf("failed MarshalIndent -> %s", err.Error())
	}
	msg := string(js)
	fmt.Printf("payload:\n%s\n", msg)
}
