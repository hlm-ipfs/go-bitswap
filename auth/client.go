package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	logging "github.com/ipfs/go-log"
	"io/ioutil"
	"net/http"
)

var log = logging.Logger("auth")

//检索文件(input: cid)
//rpc Create (basic.String) returns (RetrievalCreateResponse)
//验证检索token
//rpc Verify (basic.String) returns (basic.Empty)
func CreateRetrievalOrder(cid string, headers http.Header, order *ResponseCreateRetrievalOrder) error {
	client := &http.Client{}
	reqBody := RequestCreateRetrievalOrder{
		Value: cid,
	}
	reqBodyBytes, _ := json.Marshal(reqBody)
	req, err := http.NewRequest("POST", "http://103.44.247.16:31686/market/Retrieval/Create", bytes.NewBuffer(reqBodyBytes))
	if err != nil {
		log.Error(err)
		return err
	}
	req.Header = headers
	log.Infof("headers: %+v", req.Header)
	resp, err := client.Do(req)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return err
	}
	type Response struct {
		Status  int32  `json:"status"`
		Code    string `json:"code"`
		Message string `json:"message"`
		Value   string `json:"value"`
	}
	respBody := Response{}
	err = json.Unmarshal(body, &respBody)
	if err != nil {
		log.Error(err)
		return err
	}
	log.Infof("response: %+v", respBody)
	if respBody.Status >= 400 {
		log.Warnf("bad response %+v", respBody)
	}
	err = json.Unmarshal(bytes.NewBufferString(respBody.Value).Bytes(), order)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

type RequestCreateRetrievalOrder struct {
	Value string //cid
}

type ResponseCreateRetrievalOrder struct {
	OrderNo string
	Token   string
}

func VerifyRetrievalToken(token string) error {
	client := &http.Client{}
	reqBody := RequestVerifyRetrievalToken{
		Value: token,
	}
	reqBodyBytes, _ := json.Marshal(reqBody)
	req, err := http.NewRequest("POST", "http://103.44.247.16:31686/market/Retrieval/Verify", bytes.NewBuffer(reqBodyBytes))
	if err != nil {
		log.Error(err)
		return err
	}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return err
	}
	type Response struct {
		Status  int32  `json:"status"`
		Code    string `json:"code"`
		Message string `json:"message"`
		Value   string `json:"value"`
	}
	respBody := Response{}
	err = json.Unmarshal(body, &respBody)
	if err != nil {
		log.Error(err)
		return err
	}
	if respBody.Status != 200 {
		log.Warnf("bad response %+v", respBody)
		return errors.New(respBody.Message)
	}
	return nil
}

type RequestVerifyRetrievalToken struct {
	Value string //token
}
