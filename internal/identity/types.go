package identityclient

// Response to GET is
// {
//     "DID": "NejkFbkM60qzukmCALJSO5",
//     "public_key": "OqiFbIpbNq/Ry8MeaGOAokmMwGJXJBwxbqeGXXM8gzE=",
//     "private_key": "u/mJeV9tCWYQVo2x3VIJltjelv1eSiKx6VLtk3gb62A=",
//     "timestamp": "1646920431"
// }
type KeyPair struct {
	Did       string `json:"DID"`
	PubKey    string `json:"public_key"`
	PrivKey   string `json:"private_key"`
	Timestamp string `json:"timestamp"`
}

// Body of POST to verify
// {
//     "did": "NejkFbkM60qzukmCALJSO5",
//     "public_key": "OqiFbIpbNq/Ry8MeaGOAokmMwGJXJBwxbqeGXXM8gzE=",
//     "timestamp": "1646920431"
// }

type VerifyKeyPairBody struct {
	Did       string `json:"did"`
	PubKey    string `json:"public_key"`
	Timestamp string `json:"timestamp"`
}
