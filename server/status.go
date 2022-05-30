package server

import rmaster "github.com/RxDAF/Master/dto"

func (s *Server) updateServiceStatus(serviceName string, newStatus bool) error {
	if newStatus {
		if s.service[serviceName] != nil {
			return nil
		}
		// 否则就启动
		return s.runService(serviceName)
	}
	if s.service[serviceName] == nil {
		return nil
	}
	// 否则开始终止
	return s.shutdownService(serviceName)
}
func (s *Server) updateMasterServiceStatus(serviceName string, newStatus bool, extraInfo string) error {
	return s.connToMaster.Send(&rmaster.StatusUpdateInfo{
		StatusUpdate: &rmaster.StatusUpdateInfo_Service{
			Service: &rmaster.ServiceStatusChange{
				ServiceName: serviceName,
				NewStatus:   newStatus,
				ExtraInfo:   extraInfo,
			},
		},
	})
}
