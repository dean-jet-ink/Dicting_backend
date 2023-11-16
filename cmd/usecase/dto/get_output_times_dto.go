package dto

import "time"

type GetOutputTimesOutput struct {
	OutputTimes []*time.Time `json:"output_times"`
}
