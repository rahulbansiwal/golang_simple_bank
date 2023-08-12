package mail

import (
	"simple_bank/db/util"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSendEmailWithGmail(t *testing.T) {
	if testing.Short(){
		t.Skip()
	}

	config, err := util.LoadConfig("../")
	require.NoError(t, err)

	sender := NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)

	subject := "a test email"
	content := `
		<h1> Hello World</h1>
		<p>This is a test email </p>
		`
	to := []string{"rahulbnswl7@gmail.com"}
	//attachFiles := []string{"../go.sum"}

	err = sender.SendEmail(subject, content, to, nil, nil,nil)
	require.NoError(t, err)

}
