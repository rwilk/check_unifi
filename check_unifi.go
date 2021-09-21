/*-----------------------------------------------------------------------------
# Name:        CHECK_UNIFI 0.1.2
# Purpose:     Nagios/Icinga checker for UniFi Controller condition
#
# Author:      Rafal Wilk <rw@pcboot.pl>
#
# Created:     24-06-2021
# Modified:    21-09-2021
# Copyright:   (c) PcBoot 2021
# License:     BSD-new
-----------------------------------------------------------------------------*/

package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/alexflint/go-arg"
	"net/http"
	"os"
	"time"
)

var ctxmain = context.Background()
var conf Configuration
var itsOK = true
var failed = []DeviceBasic{}

var args struct {
	ConfFile string `arg:"-c,--config,required" help:"config file"`
}

func main() {
	if err := arg.Parse(&args); err != nil {
		fmt.Println("CHECK_UNIFI 0.1.2 for UniFi Controller")
		fmt.Println("All rights reserved. (c) PcBoot 2021")
		fmt.Println()
		arg.MustParse(&args)
	}

	if err := conf.Load(args.ConfFile); err != nil {
		handlePanic(err)
	}

	if !conf.TimeRange.InTime() {
		fmt.Println("UniFi Controller status - OK (skipped)")
		os.Exit(0)
	}

	if conf.SkipSSLVerify {
		// disable SSL cert checking
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	unifi := &APIUnifi{}

	ctx, _ := context.WithTimeout(ctxmain, time.Second*time.Duration(conf.Timeout))

	err := unifi.Login(ctx, conf.ControllerURL, conf.Username, conf.Password)
	if err != nil {
		handlePanic(err)
	}

	err, devices := unifi.GetDeviceBasic(ctx, conf.Site)
	if err != nil {
		handlePanic(err)
	}

	for _, d := range devices.Data {
		if !d.IsOK() {
			itsOK = false
			failed = append(failed, d)
		}
	}

	if itsOK {
		fmt.Println("UniFi Controller status - OK")
	} else {
		fmt.Println("UniFi Controller status - Failed")
		fmt.Println()
		for _, df := range failed {
			fmt.Println(df.String())
		}

		os.Exit(2)
	}
}

func handlePanic(err error) {
	fmt.Println("UniFi Controller status - Failed")
	fmt.Println()
	panic(err)

}
