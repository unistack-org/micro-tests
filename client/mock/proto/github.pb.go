package pb

import (
	_ "go.unistack.org/micro-proto/v3/api"
	_ "go.unistack.org/micro-proto/v3/openapiv3"
)

type LookupUserReq struct {
	Username string `json:"username,omitempty"`
}

type LookupUserRsp struct {
	Name string `json:"name,omitempty"`
}

type Error struct {
	Message          string `json:"message,omitempty"`
	DocumentationUrl string `json:"documentation_url,omitempty"`
}
