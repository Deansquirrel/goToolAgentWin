package goToolAgentWin

import (
	"errors"
	"os/exec"
	"strings"
)

type winService struct {
}

func NewWinService() *winService {
	return &winService{}
}

//重启Windows服务
func (ws *winService) Restart(name string) error {
	b, err := ws.IsRunning(name)
	if err != nil {
		return err
	}
	if b {
		err = ws.Stop(name)
		if err != nil {
			return err
		}
	}
	err = ws.Start(name)
	if err != nil {
		return err
	}
	return nil
}

//监测Windows服务是否在运行
func (ws *winService) IsRunning(name string) (bool, error) {
	cmd := exec.Command("sc", "query", name)
	out, err := cmd.CombinedOutput()
	if err != nil {
		if err.Error() == "exit status 1060" {
			return false, errors.New("指定的服务[" + name + "]未安装")
		}
		return false, err
	}
	if strings.Index(string(out), "STOPPED") > 0 {
		return false, nil
	}
	if strings.Index(string(out), "RUNNING") > 0 {
		return true, nil
	}
	return false, errors.New("状态检查失败")
}

//停止Windows服务
func (ws *winService) Stop(name string) error {
	cmd := exec.Command("net", "stop", name)
	_, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	return nil
}

//启动Windows服务
func (ws *winService) Start(name string) error {
	cmd := exec.Command("net", "start", name)
	_, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	return nil
}
