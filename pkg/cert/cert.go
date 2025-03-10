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

	"github.com/grengojbo/kubercert/pkg/shell"
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
	DryRun     bool
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

// RenewCert renew the certificate of the host
func (h *HostInfo) RenewCert(command string) {
	if len(command) != 0 {
		log.Infoln("Start renewing the certificate...")
		result, err := shell.Run(command, false, false, "")
		if err != nil {
			log.Fatalf("Command failed: %s", err.Error())
		}
		log.Infof("command: %s", command)
		log.Infof("responce: %s", result)
	} else {
		distro := h.GetKubernetesDistributive()
		switch distro {
		case "k3s":
			h.RenewCertK3s()
		default:
			log.Errorf("Unsupported distro: %s", distro)
		}
	}
}

// RenewCertK3s renew the certificate of the host for k3s
func (h *HostInfo) RenewCertK3s() {
	log.Infof("Renewing certificate for k3s")

	dryRunResponce := ""
	if h.DryRun {
		dryRunResponce = "ok..."
	}

	command := "systemctl stop k3s.service"
	log.Infoln("Stop k3s server")
	_, err := shell.Run(command, false, false, dryRunResponce)
	if err != nil {
		log.Fatalf("Command failed: %s", err.Error())
	}

	command = "systemctl start k3s.service"
	log.Infoln("start k3s server")
	_, err = shell.Run(command, false, false, dryRunResponce)
	if err != nil {
		log.Fatalf("Command failed: %s", err.Error())
	}

	// log.Infoln("Wait for k3s server to be ready")
	log.Infoln("Suncessfully renewed certificate for k3s")
}
