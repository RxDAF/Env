package server

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"os"
	"path"

	rmaster "github.com/RxDAF/Master/dto"
	"github.com/c4milo/unpackit"
)

func (s *Server) master() *rmaster.RMaster {
	var err error
	if s.mst == nil {
		s.mst, err = rmaster.NewRMaster(s.cfg.MasterServer)
		if err != nil {
			log.Println(err)
		}
	}
	return s.mst
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
func (s *Server) downloadService(serviceName string, downloadBufPath string) error {
	servicePath := path.Join(s.cfg.ServicePath, serviceName)
	data, err := s.mst.DownloadService(serviceName)
	if err != nil {
		return err
	}
	f, err := os.Create(downloadBufPath)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := f.Write(data); err != nil {
		return err
	}
	_, err = f.Seek(0, os.SEEK_SET)
	if err != nil {
		return err
	}
	// 清理之前解压的文件
	os.RemoveAll(servicePath)
	// 开始解压
	if err := os.MkdirAll(servicePath, os.ModePerm); err != nil {
		return err
	}
	// 默认都是tar.gz
	_, err = unpackit.Unpack(f, servicePath)
	return err
}
func (s *Server) checkIfNeedUpdate(serviceName string, downloadBufPath string) (bool, error) {
	// 对比md5
	f, err := os.Open(downloadBufPath)
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
	realMD5, err := s.mst.ServiceFileMD5(serviceName)
	if err != nil {
		return false, err
	}
	return realMD5 == md5str, nil
}
func (s *Server) checkLocalFile(serviceName string) error {
	bufPath := path.Join(s.cfg.DownloadBufPath, serviceName)
	ok, err := PathExists(bufPath)
	if err != nil {
		return err
	}
	if ok { // 本地存在则检查是否需要更新
		needUpdate, err := s.checkIfNeedUpdate(serviceName, bufPath)
		if err != nil {
			return err
		}
		if !needUpdate {
			return nil // 那就不需要修改
		}
		if needUpdate { // 需要更新就移除文件
			if err := os.Remove(bufPath); err != nil {
				return err
			}
		}
	}
	err = s.downloadService(serviceName, bufPath)
	return err
}
