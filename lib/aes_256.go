package lib

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

var (
	initialVectorPhrase = os.Getenv("aesiv")
	initialVector       = createHash(initialVectorPhrase)[0:16]
	keyphrase           = os.Getenv("aeskey")
	key                 = []byte(createHash(keyphrase))
)

func createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

func AESEncrypt(src string) string {
	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println("key error1", err)
	}
	if src == "" {
		fmt.Println("plain content empty")
	}
	ecb := cipher.NewCBCEncrypter(block, []byte(initialVector))
	content := []byte(src)
	content = PKCS5Padding(content, block.BlockSize())
	crypted := make([]byte, len(content))
	ecb.CryptBlocks(crypted, content)
	return base64.StdEncoding.EncodeToString(crypted)
}
func AESDecrypt(crypt string) []byte {
	data, _ := base64.StdEncoding.DecodeString(crypt)
	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println("key error1", err)
		return []byte("")
	}
	if len(data) == 0 {
		fmt.Println("plain content empty")
		return []byte("")
	}
	ecb := cipher.NewCBCDecrypter(block, []byte(initialVector))
	if len(data)%aes.BlockSize != 0 {
		return []byte("Invalid input length!")
	}
	decrypted := make([]byte, len(data))
	ecb.CryptBlocks(decrypted, []byte(data))
	return PKCS5Trimming(decrypted)
}
func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}
func PKCS5Trimming(encrypt []byte) []byte {
	padding := encrypt[len(encrypt)-1]
	return encrypt[:len(encrypt)-int(padding)]
}

// func main() {
// var plainText = "Lick it till Ice Cream I just dont know what to do!"
// fmt.Println("KEY: ", string(key))
// fmt.Println("IV:", initialVector)
// encryptedData := AESEncrypt(plainText)
// fmt.Println(encryptedData)
// decryptedText := AESDecrypt(encryptedData)
// fmt.Println(string(decryptedText))
// }

func isJSON(s string) bool {
	var js map[string]interface{}
	return json.Unmarshal([]byte(s), &js) == nil
}

type MotionJSON struct {
	Content string `json:"Content,omitempty" binding:"required"`
	Error   string `json:"Error"`
}

func AESDecryptRequestJSON(c *gin.Context) ([]byte, string) {
	var bodyBytes []byte
	var MotionReq MotionJSON
	if err := c.ShouldBindBodyWith(&MotionReq, binding.JSON); err != nil {
		MotionRes := MotionJSON{
			Error: err.Error(),
		}

		jsonRes, _ := json.Marshal(MotionRes)
		encryptRes := AESEncrypt(string(jsonRes))
		c.JSON(400, gin.H{"Content": encryptRes})
		return []byte(""), err.Error()
	}

	if c.Copy().Request.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(c.Copy().Request.Body)
	}

	// Get Request Data
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	reqJSON := string(bodyBytes)
	json.Unmarshal([]byte(reqJSON), &MotionReq)

	decryptReq := AESDecrypt(MotionReq.Content)
	checkJSON := isJSON(string(decryptReq))
	switch checkJSON {
	case false:
		MotionRes := MotionJSON{
			Error: "Invalid JSON",
		}

		jsonRes, _ := json.Marshal(MotionRes)
		encryptRes := AESEncrypt(string(jsonRes))
		c.JSON(400, gin.H{"Content": encryptRes})
		return []byte(""), "Error"
	default:
	}
	return decryptReq, ""
}
