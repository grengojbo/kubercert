package cert

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strconv"
	"text/template"
	"time"

	log "github.com/sirupsen/logrus"
)

// var timeout = 5 * time.Second

type HostInfo struct {
	Host       string
	Port       int
	Certs      []*x509.Certificate
	Expired    string
	ExpiredAt  time.Time
	ExpireDays int
}

// GetCerts returns the certificates of the host
func (h *HostInfo) GetCerts(timeout int) error {
	t := time.Duration(timeout) * time.Second
	// log.Printf("connecting to %s:%d", h.Host, h.Port)
	dialer := &net.Dialer{Timeout: t}
	conn, err := tls.DialWithDialer(
		dialer,
		"tcp",
		h.Host+":"+strconv.Itoa(h.Port),
		&tls.Config{
			InsecureSkipVerify: true,
		})
	if err != nil {
		return err
	}

	defer conn.Close()

	if err := conn.Handshake(); err != nil {
		return err
	}

	pc := conn.ConnectionState().PeerCertificates
	h.Certs = make([]*x509.Certificate, 0, len(pc))
	for _, cert := range pc {
		if cert.IsCA {
			continue
		}
		h.Certs = append(h.Certs, cert)
	}

	expireDate := time.Duration(h.ExpireDays) * 24 * time.Hour
	h.ExpiredAt = time.Now().Add(expireDate)
	dt := h.Certs[0].NotAfter.Sub(h.ExpiredAt)
	h.Expired = FmtDuration(dt)

	return nil
}

// IsExpired returns true if the certificate is expired.
func (h *HostInfo) IsExpired() bool {
	expireDate := time.Duration(h.ExpireDays) * 24 * time.Hour
	if len(h.Certs) == 0 {
		return false
	}
	return h.Certs[0].NotAfter.Before(time.Now().Add(expireDate))
}

// GetKubernetesDistributive returns the kubernetes distributive name.
func (h *HostInfo) GetKubernetesDistributive() (dist string) {
	dist = h.Certs[0].Subject.CommonName
	if len(dist) == 0 {
		return "kubernetes"
	}
	return dist
}

// show the certificates of the host
func (h *HostInfo) ShowCerts(mode string) (err error) {

	// deadline := time.Now().Add(expires)

	switch mode {
	case "json":
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		enc.SetEscapeHTML(false)
		err = enc.Encode(&h)
		// errs.Push(err)

	case "text":
		t := template.Must(template.New("").Parse(`
Host: {{ .Host }}:{{ .Port }}
ExpiredAt: {{ .ExpiredAt.Format "Jan 2, 2006 15:04" }}
Expired: {{ .Expired }}
Start expire days: {{ .ExpireDays }}
Certs:
	{{ range .Certs -}}
	Issuer: {{ .Issuer.CommonName }}
	Subject: {{ .Subject.CommonName }}
	Not Before: {{ .NotBefore.Format "Jan 2, 2006 15:04" }}
	Not After: {{ .NotAfter.Format "Jan 2, 2006 15:04" }}
	DNS names: {{ range .DNSNames }}{{ . }} {{ end }}
{{ end }}	
`))
		err = t.Execute(os.Stdout, &h)
		// errs.Push(err)

	case "none":
	}
	return err
}

// FmtDuration returns a string representing the duration in the form "72h3m".
func FmtDuration(d time.Duration) string {
	if d < 0 {
		return "-" + FmtDuration(-d)
	}
	d = d.Round(time.Minute)
	day := int(d.Hours() / 24)
	h := (d / time.Hour) % 24
	d -= h * time.Hour
	m := (d / time.Minute) % 60
	if day > 0 {
		return fmt.Sprintf("%dd %dh %dm", day, h, m)
	}
	return fmt.Sprintf("%dh %dm", h, m)
	// return fmt.Sprintf("%02dh %02dm", h, m)
}

// renew the certificate of the host
func (h *HostInfo) ReNew(command string) error {
	if len(command) != 0 {
		log.Infoln("Start renewing the certificate...")
		return nil
		// return fmt.Errorf("no command specified")
	}
	distro := h.GetKubernetesDistributive()
	log.Infof("Renewing certificate for %s", distro)
	return nil
}
