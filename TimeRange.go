package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type TimeRange struct {
	BeginTime string
	EndTime   string
	Set       bool
}

func (t TimeRange) InTime() bool {
	swapped := false
	// in time if not set
	if !t.Set {
		return true
	}

	now := time.Now()

	ts, err := time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", time.Now().Format("2006-01-02"), t.BeginTime), now.Location())
	if err != nil {
		panic(err)
	}

	te, err := time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", time.Now().Format("2006-01-02"), t.EndTime), now.Location())
	if err != nil {
		panic(err)
	}

	if te.Before(ts) {
		//te = te.Add(time.Hour * 24)
		swapped = true
	}

	if !swapped {
		if now.After(ts) && now.Before(te) {
			return true
		}
	}

	if swapped {
		if now.Before(te) {
			return true
		}

		te = te.Add(time.Hour * 24)
		if now.After(ts) && now.Before(te) {
			return true
		}
	}

	return false
}

func (t TimeRange) MarshalJSON() ([]byte, error) {
	return json.Marshal(fmt.Sprintf("%s-%s", t.BeginTime, t.EndTime))
}

func (t *TimeRange) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), "\"")

	if len(s) == 0 {
		return nil
	}
	if s == "-" {
		return nil
	}

	ss := strings.Split(s, "-")
	if len(ss) != 2 {
		return fmt.Errorf("time_range - wrong format")
	}

	// check is correct
	_, err := time.Parse("15:04", ss[0])
	if err != nil {
		return fmt.Errorf("time_range - wrong format")
	}

	_, err = time.Parse("15:04", ss[1])
	if err != nil {
		return fmt.Errorf("time_range - wrong format")
	}

	t.BeginTime = ss[0]
	t.EndTime = ss[1]
	t.Set = true

	return nil
}
