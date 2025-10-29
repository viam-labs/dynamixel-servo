# Module dynamixel 

This is a [Viam module](https://docs.viam.com/how-tos/create-module/) for controlling [Dynamixel](https://www.robotis.us/) servo motors, like the XL320, XL430, XM430, and XM540 models.

## Model viam-labs:dynamixel:servo

This model provides the Dynamixel servo motor as a [Servo component](https://docs.viam.com/dev/reference/apis/components/servo/) in Viam.

### Configuration
The following attribute template can be used to configure this model:

```json
{
  "port": <string>,
  "id": <int>
}
```

#### Attributes

The following attributes are available for this model:

| Name          | Type   | Inclusion | Description                |
|---------------|--------|-----------|----------------------------|
| `port` | string  | Required  | The path to the serial port device for communicating with the servo |
| `id` | int | Required  | The ID of the servo, might be `0` out of the box |
| `baudrate` | int | Optional  | The baud rate for serial communication. Default is `1000000`. |

#### Example Configuration

```json
{
  "port": "/dev/tty.usbserial-FT4TFT52",
  "id": 0
}
```

### DoCommand

The following commands are availble in this servo module.

#### `set_torque`

Enable / disable the torque of the servo. It is enabled on start of the module. Disabling will allow it to be moved manually.

```json
{
  "command": "set_torque",
  "enable": <bool>
}
```

#### `ping`

Check that the servo is responding to controls over the serial interface.

```json
{
  "command": "ping"
}
```
