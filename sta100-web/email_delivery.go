package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"mime"
	"net"
	"net/http"
	"net/smtp"
	"net/textproto"
	"regexp"
	"strings"
	"time"
)

type EmailSettings struct {
	Host      string `json:"host"`
	Port      int    `json:"port"`
	Security  string `json:"security"`
	Username  string `json:"username"`
	Password  string `json:"password,omitempty"`
	FromName  string `json:"fromName"`
	FromEmail string `json:"fromEmail"`
}

type EmailAttachment struct {
	Name string
	Data []byte
	Type string
}

type EmailSendRequest struct {
	To           []string `json:"to"`
	Cc           []string `json:"cc,omitempty"`
	Subject      string   `json:"subject"`
	Body         string   `json:"body"`
	Language     string   `json:"language"`
	DocumentIDs  []string `json:"documentIds,omitempty"`
	AttachRecord bool     `json:"attachRecord"`
}

var emailAddressPattern = regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)

func defaultEmailSettings() EmailSettings {
	return EmailSettings{Host: "", Port: 587, Security: "starttls", FromName: "Stratronix"}
}

func (a *businessAPI) emailSettings(ctx context.Context) EmailSettings {
	settings := defaultEmailSettings()
	_ = a.store.getSetting(ctx, "email", &settings)
	return settings
}

func publicEmailSettings(settings EmailSettings) map[string]any {
	return map[string]any{"host": settings.Host, "port": settings.Port, "security": settings.Security, "username": settings.Username, "fromName": settings.FromName, "fromEmail": settings.FromEmail, "configured": settings.Host != "" && settings.FromEmail != "" && settings.Password != ""}
}

func (a *businessAPI) emailRouter(w http.ResponseWriter, r *http.Request, parts []string) {
	if len(parts) == 0 || parts[0] == "" {
		if !allowMethod(w, r, "GET") {
			return
		}
		writeJSON(w, 200, publicEmailSettings(a.emailSettings(r.Context())))
		return
	}
	if parts[0] == "test" {
		if !allowMutation(w, r, "POST") {
			return
		}
		var request EmailSendRequest
		if err := decodeJSONBody(w, r, &request); err != nil {
			return
		}
		settings := a.emailSettings(r.Context())
		if err := validateEmailSettings(settings); err != nil {
			writeAPIError(w, 400, "EMAIL_NOT_CONFIGURED", err.Error())
			return
		}
		to := request.To
		if len(to) == 0 {
			to = []string{settings.FromEmail}
		}
		if err := sendSMTPMessage(settings, to, nil, "STA-100 SMTP 连通性测试", "STA-100 SMTP connection test\n\n"+time.Now().Format(time.RFC3339)); err != nil {
			writeAPIError(w, 502, "EMAIL_TEST_FAILED", err.Error())
			return
		}
		writeJSON(w, 200, map[string]any{"ok": true, "message": "SMTP 测试邮件已发送"})
		return
	}
	if !allowMethod(w, r, "PATCH") {
		return
	}
	var settings EmailSettings
	if err := decodeJSONBody(w, r, &settings); err != nil {
		return
	}
	if settings.Port == 0 {
		settings.Port = 587
	}
	if settings.Security == "" {
		settings.Security = "starttls"
	}
	old := a.emailSettings(r.Context())
	if settings.Password == "" {
		settings.Password = old.Password
	}
	if settings.Host == "" || settings.FromEmail == "" {
		writeAPIError(w, 400, "EMAIL_SETTINGS_INVALID", "SMTP 主机和发件人邮箱不能为空")
		return
	}
	if err := validateEmailSettings(settings); err != nil {
		writeAPIError(w, 400, "EMAIL_SETTINGS_INVALID", err.Error())
		return
	}
	if err := a.store.putSetting(r.Context(), "email", settings); err != nil {
		writeBusinessError(w, err)
		return
	}
	writeJSON(w, 200, map[string]any{"settings": publicEmailSettings(settings), "message": "邮件设置已保存"})
}

func validateEmailSettings(s EmailSettings) error {
	if s.Host == "" || s.FromEmail == "" || s.Password == "" {
		return fmt.Errorf("请填写 SMTP 主机、发件邮箱和应用专用密码")
	}
	if !emailAddressPattern.MatchString(s.FromEmail) {
		return fmt.Errorf("发件人邮箱格式无效")
	}
	if s.Port < 1 || s.Port > 65535 {
		return fmt.Errorf("SMTP 端口无效")
	}
	if s.Security != "starttls" && s.Security != "tls" && s.Security != "none" {
		return fmt.Errorf("SMTP 安全方式仅支持 STARTTLS、TLS 或无加密")
	}
	return nil
}

func (a *businessAPI) sendQuoteEmail(ctx context.Context, item Quote, request EmailSendRequest) error {
	settings := a.emailSettings(ctx)
	if err := validateEmailSettings(settings); err != nil {
		return err
	}
	if len(request.To) == 0 {
		return fmt.Errorf("至少需要一个收件人")
	}
	attachments := []EmailAttachment{}
	if request.AttachRecord {
		tpl, err := a.templateForQuote(ctx, item)
		if err != nil {
			return err
		}
		values := a.quoteTemplateValues(ctx, item, tpl)
		data, name, err := renderDynamicQuotePDF(item, values, "报价单-"+item.ID)
		if err != nil {
			return err
		}
		attachments = append(attachments, EmailAttachment{Name: name, Data: data, Type: "application/pdf"})
	}
	return sendSMTPMessage(settings, request.To, request.Cc, request.Subject, request.Body, attachments...)
}

func (a *businessAPI) sendOrderEmail(ctx context.Context, item Order, request EmailSendRequest) error {
	settings := a.emailSettings(ctx)
	if err := validateEmailSettings(settings); err != nil {
		return err
	}
	if len(request.To) == 0 {
		return fmt.Errorf("至少需要一个收件人")
	}
	attachments := []EmailAttachment{}
	for _, id := range request.DocumentIDs {
		var doc Document
		if err := a.store.get(ctx, "documents", id, &doc); err != nil {
			return fmt.Errorf("单据 %s 不存在", id)
		}
		tpl, err := a.templateForDocument(ctx, doc)
		if err != nil {
			return err
		}
		values := a.documentTemplateValues(ctx, doc, tpl)
		data, name, contentType, err := a.renderDocumentDownload(ctx, doc, tpl, values, "pdf")
		if err != nil {
			return err
		}
		attachments = append(attachments, EmailAttachment{Name: name, Data: data, Type: contentType})
	}
	return sendSMTPMessage(settings, request.To, request.Cc, request.Subject, request.Body, attachments...)
}

func sendSMTPMessage(settings EmailSettings, to, cc []string, subject, body string, attachments ...EmailAttachment) error {
	all := append(append([]string{}, to...), cc...)
	for _, address := range all {
		if !emailAddressPattern.MatchString(strings.TrimSpace(address)) {
			return fmt.Errorf("邮箱地址无效：%s", address)
		}
	}
	if subject == "" {
		return fmt.Errorf("邮件主题不能为空")
	}
	hostPort := net.JoinHostPort(settings.Host, fmt.Sprintf("%d", settings.Port))
	client, err := smtpClient(settings, hostPort)
	if err != nil {
		return err
	}
	defer client.Close()
	auth := smtp.PlainAuth("", settings.Username, settings.Password, settings.Host)
	if settings.Username != "" {
		if err := client.Auth(auth); err != nil {
			return fmt.Errorf("SMTP 鉴权失败：%w", err)
		}
	}
	from := settings.FromEmail
	if err := client.Mail(from); err != nil {
		return err
	}
	for _, address := range all {
		if err := client.Rcpt(address); err != nil {
			return err
		}
	}
	writer, err := client.Data()
	if err != nil {
		return err
	}
	message := buildMIMEMessage(settings, to, cc, subject, body, attachments)
	if _, err = writer.Write(message); err != nil {
		_ = writer.Close()
		return err
	}
	if err = writer.Close(); err != nil {
		return err
	}
	return client.Quit()
}

func smtpClient(settings EmailSettings, address string) (*smtp.Client, error) {
	if settings.Security == "tls" {
		conn, err := tls.Dial("tcp", address, &tls.Config{ServerName: settings.Host, MinVersion: tls.VersionTLS12})
		if err != nil {
			return nil, fmt.Errorf("SMTP TLS 连接失败：%w", err)
		}
		client, err := smtp.NewClient(conn, settings.Host)
		if err != nil {
			_ = conn.Close()
			return nil, err
		}
		return client, nil
	}
	client, err := smtp.Dial(address)
	if err != nil {
		return nil, fmt.Errorf("SMTP 连接失败：%w", err)
	}
	if settings.Security == "starttls" {
		if ok, _ := client.Extension("STARTTLS"); !ok {
			_ = client.Close()
			return nil, fmt.Errorf("SMTP 服务器不支持 STARTTLS")
		}
		if err := client.StartTLS(&tls.Config{ServerName: settings.Host, MinVersion: tls.VersionTLS12}); err != nil {
			_ = client.Close()
			return nil, err
		}
	}
	return client, nil
}

func buildMIMEMessage(settings EmailSettings, to, cc []string, subject, body string, attachments []EmailAttachment) []byte {
	boundary := fmt.Sprintf("sta100-%d", time.Now().UnixNano())
	from := settings.FromEmail
	if settings.FromName != "" {
		from = mime.QEncoding.Encode("UTF-8", settings.FromName) + " <" + settings.FromEmail + ">"
	}
	var out bytes.Buffer
	fmt.Fprintf(&out, "From: %s\r\nTo: %s\r\n", from, strings.Join(to, ", "))
	if len(cc) > 0 {
		fmt.Fprintf(&out, "Cc: %s\r\n", strings.Join(cc, ", "))
	}
	fmt.Fprintf(&out, "Subject: %s\r\nMIME-Version: 1.0\r\nContent-Type: multipart/mixed; boundary=%q\r\n\r\n", mime.QEncoding.Encode("UTF-8", subject), boundary)
	fmt.Fprintf(&out, "--%s\r\nContent-Type: text/plain; charset=UTF-8\r\nContent-Transfer-Encoding: 8bit\r\n\r\n%s\r\n", boundary, body)
	for _, attachment := range attachments {
		filename := mime.QEncoding.Encode("UTF-8", attachment.Name)
		fmt.Fprintf(&out, "--%s\r\n", boundary)
		h := textproto.MIMEHeader{}
		h.Set("Content-Type", attachment.Type)
		h.Set("Content-Disposition", `attachment; filename="`+filename+`"`)
		h.Set("Content-Transfer-Encoding", "base64")
		for key, values := range h {
			for _, value := range values {
				fmt.Fprintf(&out, "%s: %s\r\n", key, value)
			}
		}
		out.WriteString("\r\n")
		encoded := make([]byte, base64.StdEncoding.EncodedLen(len(attachment.Data)))
		base64.StdEncoding.Encode(encoded, attachment.Data)
		for len(encoded) > 0 {
			n := 76
			if len(encoded) < n {
				n = len(encoded)
			}
			out.Write(encoded[:n])
			out.WriteString("\r\n")
			encoded = encoded[n:]
		}
	}
	fmt.Fprintf(&out, "--%s--\r\n", boundary)
	return out.Bytes()
}
