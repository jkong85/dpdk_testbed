package ovsBridge

import (
	"errors"
	"service/log"
	"snc/driver/compAPI"
	"snc/service/template"
	"strings"
)

var instanceMap = make(map[string]*ovsBrHandler)

type ovsBridgeAPI struct {
	//	Create(compInfo []byte) (priData string, err error)
	//	Destroy(priData string, compInfo []byte) error
	//	Update(priData string, compInfo []byte) error
	//	QueryCommonInfo(enum string, priData string, compInfo []byte) (string, error)
	//	ConnectToIntf(priData string, DevCompInfo []byte, intfCompInfo []byte) error
	//	LeaveFromIntf(priData string, DevCompInfo []byte, intfCompInfo []byte) error
}

func Init() {

	log.Debugln("ovsBridge : init")
	d := ovsBridgeAPI{}
	err := compAPI.RegisterDriver("l2_dev", "ovs", &d)
	if err != nil {
		panic(err)
	}
}

func (lb *ovsBridgeAPI) Create(compInfo []byte) (bridgeName string, err error) {

	log.Debugln("@@ovsBridge : Create")

	instanceTemp := newOvsBrHandler()

	instanceTemp.parseData(compInfo)

	log.Debugln("bridgeName= ", instanceTemp.bridgeName)
	log.Debugln("instanceTemp= ", instanceTemp)
	instanceMap[instanceTemp.bridgeName] = instanceTemp
	log.Debugln("instanceMap=", instanceMap)
	for k, v := range instanceMap {
		log.Debugln("k=", k)
		log.Debugln("v=", v)
	}

	instanceTemp.newDevice()
	return bridgeName, nil
}

func (lb *ovsBridgeAPI) Destroy(bridgeName string, compInfo []byte) error {

	log.Debugln("@@ovsBridge : Destroy")
	nodeData := template.UnmarshalNode(compInfo)
	log.Debugln("nodeData=", nodeData)

	instanceTemp := getInstance(bridgeName)
	if instanceTemp == nil {
		return errors.New("ovsBHandler not found:" + bridgeName)
	}
	instanceTemp.destroyDevice()
	delete(instanceMap, instanceTemp.bridgeName)
	return nil
}

func (lb *ovsBridgeAPI) Update(bridgeName string, compInfo []byte) error {

	log.Debugln("@@ovsBridge : Update")
	nodeData := template.UnmarshalNode(compInfo)
	log.Debugln("nodeData=", nodeData)

	instanceTemp := getInstance(bridgeName)
	if instanceTemp == nil {
		return errors.New("ovsBHandler not found:" + bridgeName)
	}
	return instanceTemp.update()
}

func (lb *ovsBridgeAPI) QueryCommonInfo(bridgeName string, compInfo []byte) (string, error) {

	log.Debugln("@@ovsBridge : QueryCommonInfo")
	nodeData := template.UnmarshalNode(compInfo)
	log.Debugln("nodeData=", nodeData)

	instanceTemp := getInstance(bridgeName)
	if instanceTemp == nil {
		return "", errors.New("ovsBHandler not found:" + bridgeName)
	}
	return instanceTemp.queryCommonInfo()
}

func (lb *ovsBridgeAPI) ConnectToIntf(bridgeName string, devCompInfo []byte, intfCompInfo []byte) error {

	log.Debugln("@@ovsBridge : ConnectToIntf")
	log.Debugln("bridgeName=", bridgeName)

	instanceTemp := getInstance(bridgeName)
	//devData := template.UnmarshalNode(devCompInfo)
	//	log.Debugln("devData=", devData)
	//	log.Debugln("devData.deploywith=", devData.DeployWith)
	if instanceTemp == nil {
		return errors.New("ovsBHandler not found:" + bridgeName)
	}

	//vlanID := "100"
	var vlanID string
	intfData := template.UnmarshalNode(intfCompInfo)

	//	log.Debugln("intfData.deploywith=", intfData.DeployWith)
	instanceTemp.intfDeployWith = intfData.DeployWith

	if intfData.DeployWith == "veth" {

		externSetting, err := template.ParseVlanExternSetting(intfData.ExternSetting)
		if err != nil {
			log.Errorln(" ParseVlanExternSetting err :", err)
			return err
		}

		vlanID = externSetting.VlanID
		log.Debugln("externSetting.VlanID=", externSetting.VlanID)

		for _, value := range intfData.LinkPoints {

			//log.Debugln("value=", value)
			if value.LinkTo != "Net_Namespace" {
				instanceTemp.nicName = value.Name
				//log.Debugln("nicName=", instanceTemp.nicName)
			}

			if instanceTemp.nicName == "" {
				return errors.New("no find nicName")
			}
		}
	} else if intfData.DeployWith == "patch_port" {

		for _, value := range intfData.LinkPoints {

			//log.Debugln("value=", value)
			lt := strings.Split(value.LinkTo, ".")
			if len(lt) > 2 {
				devName := lt[1]
				if bridgeName == devName {
					instanceTemp.portName = value.Name
					//log.Debugln("portName=value= ", instanceTemp.portName)
				} else {

					instanceTemp.peerPortName = value.Name
					//log.Debugln("peerPortName=value= ", instanceTemp.peerPortName)
				}
			}
		}
	}

	instanceTemp.ConnectToIntf(bridgeName, vlanID)
	return nil
}

func (lb *ovsBridgeAPI) LeaveFromIntf(bridgeName string, devCompInfo []byte, intfCompInfo []byte) error {

	log.Debugln("@@ovsBridge : LeaveFromIntf")

	//devData := template.UnmarshalNode(devCompInfo)
	//log.Debugln("devData=", devData)

	intfData := template.UnmarshalNode(intfCompInfo)
	//log.Debugln("intfData=", intfData)

	//nicName := intfData.LinkPoints[1].Name
	//log.Debugln("nicName=", nicName)
	instanceTemp := getInstance(bridgeName)
	if instanceTemp == nil {
		return errors.New("ovsBHandler not found:" + bridgeName)
	}

	if intfData.DeployWith == "veth" {

		for _, value := range intfData.LinkPoints {

			//log.Debugln("value=", value)
			if value.LinkTo != "Net_Namespace" {
				instanceTemp.nicName = value.Name
				//log.Debugln("nicName=", instanceTemp.nicName)
			}

			if instanceTemp.nicName == "" {
				return errors.New("no find nicName")
			}
		}

		instanceTemp.LeaveFromIntf(instanceTemp.nicName)
	}
	return nil
}
