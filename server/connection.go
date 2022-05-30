package server

import (
	"log"

	rmaster "github.com/RxDAF/Master/dto"
)

func (s *Server) setupMasterConnect() error {
	var err error
	s.connToMaster, err = s.mst.SetupConnect()
	if err != nil {
		return err
	}
	if err := s.connToMaster.Send(&rmaster.StatusUpdateInfo{
		StatusUpdate: &rmaster.StatusUpdateInfo_Certification{
			Certification: &rmaster.StatusUpdateInfoCertification{
				Address: s.cfg.Address,
			},
		},
	}); err != nil {
		return err
	}
	// 否则就开启连接
	go func() {
		for {
			data, err := s.connToMaster.Recv()
			if err != nil {
				log.Println(err)
				return
			}
			service := data.GetService()
			if service != nil {
				s.updateServiceStatus(service.ServiceName, service.NewStatus)
			}
		}
	}()
	return nil
}
