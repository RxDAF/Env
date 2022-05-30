package server

import (
	"log"
	"os"
	"os/exec"
	"path"
)

func (s *Server) runServices() error {
	for _, service := range s.cfg.Roles {
		if err := s.runService(service); err != nil {
			s.shutdownAllServices()
			return err
		}
	}
	return nil
}
func (s *Server) shutdownAllServices() error {
	for service := range s.service {
		if err := s.shutdownService(service); err != nil {
			log.Println(err)
		}
	}
	s.service = nil // 释放全部内存
	return nil
}
func (s *Server) shutdownService(serviceName string) error {
	service := s.service[serviceName]
	if err := service.cmd.Process.Signal(os.Interrupt); err != nil {
		log.Println("service:", service, " ", err)
	}
	// 开始释放
	if err := service.cmd.Process.Release(); err != nil {
		return err
	}
	s.service[serviceName] = nil
	// 发送状态更新
	s.updateMasterServiceStatus(serviceName, false, "shutdown")
	return nil
}
func (s *Server) runService(serviceName string) error {
	// 开始获取运行脚本
	setupScript := path.Join(s.cfg.ServicePath, serviceName, "setup.sh")
	s.service[serviceName].cmd = exec.Cmd{
		Path: setupScript,
	}
	// 开始启动
	err := s.service[serviceName].cmd.Start()
	if err != nil {
		return err
	}
	s.updateMasterServiceStatus(serviceName, true, "run")
	return nil
}
