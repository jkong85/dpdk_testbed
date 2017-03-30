package ovsBridge

import (
	"errors"
	"service/log"
	"snc/service/template"
	"sync"
	"utils/runcmd"
)

type ovsBrHandler struct {
	sync.Mutex

	bridgeData
}

type bridgeData struct {
	bridgeName     string
	bridgeIP       string
	bridgeMac      string
	devDeployWith  string
	intfDeployWith string

	nicName      string
	portName     string
	peerPortName string
}

var dpdk_enable = 0

func newOvsBrHandler() *ovsBrHandler {

	brData := bridgeData{}
	return &ovsBrHandler{bridgeData: brData}
}

func (h *ovsBrHandler) parseData(compInfo []byte) {

	nodeStru := template.UnmarshalNode(compInfo)
	//	log.Debugln("nodeStru=", nodeStru)
	//	log.Debugln("nodeStru.Name=", nodeStru.Name)
	//	log.Debugln("nodeStru.Type=", nodeStru.Type)
	//	log.Debugln("nodeStru.PerContainer=", nodeStru.PerContainer)
	//	log.Debugln("nodeStru.CapRequire=", nodeStru.CapRequire)
	//	log.Debugln("nodeStru.LinkPoints=", nodeStru.LinkPoints)
	//	log.Debugln("nodeStru.IP=", nodeStru.IP)

	h.bridgeName = nodeStru.Name
	h.bridgeIP = nodeStru.IP

}

func getInstance(bridgeName string) *ovsBrHandler {

	if instanceTemp, found := instanceMap[bridgeName]; found {
		log.Debugln("getInstance-->instanceTemp= ", instanceTemp)
		return instanceTemp
	}
	return nil
}

func (h *ovsBrHandler) newDevice() {

	log.Debugln("ovsBHandler-->newDevice:")

	if dpdk_enable {
		log.Debugln(" New dpdk device: ")
		log.Debugln("bridgeName=", h.bridgeName)
		//ovs-vsctl add-br br0 -- set bridge br0 datapath_type=netdev
		_, err := runcmd.OpCommandDirect("ovs-vsctl add-br " + h.bridgeName + " -- set bridge " + h.bridgeName + " datapath_type=netdev ")
		if err != nil {
			log.Error("setBridge,   ", err.Error())
		}

		runcmd.OpCommandDirect("sleep 0.01")
		return
	}

	log.Debugln("bridgeName=", h.bridgeName)
	_, err := runcmd.OpCommandDirect("ovs-vsctl add-br " + h.bridgeName)
	if err != nil {
		log.Error("setBridge,   ", err.Error())
	}

	runcmd.OpCommandDirect("sleep 0.01")

	if h.bridgeIP != "" {

		log.Debugln("set bridge ip -->h.bridgeIP=", h.bridgeIP)
		_, err = runcmd.OpCommandDirect("ifconfig " + h.bridgeName + " " + h.bridgeIP)
		if err != nil {
			log.Error("setBridge, ifconfig  ", err.Error())
		}
	}

	_, err = runcmd.OpCommandDirect("ifconfig " + h.bridgeName + " up ")
	if err != nil {
		log.Error("setBridge, ifconfig  ", err.Error())
	}

}

func (h *ovsBrHandler) destroyDevice() {

	log.Debugln("ovsBHandler-->destroyDevice:")
	log.Debugln("bridgeName=", h.bridgeName)
	_, err := runcmd.OpCommandDirect("ovs-vsctl del-br " + h.bridgeName)
	if err != nil {
		log.Error("delBridge,   ", err.Error())
	}

	runcmd.OpCommandDirect("sleep 0.01")
}

func (h *ovsBrHandler) update() error {

	log.Debugln("ovsBHandler-->update:")

	return nil
}

func (h *ovsBrHandler) queryCommonInfo() (string, error) {

	log.Debugln("ovsBHandler-->queryCommonInfo:")

	log.Debugln("bridgeName=", h.bridgeName)
	return "", nil
}

func (h *ovsBrHandler) ConnectToIntf(bridgeName string, vlanID string) error {

	if h.intfDeployWith == "veth" {
		return VethConnectToIntf(bridgeName, h.nicName, vlanID)
	} else if h.intfDeployWith == "patch_port" {

		return PatchPortConnectToIntf(bridgeName, h.portName, h.peerPortName)
	}

	return nil
}

func VethConnectToIntf(bridgeName string, nicName string, vlanID string) error {

	log.Debugf("ovsBHandler-->VethConnectToIntf:--bridgeName=%s--nicName=%s", bridgeName, nicName)
	//ConnectToIntf:--bridgeName=overlay_br_tun--nicName=PairedPortToTun_link
	if vlanID == "" {
		_, err := runcmd.OpCommandDirect("ovs-vsctl add-port " + bridgeName + " " + nicName)
		if err != nil {
			return errors.New("ovs add-port fail, " + err.Error())
		}
		//runcmd.OpCommandDirect("sleep 0.1")
		runcmd.OpCommandDirect("sleep 0.01")

		if dpdk_enable {
			// then we add only one vhost-user port even there are many veth
			log.Debugf("DPDK:Add vhost-user-1 to the ovs")
			vhostName := "vhost-user-1"
			// ovs-vsctl add-port br0 vhost-user-1 -- set Interface vhost-user-1 type=dpdkvhostuser
			_, err := runcmd.OpCommandDirect("ovs-vsctl add-port " + bridgeName + " " + nicName + " -- set Interface " + nicName + " type=dpdkvhostuser ")
			if err != nil {
				return errors.New("ovs add-port fail, " + err.Error())
			} else {
				log.Debugf("Add the vhost-user-1 successfully for the first time!")
			}
			runcmd.OpCommandDirect("sleep 0.01")
		}

	} else {

		log.Debugln("vlanID = ", vlanID)
		_, err := runcmd.OpCommandDirect("ovs-vsctl add-port " + bridgeName + " " + nicName + " tag=" + vlanID)
		if err != nil {
			return errors.New("ovs add-port fail, " + err.Error())
		}
		//runcmd.OpCommandDirect("sleep 0.1")
		runcmd.OpCommandDirect("sleep 0.01")
	}

	return nil
}

func PatchPortConnectToIntf(bridgeName string, portName string, peerPortName string) error {

	log.Debugf("ovsBHandler-->PatchPortConnectToIntf:--bridgeName=%s--portName=%s peerPortName=%s", bridgeName, portName, peerPortName)

	_, err := runcmd.OpCommandDirect("ovs-vsctl set interface  " + portName + " options:peer=" + peerPortName)
	if err != nil {
		return errors.New("ovs add-port fail, " + err.Error())
	}
	runcmd.OpCommandDirect("sleep 0.01")

	return nil
}

func (h *ovsBrHandler) LeaveFromIntf(nicName string) error {
	return nil //weixu hack for overlay problem, debug later.
	log.Debugln("ovsBHandler-->LeaveFromIntf:")

	log.Debugln("bridgeName=", h.bridgeName)
	_, err := runcmd.OpCommandDirect("ovs-vsctl del-port " + h.bridgeName + " " + nicName)
	if err != nil {
		return errors.New("ovs del-port fail, " + err.Error())
	}
	runcmd.OpCommandDirect("sleep 0.01")

	return nil
}
