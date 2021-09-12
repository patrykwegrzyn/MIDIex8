# MIDIex8
Hot-plug Firmware update tool for MIDIex8 `0x1000`

Go reimplementation of [@ykcirtsyb midex-pid-changer](https://github.com/ykcirtsyb/midex-pid-changer)

## Require
  - libusb
  - [@sgorpi midex8 Linux driver](https://github.com/sgorpi/midex8)

For proper functionality it is necessary to modify original code of [midex.c](https://github.com/sgorpi/midex8/blob/master/src/kernel/sound/usb/midex/midex.c#L72) as follows:

```c
#define SB_MIDEX_VID 0x0a4e
#define SB_MIDEX8_PID 0x1001

/*******************************************************************
 * Type definitions
 *******************************************************************/

static struct usb_device_id id_table[] = {
	{ USB_DEVICE(SB_MIDEX_VID, SB_MIDEX8_PID) },
	{ },
};
```


## Build

`make build`

### Install

`sudo make install`


