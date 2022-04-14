package identityclient

// Response to GET is a simple string now

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
