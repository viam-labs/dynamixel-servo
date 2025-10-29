package dynamixel

import (
	"context"
	"fmt"

	"go.uber.org/multierr"

	servo "go.viam.com/rdk/components/servo"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/resource"

	"github.com/jacobsa/go-serial/serial"
	"go.viam.com/dynamixel/network"
	dynServo "go.viam.com/dynamixel/servo"
	"go.viam.com/dynamixel/servo/s_model"
)

var (
	Servo = resource.NewModel("viam-labs", "dynamixel", "servo")
)

func init() {
	resource.RegisterComponent(servo.API, Servo,
		resource.Registration[servo.Servo, *Config]{
			Constructor: newDynamixelServo,
		},
	)
}

type Config struct {
	Port     string `json:"port,omitempty"`
	BaudRate int    `json:"baudrate,omitempty"`
	Id       int    `json:"id,omitempty"`
}

// Validate ensures all parts of the config are valid and important fields exist.
// Returns implicit required (first return) and optional (second return) dependencies based on the config.
// The path is the JSON path in your robot's config (not the `Config` struct) to the
// resource being validated; e.g. "components.0".
func (cfg *Config) Validate(path string) ([]string, []string, error) {
	if cfg.Port == "" {
		return nil, nil, fmt.Errorf("must specify port for serial communication")
	}

	return nil, nil, nil
}

type dynamixelServo struct {
	resource.AlwaysRebuild

	name resource.Name

	logger logging.Logger
	cfg    *Config

	cancelCtx  context.Context
	cancelFunc func()

	servo *dynServo.Servo
}

func newDynamixelServo(ctx context.Context, deps resource.Dependencies, rawConf resource.Config, logger logging.Logger) (servo.Servo, error) {
	conf, err := resource.NativeConfig[*Config](rawConf)
	if err != nil {
		return nil, err
	}

	return NewServo(ctx, deps, rawConf.ResourceName(), conf, logger)

}

func NewServo(ctx context.Context, deps resource.Dependencies, name resource.Name, conf *Config, logger logging.Logger) (servo.Servo, error) {

	if conf.BaudRate == 0 {
		conf.BaudRate = 1000000
	}

	serialOptions := serial.OpenOptions{
		PortName:              conf.Port,
		BaudRate:              uint(conf.BaudRate),
		DataBits:              8,
		StopBits:              1,
		MinimumReadSize:       0,
		InterCharacterTimeout: 100,
	}
	serial, err := serial.Open(serialOptions)
	if err != nil {
		return nil, fmt.Errorf("error opening serial port: %v\n", err)
	}

	net := network.New(serial)
	svo, err := s_model.New(net, conf.Id)

	err = svo.Ping()
	if err != nil {
		return nil, fmt.Errorf("Unable to ping servo: %v", err)
	}

	err = svo.SetTorqueEnable(true)
	if err != nil {
		return nil, fmt.Errorf("Unable to enable torque on servo: %v", err)
	}

	cancelCtx, cancelFunc := context.WithCancel(context.Background())

	s := &dynamixelServo{
		name:       name,
		logger:     logger,
		cfg:        conf,
		cancelCtx:  cancelCtx,
		cancelFunc: cancelFunc,
		servo:      svo,
	}
	return s, nil
}

func (s *dynamixelServo) Name() resource.Name {
	return s.name
}

// Move moves the servo to the given angle (0-180 degrees).
// This will block until done or a new operation cancels this one.
func (s *dynamixelServo) Move(ctx context.Context, angleDeg uint32, extra map[string]interface{}) error {
	s.logger.Debug("Moving servo to: ", angleDeg)
	return s.servo.MoveTo(float64(angleDeg))
}

// Position returns the current set angle (degrees) of the servo.
func (s *dynamixelServo) Position(ctx context.Context, extra map[string]interface{}) (uint32, error) {
	position, err := s.servo.Angle()
	if err != nil {
		return 0, err
	}
	return uint32(position), nil
}

func (s *dynamixelServo) Stop(ctx context.Context, extra map[string]interface{}) error {
	s.logger.Debug("Stopping servo")
	return multierr.Combine(
		s.servo.SetTorqueEnable(false),
		s.servo.SetTorqueEnable(true),
	)
}

func (s *dynamixelServo) DoCommand(ctx context.Context, cmd map[string]interface{}) (map[string]interface{}, error) {
	switch cmd["command"] {
	case "set_torque":
		enable, ok := cmd["enable"].(bool)
		if !ok {
			return nil, fmt.Errorf("set_torque command requires 'enable' boolean parameter")
		}
		err := s.servo.SetTorqueEnable(enable)
		return map[string]interface{}{"success": err == nil}, err
	case "ping":
		err := s.servo.Ping()
		return map[string]interface{}{"success": err == nil}, err
	default:
		return nil, fmt.Errorf("unknown command: %v", cmd)
	}
}

func (s *dynamixelServo) IsMoving(ctx context.Context) (bool, error) {
	moving, err := s.servo.Moving()
	if err != nil {
		return false, err
	}
	return moving == 1, nil
}

func (s *dynamixelServo) Close(ctx context.Context) error {
	s.cancelFunc()
	s.Stop(ctx, nil)
	return nil
}
