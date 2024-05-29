package email

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"math/rand"
	"monitor/internal/svc"
	"monitor/internal/types"
	"net/smtp"
	"strings"
)

type EmailUtilLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewEmailUtilLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EmailUtilLogic {
	return &EmailUtilLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *EmailUtilLogic) SendMail(EmailInfo *types.EmailInfo, address []string, subject string, body string) (err error) {
	for _, emailUserInfo := range EmailInfo.EmailUser {
		auth := smtp.PlainAuth("", emailUserInfo.User, emailUserInfo.Password, EmailInfo.Host)
		contentType := "Content-Type: text/html; charset=UTF-8"
		for _, v := range address {
			s := fmt.Sprintf("To:%s\r\nFrom:%s<%s>\r\nSubject:%s\r\n%s\r\n\r\n%s",
				v, emailUserInfo.NickName, emailUserInfo.User, subject, contentType, body)
			msg := []byte(s)
			addr := fmt.Sprintf("%s:%s", EmailInfo.Host, EmailInfo.Port)
			err = smtp.SendMail(addr, auth, emailUserInfo.User, []string{v}, msg)
			if err != nil {
				fmt.Println(err.Error(), "发送邮件产生了异常 ************************************************************")
				return err
			}
		}
	}
	return
}

func (l *EmailUtilLogic) SendMailRandom(EmailInfo *types.EmailInfo, address []string, subject string, body string) (err error) {
	body = fmt.Sprintf("已同步发送给多人: %s  邮件内容: %s", strings.Join(address, ","), body)
	emailUserInfo := EmailInfo.EmailUser[rand.Intn(len(EmailInfo.EmailUser))]
	auth := smtp.PlainAuth("", emailUserInfo.User, emailUserInfo.Password, EmailInfo.Host)
	contentType := "Content-Type: text/html; charset=UTF-8"
	for _, v := range address {
		s := fmt.Sprintf("To:%s\r\nFrom:%s<%s>\r\nSubject:%s\r\n%s\r\n\r\n%s",
			v, emailUserInfo.NickName, emailUserInfo.User, subject, contentType, body)
		msg := []byte(s)
		addr := fmt.Sprintf("%s:%s", EmailInfo.Host, EmailInfo.Port)
		err = smtp.SendMail(addr, auth, emailUserInfo.User, []string{v}, msg)
		if err != nil {
			fmt.Println(err.Error(), "发送邮件产生了异常 ************************************************************")
			return err
		}
	}
	return
}
