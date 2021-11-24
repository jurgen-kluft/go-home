package config

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"database/sql/driver"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
)

type CryptString struct {
	String string
}

var cryptKeeperKey []byte

// MarshalJSON encrypts and marshals nested String
func (cs *CryptString) MarshalJSON() ([]byte, error) {
	encString, err := Encrypt(cs.String)
	if err != nil {
		return nil, err
	}
	return json.Marshal(encString)
}

// UnmarshalJSON encrypts and marshals nested String
func (cs *CryptString) UnmarshalJSON(b []byte) error {
	var decString string
	err := json.Unmarshal(b, &decString)
	//fmt.Println("Unmarshal CryptString", decString)
	if err != nil {
		return err
	}
	cs.String, err = Decrypt(decString)
	return err
}

// Scan implements sql.Scanner and decryptes incoming sql column data
func (cs *CryptString) Scan(value interface{}) error {
	switch v := value.(type) {
	case string:
		rawString, err := Decrypt(v)
		if err != nil {
			return err
		}
		cs.String = rawString
	case []byte:
		rawString, err := Decrypt(string(v))
		if err != nil {
			return err
		}
		cs.String = rawString
	default:
		return fmt.Errorf("couldn't scan %+v", reflect.TypeOf(value))
	}
	return nil
}

// Value implements driver.Valuer and encrypts outgoing bind values
func (cs CryptString) Value() (value driver.Value, err error) {
	return Encrypt(cs.String)
}

// SetCryptKey with user input
func SetCryptKey(secretKey []byte) error {
	keyLen := len(secretKey)
	if keyLen != 16 && keyLen != 24 && keyLen != 32 {
		return fmt.Errorf("Invalid KEY to set for GO_HOME_KEY; must be 16, 24, or 32 bytes (got %d)", keyLen)
	}
	cryptKeeperKey = secretKey
	return nil
}

// CryptKey returns a valid Crypt key
func CryptKey() []byte {
	if cryptKeeperKey == nil {
		key := os.Getenv("GO_HOME_KEY")
		if key == "" {
			fmt.Println("Error, you did not set the environment variable GO_HOME_KEY")
		} else {
			SetCryptKey([]byte(key))
		}
	}
	//fmt.Println("CryptKey:", string(cryptKeeperKey))
	return cryptKeeperKey
}

// Encrypt AES-encrypt string and then base64-encode
func Encrypt(text string) (string, error) {
	plaintext := []byte(text)

	block, err := aes.NewCipher(CryptKey())
	if err != nil {
		return "", err
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	cipher.NewCFBEncrypter(block, iv).XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	// convert to base64
	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

// Decrypt base64-decode and then AES decrypt string
func Decrypt(cryptoText string) (string, error) {
	ciphertext, err := base64.URLEncoding.DecodeString(cryptoText)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(CryptKey())
	if err != nil {
		return "", err
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	if byteLen := len(ciphertext); byteLen < aes.BlockSize {
		return "", fmt.Errorf("invalid cipher size %d", byteLen)
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	// XORKeyStream can work in-place if the two arguments are the same.
	cipher.NewCFBDecrypter(block, iv).XORKeyStream(ciphertext, ciphertext)

	return string(ciphertext), nil
}
