package monitor

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"

	"monitor/internal/svc"
	"monitor/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type EncryptLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewEncryptLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EncryptLogic {
	return &EncryptLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *EncryptLogic) Encrypt(req *types.EncryptRes) (resp *types.EncryptResp, err error) {
	resp = &types.EncryptResp{
		Code: "10000",
		Msg:  "success",
		Data: &types.EncryptRp{
			Encrypted: make(map[string]string),
			OriString: make(map[string]string),
		},
	}
	for _, oristr := range req.OriString {
		resp.Data.Encrypted[oristr], _ = Encrypt(oristr, svc.Keystr)
	}
	for _, encryped := range req.Encrypted {
		oristr, _ := Decrypt(encryped, svc.Keystr)
		resp.Data.OriString[oristr] = encryped
	}
	return resp, nil
}
func Encrypt(msg string, keystr string) (string, error) {
	key := []byte(keystr)
	src := []byte(msg)
	//生成cipher.Block 数据块
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	} else if len(src) == 0 {
		return "", errors.New("src is empty")
	}

	//填充内容，如果不足16位字符
	blockSize := block.BlockSize()
	originData := pad(src, blockSize)

	//加密，输出到[]byte数组
	crypted := make([]byte, aes.BlockSize+len(originData))
	iv := crypted[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", nil
	}
	//加密方式
	blockMode := cipher.NewCBCEncrypter(block, iv)
	blockMode.CryptBlocks(crypted[aes.BlockSize:], originData)
	return base64.StdEncoding.EncodeToString(crypted), nil
}

func pad(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func Decrypt(src, keystr string) (string, error) {
	defer func() {
		if e := recover(); e != nil {
			return
		}
	}()
	key := []byte(keystr)
	decode_data, err := base64.StdEncoding.DecodeString(src)
	if err != nil {
		return "", nil
	}
	//生成密码数据块cipher.Block
	block, _ := aes.NewCipher(key)
	//解密模式
	blockMode := cipher.NewCBCDecrypter(block, decode_data[:aes.BlockSize])
	//输出到[]byte数组
	origin_data := make([]byte, len(decode_data)-aes.BlockSize)
	blockMode.CryptBlocks(origin_data, decode_data[aes.BlockSize:])
	//去除填充,并返回
	return string(unpad(origin_data)), nil
}

func unpad(ciphertext []byte) []byte {
	length := len(ciphertext)
	//去掉最后一次的padding
	unpadding := int(ciphertext[length-1])
	return ciphertext[:(length - unpadding)]
}
