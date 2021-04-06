package email

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"log"
	"os"
	"reflect"

	//go get -u github.com/aws/aws-sdk-go
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/joho/godotenv"
	gomail "gopkg.in/gomail.v2"
)

const (
    // The character encoding for the email.
    CharSet = "UTF-8"
	defaultTemplate = `
	<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN"
    "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html xmlns="http://www.w3.org/1999/xhtml">
<head>
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <meta http-equiv="Content-Type" content="text/html; charset=UTF-8">
  <meta name="color-scheme" content="light">
  <meta name="supported-color-schemes" content="light">
</head>
<body
    style="box-sizing: border-box; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif, 'Apple Color Emoji', 'Segoe UI Emoji', 'Segoe UI Symbol'; position: relative; -webkit-text-size-adjust: none; background-color: #ffffff; color: #718096; height: 100%; line-height: 1.4; margin: 0; padding: 0; width: 100% !important;">
<style>
    @media only screen and (max-width: 600px) {
        .inner-body {
            width: 100% !important;
        }

        .footer {
            width: 100% !important;
        }
    }

    @media only screen and (max-width: 500px) {
        .button {
            width: 100% !important;
        }
    }
</style>

<table class="wrapper" width="100%" cellpadding="0" cellspacing="0" role="presentation"
       style="box-sizing: border-box; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif, 'Apple Color Emoji', 'Segoe UI Emoji', 'Segoe UI Symbol'; position: relative; -premailer-cellpadding: 0; -premailer-cellspacing: 0; -premailer-width: 100%; background-color: #edf2f7; margin: 0; padding: 0; width: 100%;">
  <tbody>
  <tr>
    <td align="center"
        style="box-sizing: border-box; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif, 'Apple Color Emoji', 'Segoe UI Emoji', 'Segoe UI Symbol'; position: relative;">
      <table class="content" width="100%" cellpadding="0" cellspacing="0" role="presentation"
             style="box-sizing: border-box; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif, 'Apple Color Emoji', 'Segoe UI Emoji', 'Segoe UI Symbol'; position: relative; -premailer-cellpadding: 0; -premailer-cellspacing: 0; -premailer-width: 100%; margin: 0; padding: 0; width: 100%;">
        <tbody>
        <tr>
          <td class="header"
              style="box-sizing: border-box; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif, 'Apple Color Emoji', 'Segoe UI Emoji', 'Segoe UI Symbol'; position: relative; padding: 25px 0; text-align: center;">
			{{with .HeaderImageLink}}
			  <img src="{{.}}" class="logo" alt="Logo"
                   style="box-sizing: border-box; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif, 'Apple Color Emoji', 'Segoe UI Emoji', 'Segoe UI Symbol'; position: relative; max-width: 100%; border: none; height: 75px; width: 75px;">
            {{end}}
			{{with .HeaderTitle}}
				<div style="box-sizing: border-box; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif, 'Apple Color Emoji', 'Segoe UI Emoji', 'Segoe UI Symbol'; position: relative;">
                	{{.}}
              	</div>
			{{end}}
            </a>
          </td>
        </tr>

        <!-- Email Body -->
        <tr>
          <td class="body" width="100%" cellpadding="0" cellspacing="0"
              style="box-sizing: border-box; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif, 'Apple Color Emoji', 'Segoe UI Emoji', 'Segoe UI Symbol'; position: relative; -premailer-cellpadding: 0; -premailer-cellspacing: 0; -premailer-width: 100%; background-color: #edf2f7; border-bottom: 1px solid #edf2f7; border-top: 1px solid #edf2f7; margin: 0; padding: 0; width: 100%;">
            <table class="inner-body" align="center" width="570" cellpadding="0" cellspacing="0" role="presentation"
                   style="box-sizing: border-box; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif, 'Apple Color Emoji', 'Segoe UI Emoji', 'Segoe UI Symbol'; position: relative; -premailer-cellpadding: 0; -premailer-cellspacing: 0; -premailer-width: 570px; background-color: #ffffff; border-color: #e8e5ef; border-radius: 2px; border-width: 1px; box-shadow: 0 2px 0 rgba(0, 0, 150, 0.025), 2px 4px 0 rgba(0, 0, 150, 0.015); margin: 0 auto; padding: 0; width: 570px;">
              <!-- Body content -->
              <tbody>
              <tr>
                <td class="content-cell"
                    style="box-sizing: border-box; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif, 'Apple Color Emoji', 'Segoe UI Emoji', 'Segoe UI Symbol'; position: relative; max-width: 100vw; padding: 32px;">
                  <h1 style="box-sizing: border-box; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif, 'Apple Color Emoji', 'Segoe UI Emoji', 'Segoe UI Symbol'; position: relative; color: #3d4852; font-size: 18px; font-weight: bold; margin-top: 0; text-align: left;">
                    {{.Greeting}}
                  </h1>
                  {{range .Items}}
				  	{{$instanceOf := instanceOf .}}
                    {{ if eq $instanceOf "Paragraph"}}
                      <p style="box-sizing: border-box; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif, 'Apple Color Emoji', 'Segoe UI Emoji', 'Segoe UI Symbol'; position: relative; font-size: 16px; line-height: 1.5em; margin-top: 0; text-align: left;">
                        {{.Text}}
                      </p>
                    {{ else if eq $instanceOf "Button"}}
                      <table class="action" align="center" width="100%" cellpadding="0" cellspacing="0" role="presentation"
                         style="box-sizing: border-box; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif, 'Apple Color Emoji', 'Segoe UI Emoji', 'Segoe UI Symbol'; position: relative; -premailer-cellpadding: 0; -premailer-cellspacing: 0; -premailer-width: 100%; margin: 30px auto; padding: 0; text-align: center; width: 100%;">
                    <tbody>
                    <tr>
                      <td align="center"
                          style="box-sizing: border-box; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif, 'Apple Color Emoji', 'Segoe UI Emoji', 'Segoe UI Symbol'; position: relative;">
                        <table width="100%" border="0" cellpadding="0" cellspacing="0" role="presentation"
                               style="box-sizing: border-box; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif, 'Apple Color Emoji', 'Segoe UI Emoji', 'Segoe UI Symbol'; position: relative;">
                          <tbody>
                          <tr>
                            <td align="center"
                                style="box-sizing: border-box; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif, 'Apple Color Emoji', 'Segoe UI Emoji', 'Segoe UI Symbol'; position: relative;">
                              <table border="0" cellpadding="0" cellspacing="0" role="presentation"
                                     style="box-sizing: border-box; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif, 'Apple Color Emoji', 'Segoe UI Emoji', 'Segoe UI Symbol'; position: relative;">
                                <tbody>
                                <tr>
                                  <td style="box-sizing: border-box; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif, 'Apple Color Emoji', 'Segoe UI Emoji', 'Segoe UI Symbol'; position: relative;">
                                    <a href="{{.Link}}" class="button button-primary"
                                       target="_blank" rel="noopener"
                                       style="box-sizing: border-box; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif, 'Apple Color Emoji', 'Segoe UI Emoji', 'Segoe UI Symbol'; position: relative; -webkit-text-size-adjust: none; border-radius: 4px; color: #fff; display: inline-block; overflow: hidden; text-decoration: none; background-color: #2d3748; border-bottom: 8px solid #2d3748; border-left: 18px solid #2d3748; border-right: 18px solid #2d3748; border-top: 8px solid #2d3748;">
									   {{.Text}}
									</a>
                                  </td>
                                </tr>
                                </tbody>
                              </table>
                            </td>
                          </tr>
                          </tbody>
                        </table>
                      </td>
                    </tr>
                    </tbody>
                  </table>
                    {{end}}
                  {{end}}
                    <p style="box-sizing: border-box; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif, 'Apple Color Emoji', 'Segoe UI Emoji', 'Segoe UI Symbol'; position: relative; font-size: 16px; line-height: 1.5em; margin-top: 0; text-align: left;">
                      {{range .Farewells}}
                        {{.}}<br>
                      {{end}}
                    </p>
                </td>
              </tr>
              </tbody>
            </table>
          </td>
        </tr>

        {{with .CopyRight}}
          <tr>
          <td style="box-sizing: border-box; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif, 'Apple Color Emoji', 'Segoe UI Emoji', 'Segoe UI Symbol'; position: relative;">
            <table class="footer" align="center" width="570" cellpadding="0" cellspacing="0" role="presentation"
                   style="box-sizing: border-box; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif, 'Apple Color Emoji', 'Segoe UI Emoji', 'Segoe UI Symbol'; position: relative; -premailer-cellpadding: 0; -premailer-cellspacing: 0; -premailer-width: 570px; margin: 0 auto; padding: 0; text-align: center; width: 570px;">
              <tbody>
              <tr>
                <td class="content-cell" align="center"
                    style="box-sizing: border-box; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif, 'Apple Color Emoji', 'Segoe UI Emoji', 'Segoe UI Symbol'; position: relative; max-width: 100vw; padding: 32px;">
                  <p style="box-sizing: border-box; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif, 'Apple Color Emoji', 'Segoe UI Emoji', 'Segoe UI Symbol'; position: relative; line-height: 1.5em; margin-top: 0; color: #b0adc5; font-size: 12px; text-align: center;">
                    {{.}}
                  </p>
                </td>
              </tr>
              </tbody>
            </table>
          </td>
        </tr>
        {{end}}
        </tbody>
      </table>
    </td>
  </tr>
  </tbody>
</table>

</body>
</html>
	`
)

type Message struct {
	Sender string
	Recipient string
	Subject string
	Greeting string
	HtmlBody string
	TextBody string
	Farewells []string
	Items []Item
	Template string
	HeaderTitle string
	HeaderImageLink string
	CopyRight string
	ReplyTo string
	RawData []byte
	Attachments []string
}



type Paragraph struct {
	Text string
}

type Button struct {
	Text string
	Link string
}

type Item interface{}

func (message *Message) WithGreeting(greeting string) *Message {
	message.Greeting = greeting

	return message
}

func (message *Message) WithSubject(subject string) *Message {
	message.Subject = subject
	return message
}

func (message *Message) WithHtmlBody(html string) *Message {
	message.HtmlBody = html
	return message
}

func (message *Message) WithParagraph(text string) *Message {
	message.Items = append(
		message.Items,
		Paragraph{Text: text},
	)

	return message
}

func (message *Message) WithButton(text string, link string) *Message {
	message.Items = append(
		message.Items, 
		Button{Text: text, Link: link},
	)

	return message
}

func (message *Message) WithFarewells(Farewells []string) *Message {
	message.Farewells = Farewells

	return message
}

func (message *Message) WithCopyRight(copyRight string) *Message {
	message.CopyRight = copyRight
	return message
}

func (message *Message) WithHeaderImageLink(headerImageLink string) *Message {
	message.HeaderImageLink = headerImageLink
	return message
}

func (message *Message) WithHeaderTitle(headerTitle string) *Message {
	message.HeaderTitle = headerTitle
	return message
}

func (message *Message) WithAttachment(filePath string) *Message {
	message.Attachments = append(
		message.Attachments,
		filePath,
	)

	return message
}

func instanceOf(item Item) string {
	return reflect.TypeOf(item).Name()
}

func (message *Message) SendTo(recipient string) error {
	message.Recipient = recipient

	if message.Template == "" {
		message.Template = defaultTemplate
	}

	if message.Sender == "" {
		message.Sender = GetEnvWithKey("EMAIL_FROM")
		if message.Sender == "" {
			return errors.New("Sender Email should be specified in env or in Message Struct")
		}
	}

	// Custom function map
    funcMap := template.FuncMap{
        "instanceOf": instanceOf,
    }

	// Create a new template and parse the email into it.
	t := template.Must(template.New("emailBody").Funcs(funcMap).Parse(message.Template))

	var htmlBody bytes.Buffer
	err := t.Execute(&htmlBody, message)

	message.HtmlBody = htmlBody.String()

	if err != nil {
		log.Println("executing template:", err)
		return err
	}

	msg := gomail.NewMessage()
	msg.SetHeader("From", message.Sender)
	msg.SetHeader("To", message.Recipient)
	msg.SetHeader("Subject", message.Subject)
	msg.SetBody("text/html", message.HtmlBody)

	for _, attachment := range message.Attachments {
		msg.Attach(attachment)
	}
	

	var emailRaw bytes.Buffer
	msg.WriteTo(&emailRaw)

	message.RawData = emailRaw.Bytes()

	send(message)

	return nil
}

//GetEnvWithKey : get env value
func GetEnvWithKey(key string) string {
	return os.Getenv(key)
}

func loadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
		os.Exit(1)
	}
}

func init() {
	loadEnv()
}

func send(message *Message) error {
	// Create a new session in the us-west-2 region.
    // Replace us-west-2 with the AWS Region you're using for Amazon SES.
    sess, err := session.NewSession(
		&aws.Config{
        	Region:aws.String(GetEnvWithKey("SES_REGION")),
			Credentials: credentials.NewStaticCredentials(
				GetEnvWithKey("SES_KEY"),
				GetEnvWithKey("SES_SECRET"),
				"", // a token will be created when the session it's used.
			),
		},
    )
    
    // Create an SES session.
    svc := ses.New(sess)
    
    // Assemble the email.
    input := &ses.SendRawEmailInput{
        Destinations: []*string{
            aws.String(message.Recipient),
        },
        RawMessage: &ses.RawMessage{
            Data: message.RawData,
        },
        Source: aws.String(message.Sender),
            // Uncomment to use a configuration set
            //ConfigurationSetName: aws.String(ConfigurationSet),
    }

    // Attempt to send the email.
    result, err := svc.SendRawEmail(input)
    
    // Display error messages if they occur.
    if err != nil {
        if aerr, ok := err.(awserr.Error); ok {
            switch aerr.Code() {
            case ses.ErrCodeMessageRejected:
                fmt.Println(ses.ErrCodeMessageRejected, aerr.Error())
            case ses.ErrCodeMailFromDomainNotVerifiedException:
                fmt.Println(ses.ErrCodeMailFromDomainNotVerifiedException, aerr.Error())
            case ses.ErrCodeConfigurationSetDoesNotExistException:
                fmt.Println(ses.ErrCodeConfigurationSetDoesNotExistException, aerr.Error())
            default:
                fmt.Println(aerr.Error())
            }
        } else {
            // Print the error, cast err to awserr.Error to get the Code and
            // Message from an error.
            fmt.Println(err.Error())
        }
    
        return err
    }
    
    fmt.Println("Email Sent to address: " + message.Recipient)
    fmt.Println(result)

	return nil
}