package client


func toBytes(msg string) []byte {
	return []byte(msg)
}

func toString (msgBytes []byte) string {
	return string(msgBytes)
}