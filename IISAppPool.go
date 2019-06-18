package goToolAgentWin

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

type iisAppPool struct {
	AppPoolPath  string
	StartTimeout int
	StartDelay   int
}

func NewIISAppPool(appPoolPath string, startTimeout int, startDelay int) *iisAppPool {
	return &iisAppPool{
		AppPoolPath:  appPoolPath,
		StartTimeout: startTimeout,
		StartDelay:   startDelay,
	}
}

//重启IIS应用程序池
func (ap *iisAppPool) Restart(name string) error {
	exist, err := ap.isExist(name)
	if err != nil {
		return err
	}
	if !exist {
		return errors.New("AppPool is not exist")
	}
	b, err := ap.IsRunning(name)
	if err != nil {
		return err
	}
	if b {
		err = ap.Stop(name)
		if err != nil {
			return err
		}
		time.Sleep(time.Second * 5)
	}
	outTime := time.Now().Add(time.Duration(ap.StartTimeout * 1000 * 1000 * 1000))
	for {
		err = ap.Start(name)
		if err != nil {
			return err
		}
		time.Sleep(time.Duration(ap.StartDelay * 1000 * 1000 * 1000))
		check, err := ap.IsRunning(name)
		if err != nil {
			return err
		}
		if check {
			return nil
		}
		if time.Now().After(outTime) {
			errMsg := fmt.Sprintf("Restart AppPool  %s timeout", name)
			return errors.New(errMsg)
		}
	}
}

//检测IIS应用程序池是否存在
func (ap *iisAppPool) isExist(name string) (bool, error) {
	cmd := exec.Command(ap.AppPoolPath, "list", "appPools")
	out, err := cmd.CombinedOutput()
	if err != nil {
		errMsg := fmt.Sprintf("检测IIS应用程序池是否存在时遇到错误：%s", err.Error())
		return false, errors.New(errMsg)
	}
	if strings.Index(strings.ToLower(string(out)), strings.ToLower(name)) > 0 {
		return true, nil
	}
	return false, nil
}

//检测IIS应用程序池是否在运行
func (ap *iisAppPool) IsRunning(name string) (bool, error) {
	exist, err := ap.isExist(name)
	if err != nil {
		return false, err
	}
	if !exist {
		return false, errors.New("AppPool is not exist")
	}
	cmd := exec.Command(ap.AppPoolPath, "list", "appPools", "/state:started")
	out, err := cmd.CombinedOutput()
	if err != nil {
		errMsg := fmt.Sprintf("检测IIS应用程序池是否在运行时遇到错误：%s", err.Error())
		return false, errors.New(errMsg)
	}
	outStr := string(out)
	outStr = strings.ToLower(outStr)
	if strings.Index(outStr, strings.ToLower(name)) > 0 {
		return true, nil
	}
	return false, nil
}

//启动IIS应用程序池
func (ap *iisAppPool) Start(name string) error {
	exist, err := ap.isExist(name)
	if err != nil {
		return err
	}
	if !exist {
		return errors.New("AppPool is not exist")
	}
	running, err := ap.IsRunning(name)
	if err != nil {
		return err
	}
	if running {
		return nil
	}
	cmd := exec.Command(ap.AppPoolPath, "start", "appPool", fmt.Sprintf("/appPool.name:%s", name))
	_, err = cmd.CombinedOutput()
	if err != nil {
		errMsg := fmt.Sprintf("启动应用程序池时遇到错误：%s", err.Error())
		return errors.New(errMsg)
	}
	return nil
}

//停止IIS应用程序池
func (ap *iisAppPool) Stop(name string) error {
	exist, err := ap.isExist(name)
	if err != nil {
		return err
	}
	if !exist {
		return errors.New("AppPool is not exist")
	}
	running, err := ap.IsRunning(name)
	if err != nil {
		return err
	}
	if !running {
		return nil
	}
	cmd := exec.Command(ap.AppPoolPath, "stop", "appPool", fmt.Sprintf("/appPool.name:%s", name))
	_, err = cmd.CombinedOutput()
	if err != nil {
		errMsg := fmt.Sprintf("停止应用程序池时遇到错误：%s", err.Error())
		return errors.New(errMsg)
	}
	return nil
}
