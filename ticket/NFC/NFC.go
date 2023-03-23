import (
"crypto/sha256"
"encoding/hex"
"encoding/json"
"golang.org/x/crypto/hkdf"
)

func CreateNFCTagSignature(ticket *Ticket, secretKey string) string {
jsonData, _ := json.Marshal(ticket)
data := []byte(secretKey + string(jsonData))
hash := sha256.Sum256(data)

salt := []byte("nfctagsalt") // A unique salt value
info := []byte("nfctaginfo") // A unique info value

// Use HKDF to derive a key from the secret key and hash value
key := make([]byte, 32)
hkdf.New(sha256.New, []byte(secretKey), salt, info).Read(key)

// Use the derived key to create a HMAC-SHA256 signature of the hash value
mac := hmac.New(sha256.New, key)
mac.Write(hash[:])
signature := mac.Sum(nil)

return hex.EncodeToString(signature)
}
