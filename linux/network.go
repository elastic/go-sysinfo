package linux

import (
	"net"
)

func Network() (ips, macs []string, err error) {
	ifcs, err := net.Interfaces()
	if err != nil {
		return nil, nil, err
	}

	ips = make([]string, 0, len(ifcs))
	macs = make([]string, 0, len(ifcs))
	for _, ifc := range ifcs {
		addrs, err := ifc.Addrs()
		if err != nil {
			return nil, nil, err
		}
		for _, addr := range addrs {
			ips = append(ips, addr.String())
		}

		mac := ifc.HardwareAddr.String()
		if mac != "" {
			macs = append(macs, mac)
		}
	}

	return ips, macs, nil
}
