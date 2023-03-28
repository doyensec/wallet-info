package utils

import (
	"errors"
	"fmt"
	neturl "net/url"
	"strings"
)

func GetHost(url string) string {
	if url == "" {
		return ""
	}

	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = fmt.Sprintf("https://%s", url)
	}

	parsed, err := neturl.Parse(url)
	if err != nil {
		return ""
	}

	// can return empty string
	return parsed.Hostname()
}

func GetTldPlusOne(host string) (string, error) {
	parts := strings.Split(host, ".")
	len := len(parts)

	if len < 2 {
		return "", errors.New("cannot parse host")
	}

	if len == 2 {
		return host, nil
	}

	return fmt.Sprintf("%s.%s", parts[len-2], parts[len-1]), nil
}
