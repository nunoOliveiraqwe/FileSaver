package main

import (
	"fmt"
	"github.com/beevik/ntp"
	"os"
	"time"
)


func GetNtpTime() string{
	ntpTime,err := ntp.Time("0.beevik-ntp.pool.ntp.org")
	if err != nil {
		fmt.Fprint(os.Stderr,err.Error())
		ntpTime,err := ntp.Time("pool.ntp.org")
		if err != nil {
			fmt.Fprintf(os.Stderr,err.Error())
			return string(time.Now().Local().Format(time.RFC850))
		} else{
			return ntpTime.Format(time.RFC850)
		}
	} else{
		return ntpTime.Format(time.RFC850)
	}
}
