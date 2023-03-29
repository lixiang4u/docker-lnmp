package model

type VirtualHost struct {
	Id      string `yaml:"id" json:"id" toml:"id"`                   // ID全局唯一
	Name    string `yaml:"name" json:"name" toml:"name"`             // 不带特殊字符的名称，全局唯一
	Domain  string `yaml:"domain" json:"domain" toml:"domain"`       // 虚拟主机域名
	Root    string `yaml:"root" json:"root" toml:"root"`             // 项目根目录（本地机器目录），用于docker的volumes映射
	WebRoot string `yaml:"web_root" json:"web_root" toml:"web_root"` // 项目web服务的根目录（本地机器目录），相对于root位置，一般为："/", "/public"这类
	Port    int    `yaml:"port" json:"port" toml:"port"`             // 虚拟主机端口
}
