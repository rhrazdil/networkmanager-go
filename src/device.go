package networkmanager

import (
	"github.com/godbus/dbus"
)

const (
	GetDevicesCall           = "org.freedesktop.NetworkManager.GetDevices"
	GetDeviceByIpIfaceCall   = "org.freedesktop.NetworkManager.GetDeviceByIpIface"
	DeviceInterfaceProperty  = "org.freedesktop.NetworkManager.Device.Interface"
	DeviceDeviceTypeProperty = "org.freedesktop.NetworkManager.Device.DeviceType"
)

type DeviceType string

const (
	DeviceTypeUnknown     DeviceType = "unknown"
	DeviceTypeGeneric     DeviceType = "generic"
	DeviceTypeEthernet    DeviceType = "ethernet"
	DeviceTypeWifi        DeviceType = "wifi"
	DeviceTypeUnused1     DeviceType = "unused1"
	DeviceTypeUnused2     DeviceType = "unused2"
	DeviceTypeBluetooth   DeviceType = "bluetooth"
	DeviceTypeOlcpMesh    DeviceType = "olcp-mesh"
	DeviceTypeWimax       DeviceType = "wimax"
	DeviceTypeModem       DeviceType = "modem"
	DeviceTypeInifiniband DeviceType = "inifiniband"
	DeviceTypeBond        DeviceType = "bond"
	DeviceTypeVlan        DeviceType = "vlan"
	DeviceTypeAdsl        DeviceType = "adsl"
	DeviceTypeBridge      DeviceType = "bridge"
	DeviceTypeTeam        DeviceType = "team"
	DeviceTypeTun         DeviceType = "tun"
	DeviceTypeIpTunnel    DeviceType = "ip-tunnel"
	DeviceTypeMacvlan     DeviceType = "macvlan"
	DeviceTypeVxlan       DeviceType = "vxlan"
	DeviceTypeVeth        DeviceType = "veth"
	DeviceTypeMacsec      DeviceType = "macsec"
	DeviceTypeDummy       DeviceType = "dummy"
)

var deviceTypeByNmDeviceType = map[uint32]DeviceType{
	0:  DeviceTypeUnknown,
	14: DeviceTypeGeneric,
	1:  DeviceTypeEthernet,
	2:  DeviceTypeWifi,
	3:  DeviceTypeUnused1,
	4:  DeviceTypeUnused2,
	5:  DeviceTypeBluetooth,
	6:  DeviceTypeOlcpMesh,
	7:  DeviceTypeWimax,
	8:  DeviceTypeModem,
	9:  DeviceTypeInifiniband,
	10: DeviceTypeBond,
	11: DeviceTypeVlan,
	12: DeviceTypeAdsl,
	13: DeviceTypeBridge,
	15: DeviceTypeTeam,
	16: DeviceTypeTun,
	17: DeviceTypeIpTunnel,
	18: DeviceTypeMacvlan,
	19: DeviceTypeVxlan,
	20: DeviceTypeVeth,
	21: DeviceTypeMacsec,
	22: DeviceTypeDummy,
}

type Device struct {
	Interface    string
	Type         DeviceType
	deviceObject dbus.BusObject
}

func (client *Client) newDeviceFromPath(devicePath dbus.ObjectPath) *Device {
	deviceObject := client.conn.Object(InterfacePath, devicePath)

	device := new(Device)
	device.deviceObject = deviceObject

	interfacePropertyVariant, _ := deviceObject.GetProperty(DeviceInterfaceProperty)
	device.Interface = interfacePropertyVariant.Value().(string)

	deviceTypePropertyVariant, _ := deviceObject.GetProperty(DeviceDeviceTypeProperty)
	nmDeviceType := deviceTypePropertyVariant.Value().(uint32)
	device.Type = deviceTypeByNmDeviceType[nmDeviceType]

	return device
}

func (client *Client) GetDevices() []*Device {
	call := client.conn.Object(InterfacePath, ObjectPath).Call(GetDevicesCall, 0)
	check(call.Err)
	var devicePaths []dbus.ObjectPath
	call.Store(&devicePaths)

	devices := make([]*Device, 0, len(devicePaths))
	for _, devicePath := range devicePaths {
		device := client.newDeviceFromPath(devicePath)
		devices = append(devices, device)
	}
	return devices
}

func (client *Client) GetDeviceByIpIface(ifname string) *Device {
	call := client.conn.Object(InterfacePath, ObjectPath).Call(GetDeviceByIpIfaceCall, 0, ifname)
	check(call.Err)
	var devicePath dbus.ObjectPath
	call.Store(&devicePath)
	return client.newDeviceFromPath(devicePath)
}