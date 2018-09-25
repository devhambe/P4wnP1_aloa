package service

import (
	"github.com/mame82/P4wnP1_go/service/datastore"
)


type Service struct {
	SubSysState          interface{}
	SubSysLogging        interface{}
	SubSysNetwork *NetworkManager
	SubSysDataStore      *datastore.Store
	SubSysEvent          *EventManager
	SubSysUSB            *UsbGadgetManager
	SubSysLed            *LedService
	SubSysWifi           *WiFiService
	SubSysBluetooth      *BtService
	SubSysRPC            *server
	SubSysTriggerActions interface{}
}

func NewService() (svc *Service, err error) {
	svc = &Service{}
	svc.SubSysLed = NewLedService()
	svc.SubSysNetwork, err = NewNetworkManager()
	if err != nil { return nil,err}
	svc.SubSysUSB,err = NewUSBGadgetManager(svc)
	if err != nil { return nil,err}
	svc.SubSysWifi = NewWifiService(svc)

	svc.SubSysRPC = NewRpcServerService(svc)
	return
}

func (s *Service) Start() {
	s.SubSysLed.Start()
	s.SubSysRPC.StartRpcServerAndWeb("0.0.0.0", "50051", "8000", "/usr/local/P4wnP1/www") //start gRPC service
}

func (s *Service) Stop() {
	s.SubSysLed.Stop()
}