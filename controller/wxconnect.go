package controllers

import (
	"crypto/sha1"
	"fmt"
	"github.com/astaxie/beego"
	"io"
	"sort"
	"strings"
)

type WxConnectController struct {
	beego.Controller
}

const Token = "solomagic695981642"

func (w *WxConnectController) Get() {
	timestamp, nonce, signatureIn := w.GetString("timestamp"), w.GetString("nonce"), w.GetString("signature")
	signatureGen := makeSignature(timestamp, nonce)
	if signatureGen != signatureIn {
		fmt.Printf("signatureGen != signatureIn signatureGen=%s,signatureIn=%s\n", signatureGen, signatureIn)
	} else {
		echostr := w.GetString("echostr")
		w.Ctx.WriteString(echostr)
	}
}

func makeSignature(timestamp, nonce string) string {
	sl := []string{Token, timestamp, nonce}
	sort.Strings(sl)
	s := sha1.New()
	io.WriteString(s, strings.Join(sl, ""))
	return fmt.Sprintf("%x", s.Sum(nil))
}
