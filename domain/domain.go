package domain

import (
	"time"

	werr "github.com/doyensec/wallet-info/errors"
	"github.com/doyensec/wallet-info/utils"
	"github.com/likexian/whois"
	whoisparser "github.com/likexian/whois-parser"
)

type DomainRecordInfo struct {
	Created time.Time
	Updated time.Time
	Expires time.Time
}

func GetDomainRecordInfo(host string) (*DomainRecordInfo, error) {
	tldPlusOne, err := utils.GetTldPlusOne(host)
	if err != nil {
		utils.Logger.Errorw("error building tld+1", err)
		return nil, &werr.WhoIsLookupError{}
	}

	whois_raw, err := whois.Whois(tldPlusOne)
	if err != nil {
		utils.Logger.Errorw("whois lookup failed", err)
		return nil, &werr.WhoIsLookupError{}
	}

	parsed, err := whoisparser.Parse(whois_raw)
	if err != nil {
		utils.Logger.Errorw("whois parsing failed", err)
		return nil, &werr.WhoIsParsingError{}
	}

	created, err := time.Parse(time.RFC3339, parsed.Domain.CreatedDate)
	if err != nil {
		utils.Logger.Errorw("whois created date missing", err)
		return nil, &werr.WhoIsParsingError{}
	}

	updated, err := time.Parse(time.RFC3339, parsed.Domain.UpdatedDate)
	if err != nil {
		utils.Logger.Errorw("whois updated date missing", err)
		return nil, &werr.WhoIsParsingError{}
	}

	expires, err := time.Parse(time.RFC3339, parsed.Domain.ExpirationDate)
	if err != nil {
		utils.Logger.Errorw("whois expiration date missing", err)
		return nil, &werr.WhoIsParsingError{}
	}

	info := &DomainRecordInfo{
		Created: created,
		Updated: updated,
		Expires: expires,
	}

	return info, nil
}
