package tls

import (
	"crypto/tls"
	"fmt"
	"time"

	"github.com/doyensec/wallet-info/utils"
)

type TlsInfo struct {
	IssuedOn  time.Time
	ExpiresOn time.Time
}

func (tls *TlsInfo) IsValid() bool {
	now := time.Now()
	return tls.IssuedOn.Before(now) && tls.ExpiresOn.After(now)
}

func GetInfo(domain string) (*TlsInfo, error) {
	utils.Logger.Infow("getting tls info", "domain", domain)

	conn, err := tls.Dial("tcp", fmt.Sprintf("%v:443", domain), nil)
	if err != nil {
		utils.Logger.Errorw("dial failed", err)
		return nil, err
	}

	err = conn.VerifyHostname(domain)
	if err != nil {
		utils.Logger.Errorw("hostname verification failed", err)
		return nil, err
	}

	info := &TlsInfo{
		IssuedOn:  conn.ConnectionState().PeerCertificates[0].NotBefore,
		ExpiresOn: conn.ConnectionState().PeerCertificates[0].NotAfter,
	}
	return info, nil
}
