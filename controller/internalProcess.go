package controller

import (
	"errors"

	"github.com/chilts/sid"

	"go.uber.org/zap"
)

type InternalProcess struct {
	Sugar *zap.SugaredLogger
	ID    string
}

type InternalProcessConfig struct {
	Sugar *zap.SugaredLogger
}

func (ip *InternalProcess) NewProcess(configRaw interface{}) error {
	config, ok := configRaw.(InternalProcessConfig)
	if !ok {
		return errors.New("config type mismatch")
	}

	ip.Sugar = config.Sugar
	ip.ID = sid.IdBase64()

	ip.Sugar.Infow("Succesfully created new internal process", "id", ip.ID)
	return nil
}

func (ip *InternalProcess) Stop() error {
	return nil
}

func (ip *InternalProcess) Halt() error {
	return nil
}

func (ip *InternalProcess) GetState() (StateSchema, interface{}, error) {
	return Internal, nil, nil
}
