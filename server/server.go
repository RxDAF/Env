package server

import (
	"github.com/RxDAF/Env/cfg"
	rmaster "github.com/RxDAF/Master/dto"
)

type Server struct {
	cfg *cfg.Configure
	mst *rmaster.RMaster
}

func NewServer(cfg *cfg.Configure) *Server {
	return &Server{
		cfg: cfg,
	}
}
func (s *Server) Run() error {
	// 先检测本地是否有下载的代码文件
	if err := s.checkLocalFiles(); err != nil {
		return err
	}
	return nil
}
