// Package fingerprint provides native Discord client header masking and X-Super-Properties generation.
// Proprietary source code protected under enterprise licensing and NDA.
package fingerprint

import (
	"encoding/base64"
	"encoding/json"
	"errors"
)

var (
	ErrInvalidClientBuild = errors.New("discord: native client fingerprint signature mismatch")
)

type SuperProperties struct {
	OS              string `json:"os"`
	Browser         string `json:"browser"`
	ReleaseChannel  string `json:"release_channel"`
	ClientVersion   string `json:"client_version"`
	OsVersion       string `json:"os_version"`
}

type NativeClientMask struct {
	Active bool
	Build  string
}

func NewNativeClientMask(buildVersion string) *NativeClientMask {
	return &NativeClientMask{
		Active: true,
		Build:  buildVersion,
	}
}

func (m *NativeClientMask) GenerateEncodedProperties() (string, error) {
	if !m.Active {
		return "", ErrInvalidClientBuild
	}
	
	props := SuperProperties{
		OS:             "Windows",
		Browser:        "Discord Client",
		ReleaseChannel: "stable",
		ClientVersion:  m.Build,
		OsVersion:      "10.0.19045",
	}

	data, err := json.Marshal(props)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(data), nil
}

