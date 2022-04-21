package session

import (
	"errors"
	"strings"
	"sync"
	"time"

	. "github.com/chefsgo/base"
	"github.com/chefsgo/chef"
)

var (
	errNotFound = errors.New("NotFound")
)

type (
	defaultSessionDriver  struct{}
	defaultSessionConnect struct {
		name     string
		config   chef.SessionConfig
		setting  defaultSessionSetting
		sessions sync.Map
	}
	defaultSessionSetting struct {
	}
	defaultSessionValue struct {
		Value  Map
		Expiry time.Time
	}
)

//连接
func (driver *defaultSessionDriver) Connect(name string, config chef.SessionConfig) (chef.SessionConnect, error) {
	setting := defaultSessionSetting{}

	return &defaultSessionConnect{
		name: name, config: config, setting: setting,
		sessions: sync.Map{},
	}, nil
}

//打开连接
func (connect *defaultSessionConnect) Open() error {
	return nil
}

//关闭连接
func (connect *defaultSessionConnect) Close() error {
	return nil
}

//查询会话，
func (connect *defaultSessionConnect) Read(id string) (Map, error) {
	if value, ok := connect.sessions.Load(id); ok {
		if vv, ok := value.(defaultSessionValue); ok {
			if vv.Expiry.Unix() > time.Now().Unix() {
				return vv.Value, nil
			} else {
				//过期了就删除
				connect.Delete(id)
			}
		}
	}
	return nil, errNotFound
}

//更新会话
func (connect *defaultSessionConnect) Write(id string, val Map, expiry time.Duration) error {
	now := time.Now()

	if expiry <= 0 {
		expiry = connect.config.Expiry
	}

	value := defaultSessionValue{
		Value: val, Expiry: now.Add(expiry),
	}

	connect.sessions.Store(id, value)

	return nil
}

//删除会话
func (connect *defaultSessionConnect) Delete(id string) error {
	connect.sessions.Delete(id)
	return nil
}

//清空会话
func (connect *defaultSessionConnect) Clear(prefix string) error {
	connect.sessions.Range(func(k, _ Any) bool {
		if strings.HasPrefix(k, prefix) {
			connect.sessions.Delete(k)
		}
		return true
	})
	return nil
}
