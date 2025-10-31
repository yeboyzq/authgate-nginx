/*
Copyright (c) 2025 authgate-nginx
authgate-nginx is licensed under Mulan PSL v2.
You can use this software according to the terms and conditions of the Mulan PSL v2.
You may obtain a copy of Mulan PSL v2 at:
        http://license.coscl.org.cn/MulanPSL2
THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
See the Mulan PSL v2 for more details.
*/

package modules

import (
	"crypto/tls"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/go-ldap/ldap/v3"
	"github.com/yeboyzq/authgate-nginx/app/modules/config"
	"github.com/yeboyzq/authgate-nginx/app/modules/log"
)

const DefaultTimeout time.Duration = 30 * time.Second

var Ldap *LdapInstance

type LdapInstance struct {
	url          string
	skipVerify   bool
	bindDN       string
	bindPassword string
	baseDN       string
	filter       string
}

type UserInfo struct {
	Username string
	Email    string
	DN       string
	Groups   []string
}

// LdapInit 实例化ldap
func LdapInit() {
	Ldap = &LdapInstance{
		url:          config.Cfg.GetString("base.ldap.url"),
		skipVerify:   config.Cfg.GetBool("base.ldap.skipVerify"),
		bindDN:       config.Cfg.GetString("base.ldap.bindDn"),
		bindPassword: config.Cfg.GetString("base.ldap.bindPassword"),
		baseDN:       config.Cfg.GetString("base.ldap.userBaseDn"),
		filter:       config.Cfg.GetString("base.ldap.filter"),
	}

	log.Info("LDAP组件初始化完成.")
}

// GetLdapConn 获取ldap连接
func (s *LdapInstance) GetLdapConn() (*ldap.Conn, error) {
	conn, err := ldap.DialURL(s.url)
	if err != nil {
		return nil, err
	}
	// defer conn.Close()
	// 设置连接超时
	conn.SetTimeout(DefaultTimeout)
	ldapHost, err := extractIPFromLDAPURL(s.url)
	if err == nil && ldapHost != "" {
		err = conn.StartTLS(&tls.Config{
			InsecureSkipVerify: s.skipVerify,
			ServerName:         ldapHost,
		})
		if err != nil {
			log.Warn(err.Error())
		}
	}

	return conn, err
}

// Authenticate 认证用户
func (s *LdapInstance) Authenticate(username, password string) (bool, *UserInfo, error) {
	conn, err := s.GetLdapConn()
	if err != nil {
		return false, nil, err
	}

	// 首先使用bind DN进行绑定（如果需要）
	if s.bindDN != "" && s.bindPassword != "" {
		if err := conn.Bind(s.bindDN, s.bindPassword); err != nil {
			return false, nil, err
		}
	}

	// 搜索用户
	searchRequest := ldap.NewSearchRequest(
		s.baseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		s.formatFilter(username),
		[]string{"dn", "cn", "mail", "memberOf"},
		nil,
	)

	sr, err := conn.Search(searchRequest)
	if err != nil {
		return false, nil, err
	}

	if len(sr.Entries) != 1 {
		return false, nil, ldap.NewError(ldap.LDAPResultInvalidCredentials, err)
	}

	userDN := sr.Entries[0].DN

	// 使用用户DN和密码进行绑定验证
	if err := conn.Bind(userDN, password); err != nil {
		return false, nil, err
	}

	// 获取用户信息
	userInfo := &UserInfo{
		Username: username,
		DN:       userDN,
	}

	if mail := sr.Entries[0].GetAttributeValue("mail"); mail != "" {
		userInfo.Email = mail
	}

	// 获取用户组
	if groups := sr.Entries[0].GetAttributeValues("memberOf"); len(groups) > 0 {
		userInfo.Groups = groups
	}

	return true, userInfo, nil
}

// formatFilter 格式化Filter
func (s *LdapInstance) formatFilter(username string) string {
	return fmt.Sprintf(s.filter, username)
}

// extractIPFromLDAPURL 返回服务器ip
func extractIPFromLDAPURL(ldapURL string) (string, error) {
	u, err := url.Parse(ldapURL)
	if err != nil {
		return "", err
	}

	host := u.Host
	// 处理可能没有端口的情况
	if strings.Contains(host, ":") {
		return strings.Split(host, ":")[0], nil
	}
	return host, nil
}
