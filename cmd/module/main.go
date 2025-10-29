package main

import (
	"dynamixel"
	servo "go.viam.com/rdk/components/servo"
	"go.viam.com/rdk/module"
	"go.viam.com/rdk/resource"
)

func main() {
	// ModularMain can take multiple APIModel arguments, if your module implements multiple models.
	module.ModularMain(resource.APIModel{servo.API, dynamixel.Servo})
}
