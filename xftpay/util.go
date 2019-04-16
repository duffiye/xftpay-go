package xftpay

import (
	"bytes"
	"crypto"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"log"
	"reflect"
	"sort"
	"strings"
)

// 简单封装了json的Marshal功能
// pingpp.JsonEncode(param1)
func JsonEncode(v interface{}) ([]byte, error) {
	return json.Marshal(&v)
}

// 简单封装了json的UnMarshal功能
// Example pingpp.JsonDecode(param1, param2)
// param1：需要转换成结构体的json数据
// param2：转换后数据容器
func JsonDecode(p []byte, v interface{}) error {
	obj := json.NewDecoder(bytes.NewBuffer(p))
	obj.UseNumber()
	return obj.Decode(&v)
}

//用商户的私钥去生成签名目前在创建订单的时候使用
func GenSign(data []byte, privateKey []byte) (sign []byte, err error) {
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, errors.New("private key error")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	hashFunc := crypto.SHA256
	h := hashFunc.New()
	h.Write(data)
	hashed := h.Sum(nil)
	return rsa.SignPKCS1v15(rand.Reader, priv, hashFunc, hashed)
}

//用ping++公钥去验证签名目前在Webhook时候使用
func Verify(data []byte, publicKey []byte, sign []byte) (err error) {
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return errors.New("public key error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return err
	}
	pub := pubInterface.(*rsa.PublicKey)

	hashFunc := crypto.SHA256
	h := hashFunc.New()
	h.Write(data)
	hashed := h.Sum(nil)
	return rsa.VerifyPKCS1v15(pub, hashFunc, hashed, sign)
}

func GenMD5Sign(data interface{}, key string) (d interface{}, err error) {
	paramMap := Struct2Map(data)

	fmt.Println(paramMap)

	var keys []string
	for k := range paramMap {
		if k == "sign" {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)

	dataToBeSgin := ""

	for _, v := range keys {
		if paramMap[v] == "" || paramMap[v] == "null" {
			continue
		}
		fmt.Printf("v is %q\n", v)
		dataToBeSgin += fmt.Sprintf("%s=%s&", v, paramMap[v])
	}

	dataToBeSgin += fmt.Sprintf("key=%s", key)

	v := reflect.ValueOf(data).Elem()
	md5Str := MD5(dataToBeSgin)
	signField := v.FieldByName("Sign")

	log.Printf("Sign Field %s\n", signField)

	log.Printf("Can Set %v\n", signField.CanSet())

	signField.SetString(md5Str)

	return data, nil

}

func Struct2Map(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj).Elem()
	v := reflect.ValueOf(obj).Elem()

	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		tag := t.Field(i).Tag.Get("json")
		fmt.Printf("Tag Name is %s\n", tag)
		//t.Field(i).Name
		data[tag] = v.Field(i).Interface()

	}
	return data
}

func MD5(str string) string {
	log.Printf("待签名字符串 :%s\n", str)
	data := []byte(str)
	has := md5.Sum(data)
	md5str1 := fmt.Sprintf("%x", has) //将[]byte转成16进制
	log.Printf("md5 %s\n", md5str1)
	return strings.ToUpper(md5str1)
}
