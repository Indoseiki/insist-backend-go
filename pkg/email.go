package pkg

import (
	"fmt"

	"gopkg.in/gomail.v2"
)

type EmailSender struct {
	SMTPHost string
	SMTPPort int
	Email    string
	Password string
}

func NewEmailSender(host string, port int, email, password string) *EmailSender {
	return &EmailSender{
		SMTPHost: host,
		SMTPPort: port,
		Email:    email,
		Password: password,
	}
}

func (es *EmailSender) SendEmail(to string, subject string, typeBody string, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", es.Email)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody(typeBody, body)

	d := gomail.NewDialer(es.SMTPHost, es.SMTPPort, es.Email, es.Password)

	return d.DialAndSend(m)
}

func TemplateSendOTP(name string, username string, code string, token string) string {
	html := fmt.Sprintf(`
	<div
      style="
        display: flex;
        justify-content: center;
        align-items: center;
        font-family: Arial, sans-serif;
        background-color: #f8fbff;
      "
    >
    <table
      width="700"
      border="0"
      align="center"
      cellpadding="0"
      cellspacing="0"
      style="
        background-color: #fff;
        border-collapse: collapse;
        box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
      "
    >
      <tbody>
        <tr>
          <td align="center" valign="middle" style="padding: 33px 0">
            <img
              src="https://images.glints.com/unsafe/720x0/glints-dashboard.s3.amazonaws.com/company-logo/6f02624d065682496253c50641c72d29.png"
              alt="Company Logo"
              width="150"
              style="border: 0"
            />
          </td>
        </tr>
        <tr>
          <td style="padding: 0 30px; background: #fff">
            <table width="100%%" border="0" cellspacing="0" cellpadding="0">
              <tbody>
                <tr>
                  <td
                    style="
                      border-bottom: 3px solid #e6e6e6;
                      font-size: 18px;
                      padding: 20px 0 5px;
					  font-weight: bold;
                    "
                  >       
					Reset Password & Aktivasi 2FA
                  </td>
                </tr>

                <tr>
                  <td
                    style="
                      font-size: 14px;
                      line-height: 30px;
                      color: #666;
                      padding: 20px 0;
                    "
                  >
                    <strong>Hai, %s!</strong>
                    <br />
                    Anda kini telah resmi terdaftar sebagai vendor PT. Indoseiki
                    Metalutama di sistem INSIST (Indoseiki Integration System).
                    Untuk keamanan, Anda diharuskan membuat password baru dengan 
					melakukan reset melalui tautan berikut:
					<a
                    	style="color: #000080; text-decoration: none"
                        href="http://localhost:5173/reset-password/%s"
                        target="_blank"
                        ><strong>Reset Password</strong></a
                    >
                  </td>
                </tr>

				<tr>
                    <td style="font-size: 14px; line-height: 30px; color: #666">
                      Detail akun Anda sebagai berikut:
                    </td>
                  </tr>

                <tr style="display: flex; justify-content: center">
                  <td style="font-size: 14px; line-height: 30px; color: #666">
                    <table
                      width="100%%"
                      border="0"
                      cellpadding="0"
                      cellspacing="0"
                    >
                      <tr>
                        <td style="width: 70px; font-size: 14px; color: #666">
                          <strong>Username</strong>
                        </td>
                        <td style="width: 10px; font-size: 14px; color: #666">
                          :
                        </td>
                        <td style="font-size: 14px; color: #666">%s</td>
                      </tr>
                      <tr>
                        <td style="width: 70px; font-size: 14px; color: #666">
                          <strong>Code</strong>
                        </td>
                        <td style="width: 10px; font-size: 14px; color: #666">
                          :
                        </td>
                        <td style="font-size: 14px; color: #666">%s</td>
                      </tr>
                    </table>
                  </td>
                </tr>

				<tr>
                    <td align="center" style="padding-top: 20px">
                      <img
                        src="https://insist.indoseiki.com:4004/images/20241011_rizki.png"
                        alt="QR Code"
                        style="width: 150px; height: 150px; border: 0"
                      />
                    </td>
                </tr>

                <tr>
                  <td style="font-size: 14px; line-height: 30px; color: #666; padding: 20px 0;">
                    Untuk meningkatkan keamanan akun Anda, silakan scan kode QR
                    yang telah kami kirimkan menggunakan aplikasi <strong>Google
                    Authenticator</strong>, atau masukkan kode yang tersedia untuk
                    mengaktifkan autentikasi dua faktor (2FA). Ini akan
                    menambahkan lapisan perlindungan ekstra pada akun Anda.
                  </td>
                </tr>

				<tr>
                  <td
                    style="
                      font-size: 14px;
                      line-height: 30px;
                      color: #666;
                    "
                  >
                    Silakan login melalui sistem:
                    <a
                      style="color: #000080; text-decoration: none"
                      href="https://insist.indoseiki.com:4004/"
                      target="_blank"
                      ><strong>INSIST (Indoseiki Integration System)</strong></a
                    >
                  </td>
                </tr>

                <tr>
                  <td
                    style="
                      font-size: 14px;
                      line-height: 30px;
                      color: #666;
					  padding: 20px 0;
                    "
                  >
                    Jika mengalami kendala saat login, Anda dapat menghubungi
                    Administrator melalui email
                    <a
                      style="color: #000080; text-decoration: none"
                      href="mailto:email@example.com"
                      >email@example.com</a
                    >
                    atau
                    <a
                      style="color: #000080; text-decoration: none"
                      href="mailto:email@example.com"
                      >email@example.com</a
                    >
                    atau telepon ke
                    <a
                      style="color: #000080; text-decoration: none"
                      href="tel:+6281234567890"
                      >+62 812-3456-7890</a
                    >
                    .
                  </td>
                </tr>

                <tr>
                  <td
                    style="
                      font-size: 14px;
                      line-height: 30px;
                      color: #666;
                      padding: 20px 0;
                    "
                  >
                    <span style="font-size: 14px; color: red; font-weight: bold"
                      >Peringatan Keamanan:</span
                    >
                    <ul>
                      <li>Jaga kerahasiaan detail akun Anda.</li>
                      <li>
                        Jangan berbagi informasi akun Anda dengan siapapun.
                      </li>
                      <li>Rutinlah mengganti password untuk keamanan.</li>
                    </ul>
                  </td>
                </tr>

                <tr>
                  <td
                    style="
                      padding: 10px 0 15px 0;
                      font-size: 12px;
                      color: #999;
                      line-height: 20px;
                    "
                  >
                    Email ini terkirim otomatis.
                  </td>
                </tr>
              </tbody>
            </table>
          </td>
        </tr>
        <tr>
          <td
            align="center"
            style="font-size: 12px; color: #999; padding: 20px 0"
          >
            PT. Indoseiki Metalutama
            <br />
            Website Resmi:
            <a
              style="color: #999; text-decoration: none"
              href="https://www.indoseiki.com/"
              target="_blank"
              >www.indoseiki.com</a
            >
          </td>
        </tr>
      </tbody>
    </table>
    </div>
`, name, token, username, code)

	return html
}
