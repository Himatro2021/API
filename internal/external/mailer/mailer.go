package mailer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Himatro2021/API/internal/config"
	"github.com/Himatro2021/API/internal/helper"
	"github.com/Himatro2021/API/internal/model"
	"github.com/kumparan/go-utils"
	"github.com/sirupsen/logrus"
)

// Mailer use this if you needed mailing capability
type Mailer struct {
	userRepo model.UserRepository
}

// NewMailer returns new mailer instance
func NewMailer(userRepo model.UserRepository) *Mailer {
	return &Mailer{
		userRepo,
	}
}

// SendInvitationEmail send invitation based on it's invitation model
func (m *Mailer) SendInvitationEmail(invitation *model.UserInvitation) {
	subject := "Undangan untuk menjadi Pengguna Website Himatro Unila"
	content := m.generateInvitationEmailContent(invitation.GenerateInvitationLink())

	payload := map[string]string{
		"receipient_email": invitation.Email,
		"receipient_name":  invitation.Name,
		"sender_email":     "himatro@unila.ac.id",
		"sender_name":      "Himatro Server",
		"subject":          subject,
		"content":          content,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		logrus.Error("payload error:", err.Error())
		return
	}

	resp, err := http.Post(config.MailServiceURL(), "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		logrus.Error("error when requesting to send invitation email:", err.Error())
		helper.LogIfErr(m.userRepo.MarkInvitationStatus(context.Background(), invitation, model.InvitationStatusFailed))
		return
	}
	defer helper.WrapCloser(resp.Body.Close)

	if err := invitation.Encrypt(); err != nil {
		logrus.Error("Err when encrypting invitation data after send email: ", utils.Dump(invitation))
	}

	if resp.Status != "200 OK" {
		logrus.Error("mail service returned not OK status: ", resp)
		helper.LogIfErr(m.userRepo.MarkInvitationStatus(context.Background(), invitation, model.InvitationStatusFailed))
		return
	}

	bodyByte, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.Error("failed to ready body response from mail service", err.Error())
	}

	m.addMailServiceIDToMailingList(bodyByte, invitation)
	helper.LogIfErr(m.userRepo.MarkInvitationStatus(context.Background(), invitation, model.InvitationStatusSent))
}

func (m *Mailer) generateInvitationEmailContent(link string) string {
	return fmt.Sprintf(`
		<html>
			<h3>Selamat Anda Telah Menerima Undangan untuk Menjadi Pengguna Website Himatro</h3>
			<p>Silahkan klik link berikut ini untuk membuat akun anda di website Himatro</p>
			<a href="%s">klik link ini</a>
		</html>
	`, link)
}

func (m *Mailer) addMailServiceIDToMailingList(resp []byte, invitation *model.UserInvitation) {
	type target struct {
		ID int64 `json:"id"`
	}

	data := &target{}
	helper.PanicIfErr(json.Unmarshal(resp, data))

	invitation.MailServiceID = data.ID
}
