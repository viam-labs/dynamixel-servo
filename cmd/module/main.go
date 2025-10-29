package main

import (
	"dynamixel"
	"go.viam.com/rdk/module"
	"go.viam.com/rdk/resource"
	servo "go.viam.com/rdk/components/servo"
)

func main() {
	// ModularMain can take multiple APIModel arguments, if your module implements multiple models.
	module.ModularMain(resource.APIModel{ servo.API, dynamixel.Servo})
}
