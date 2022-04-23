package session_default

import (
	"errors"
	"strings"
	"sync"
	"time"

	. "github.com/chefsgo/base"
	"github.com/chefsgo/session"
)

var (
	errNotFound = errors.New("NotFound")
)

type (
	defaultDriver  struct{}
	defaultConnect struct {
		name     string
		config   session.Config
		setting  defaultSetting
		sessions sync.Map
	}
	defaultSetting struct {
	}
	defaultValue struct {
		Value  Map
		Expiry time.Time
	}
)

//连接
func (driver *defaultDriver) Connect(name string, config session.Config) (session.Connect, error) {
	setting := defaultSetting{}

	return &defaultConnect{
		name: name, config: config, setting: setting,
		sessions: sync.Map{},
	}, nil
}

//打开连接
func (connect *defaultConnect) Open() error {
	return nil
}

//关闭连接
func (connect *defaultConnect) Close() error {
	return nil
}

//查询会话，
func (connect *defaultConnect) Read(id string) (Map, error) {
	if value, ok := connect.sessions.Load(id); ok {
		if vv, ok := value.(defaultValue); ok {
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
func (connect *defaultConnect) Write(id string, val Map, expiry time.Duration) error {
	now := time.Now()

	if expiry <= 0 {
		expiry = connect.config.Expiry
	}

	value := defaultValue{
		Value: val, Expiry: now.Add(expiry),
	}

	connect.sessions.Store(id, value)

	return nil
}

//删除会话
func (connect *defaultConnect) Delete(id string) error {
	connect.sessions.Delete(id)
	return nil
}

//清空会话
func (connect *defaultConnect) Clear(prefix string) error {
	connect.sessions.Range(func(k, _ Any) bool {
		if strings.HasPrefix(k.(string), prefix) {
			connect.sessions.Delete(k)
		}
		return true
	})
	return nil
}
