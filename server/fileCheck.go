package server

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"

	rmaster "github.com/RxDAF/Master"
)

func (s *Server) master() *rmaster.RMaster {

}
func (s *Server) checkLocalFiles() error {
	for _, serviceName := range s.cfg.Roles {
		if err := s.checkLocalFile(serviceName); err != nil {
			return err
		}
	}
	return nil
}
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
func (s *Server) downloadService(serviceName string, path string) error {

}
func (s *Server) checkIfNeedUpdate(serviceName string, path string) (bool, error) {
	// 对比md5
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	md5Handle := md5.New()         //创建 md5 句柄
	_, err = io.Copy(md5Handle, f) //将文件内容拷贝到 md5 句柄中
	if nil != err {
		return false, err
	}
	md := md5Handle.Sum(nil)        //计算 MD5 值，返回 []byte
	md5str := fmt.Sprintf("%x", md) //将 []byte 转为 string
	// 开始获取远程服务器的md5

}
func (s *Server) checkLocalFile(serviceName string) error {
	path := s.cfg.DownloadBufPath + string(os.PathSeparator) + serviceName
	ok, err := PathExists(path)
	if err != nil {
		return err
	}
	if !ok {
		if err := s.downloadService(serviceName, path); err != nil {
			return err
		}
	} else {
		needUpdate, err := s.checkIfNeedUpdate(serviceName, path)
		if err != nil {
			return err
		}
		if needUpdate {
			if err := os.Remove(path); err != nil {
				return err
			}
			s.downloadService(serviceName, path)
		}
	}
}
