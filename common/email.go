package common

import (
	"log"

	"github.com/luxingwen/go-realworld/config"
	gomail "gopkg.in/gomail.v2"
)

/***
  *  发送者信息
***/
type Mail struct {
	Address string //收件人信息
	Name    string //收件人昵称
	Subject string //邮件标题
	Content string //邮件内容
}

func SendEmail(ma Mail) {
	d := gomail.NewDialer(config.EmaiConf.DefaultAdress, config.EmaiConf.DefaultPort, config.EmaiConf.DefaultUser, config.EmaiConf.DefaultPasswd)
	server, err := d.Dial()
	if err != nil {
		panic(err)
	}

	m := gomail.NewMessage()

	m.SetHeader("From", config.EmaiConf.DefaultUser)
	m.SetHeader("To", ma.Address)
	m.SetAddressHeader("To", ma.Address, ma.Name)
	m.SetHeader("Subject", ma.Subject)
	// m.SetBody("text/html", "<style type='text/css'>body{min-width: 800px;}#email{width: 600px;margin: auto;}a{color: #286090;text-decoration: none}#title{color: #4cae4c;font-size: x-large}#button{display: inline-block;height: 30px;width: 100px;background-color: #4cae4c;color: white;text-align: center;line-height: 30px;text-decoration: none; border-radius: 6px}</style><div id='email'><p id='title'>请确认您的邮箱地址</p><p>点击下方按钮或者复制链接进行邮箱验证：</p><p><a href='http://x.biggerforum.org/'>http://x.biggerforum.org/</a></p><a id='button' href='http://x.biggerforum.org/'>马上验证</a></div>")
	m.SetBody("text/html", ma.Content)
	if err := gomail.Send(server, m); err != nil {
		log.Printf("Could not send email to %q: %v", ma.Address, err)
	}
	m.Reset()
}
