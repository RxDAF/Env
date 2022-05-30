package server

import (
	"os/exec"

	"github.com/RxDAF/Env/cfg"
	rmaster "github.com/RxDAF/Master/dto"
)

type Server struct {
	cfg          *cfg.Configure
	mst          *rmaster.RMaster
	service      map[string]*ServiceProgress
	connToMaster rmaster.RMaster_StatusUpdateClient
}
type ServiceProgress struct {
	cmd exec.Cmd
}

func NewServer(cfg *cfg.Configure) *Server {
	return &Server{
		cfg:     cfg,
		service: map[string]*ServiceProgress{},
	}
}
func (s *Server) Run() error {
	// 与Master建立持久连接
	s.setupMasterConnect()
	// 先检测本地是否有下载的代码文件
	if err := s.checkLocalFiles(); err != nil {
		return err
	}
	// 开始依次运行
	if err := s.runServices(); err != nil {
		return err
	}
	return nil
}
