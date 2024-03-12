/*
@Time: 2020-07-23 15:25
@Auth: xujiang
@File: proxy.go
@Software: GoLand
@Desc: TODO
*/
package config

var ProxyConfig proxyConfig

type proxyConfig struct {
	NciicProxy           string `goblet:"nciic_proxy,"`
	TaofuNciicProxy        string `goblet:"taofu_nciic_proxy,"`
}