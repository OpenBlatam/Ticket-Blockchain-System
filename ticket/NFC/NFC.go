import (
    "crypto/hmac"
    "crypto/rand"
    "crypto/sha256"
    "encoding/hex"
    "encoding/json"
    "io"
)

func CreateNFCTagSignature(ticket *Ticket, secretKey string) (string, error) {
    jsonData, err := json.Marshal(ticket)
    if err != nil {
        return "", err
    }

    data := []byte(secretKey + string(jsonData))

    // Reuse hash for HMAC-SHA256 operation
    hash := sha256.Sum256(data)

    salt := make([]byte, 32)
    if _, err := io.ReadFull(rand.Reader, salt); err != nil {
        return "", err
    }

    info := []byte("nfctaginfo")

    // Use HKDF to derive a key from the secret key and hash value
    key := make([]byte, 32)
    hkdf := hkdf.New(sha256.New, []byte(secretKey), salt, info)
    if _, err := io.ReadFull(hkdf, key); err != nil {
        return "", err
    }

    // Use the derived key to create a HMAC-SHA256 signature of the hash value
    mac := hmac.New(sha256.New, key)
    mac.Write(hash[:])
    signature := mac.Sum(nil)

    return hex.EncodeToString(signature), nil
}
