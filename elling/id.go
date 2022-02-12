package elling

import (
	"github.com/Elytrium/elling/config"
	"github.com/rs/zerolog/log"
	"github.com/sony/sonyflake"
	"strconv"
	"time"
)

var ID *sonyflake.Sonyflake

func InitID() {
	var MachineID func() (uint16, error)
	if config.AppConfig.MachineID != "ip" {
		MachineID = func() (uint16, error) {
			id, err := strconv.ParseUint(config.AppConfig.MachineID, 10, 16)
			return uint16(id), err
		}
	}

	ID = sonyflake.NewSonyflake(sonyflake.Settings{
		StartTime:      time.UnixMilli(config.AppConfig.StartTime),
		MachineID:      MachineID,
		CheckMachineID: nil,
	})
}

func NextID() uint64 {
	id, err := ID.NextID()
	if err != nil {
		log.Error().Err(err).Msg("Generating ID")
	}
	return id
}
