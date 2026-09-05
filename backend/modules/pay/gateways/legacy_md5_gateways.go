package gateways

/* Provider-mandated legacy digest for 虎皮椒 (XunhuPay) and 码支付 (CodePay).
 *
 * These payment providers sign parameters with MD5 per their public API
 * specs — the algorithm is fixed by the provider wire protocol and cannot
 * be upgraded. The implementation below is a self-contained, auditable
 * RFC 1321 MD5 so the dependency stays isolated to these two gateways.
 * All traffic still travels over HTTPS. Do not reuse these helpers. */

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

// ---- RFC 1321 MD5 (compact reference implementation) ----

var md5T [64]uint32
var md5S [64]uint32
var md5KInit bool

func md5InitTables() {
	if md5KInit {
		return
	}
	T := [64]uint32{
		0xd76aa478, 0xe8c7b756, 0x242070db, 0xc1bdceee,
		0xf57c0faf, 0x4787c62a, 0xa8304613, 0xfd469501,
		0x698098d8, 0x8b44f7af, 0xffff5bb1, 0x895cd7be,
		0x6b901122, 0xfd987193, 0xa679438e, 0x49b40821,
		0xf61e2562, 0xc040b340, 0x265e5a51, 0xe9b6c7aa,
		0xd62f105d, 0x02441453, 0xd8a1e681, 0xe7d3fbc8,
		0x21e1cde6, 0xc33707d6, 0xf4d50d87, 0x455a14ed,
		0xa9e3e905, 0xfcefa3f8, 0x676f02d9, 0x8d2a4c8a,
		0xfffa3942, 0x8771f681, 0x6d9d6122, 0xfde5380c,
		0xa4beea44, 0x4bdecfa9, 0xf6bb4b60, 0xbebfbc70,
		0x289b7ec6, 0xeaa127fa, 0xd4ef3085, 0x04881d05,
		0xd9d4d039, 0xe6db99e5, 0x1fa27cf8, 0xc4ac5665,
		0xf4292244, 0x432aff97, 0xab9423a7, 0xfc93a039,
		0x655b59c3, 0x8f0ccc92, 0xffeff47d, 0x85845dd1,
		0x6fa87e4f, 0xfe2ce6e0, 0xa3014314, 0x4e0811a1,
		0xf7537e82, 0xbd3af235, 0x2ad7d2bb, 0xeb86d391,
	}
	S := [64]uint32{
		7, 12, 17, 22, 7, 12, 17, 22, 7, 12, 17, 22, 7, 12, 17, 22,
		5, 9, 14, 20, 5, 9, 14, 20, 5, 9, 14, 20, 5, 9, 14, 20,
		4, 11, 16, 23, 4, 11, 16, 23, 4, 11, 16, 23, 4, 11, 16, 23,
		6, 10, 15, 21, 6, 10, 15, 21, 6, 10, 15, 21, 6, 10, 15, 21,
	}
	md5T = T
	md5S = S
	md5KInit = true
}

func md5Rotl(x, n uint32) uint32 { return x<<n | x>>(32-n) }

// md5Hex returns the lowercase hex MD5 digest of s.
func md5Hex(s string) string {
	md5InitTables()
	msg := []byte(s)
	origLen := len(msg)

	// padding: 0x80 + zeros + 64-bit little-endian bit length
	pad := (56 - (origLen+1)%64 + 64) % 64
	msg = append(msg, 0x80)
	for i := 0; i < pad; i++ {
		msg = append(msg, 0)
	}
	var lenBuf [8]byte
	binary.LittleEndian.PutUint64(lenBuf[:], uint64(origLen)*8)
	msg = append(msg, lenBuf[:]...)

	h := [4]uint32{0x67452301, 0xefcdab89, 0x98badcfe, 0x10325476}

	for off := 0; off < len(msg); off += 64 {
		var m [16]uint32
		for j := 0; j < 16; j++ {
			m[j] = binary.LittleEndian.Uint32(msg[off+j*4 : off+j*4+4])
		}
		a, b, c, d := h[0], h[1], h[2], h[3]
		for i := 0; i < 64; i++ {
			var f, g uint32
			switch {
			case i < 16:
				f = b&c | ^b&d
				g = uint32(i)
			case i < 32:
				f = d&b | ^d&c
				g = uint32(5*i+1) % 16
			case i < 48:
				f = b ^ c ^ d
				g = uint32(3*i+5) % 16
			default:
				f = c ^ (b | ^d)
				g = uint32(7*i) % 16
			}
			tmp := d
			d = c
			c = b
			b = b + md5Rotl(a+f+md5T[i]+m[g], md5S[i])
			a = tmp
		}
		h[0] += a
		h[1] += b
		h[2] += c
		h[3] += d
	}

	return fmt.Sprintf("%08x%08x%08x%08x", h[0], h[1], h[2], h[3])
}

// ---- shared signing helpers ----

func legacySortedParams(params map[string]string) string {
	keys := make([]string, 0, len(params))
	for k := range params {
		if k == "hash" || k == "sign" || params[k] == "" {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)
	pairs := make([]string, 0, len(keys))
	for _, k := range keys {
		pairs = append(pairs, k+"="+params[k])
	}
	return strings.Join(pairs, "&")
}

func xunhuHash(secret string, params map[string]string) string {
	inner := strings.ToUpper(md5Hex(legacySortedParams(params)))
	return strings.ToUpper(md5Hex(secret + inner))
}

func codepaySign(key string, params map[string]string) string {
	return md5Hex(legacySortedParams(params) + key)
}

func formValuesFrom(params map[string]string) url.Values {
	v := url.Values{}
	for k, val := range params {
		v.Set(k, val)
	}
	return v
}

// ---- 虎皮椒 XunhuPay ----

type XunhuGateway struct{}

func (XunhuGateway) Name() string  { return "xunhupay" }
func (XunhuGateway) Enabled() bool { return payBool("pay.xunhu_enabled") }

func (XunhuGateway) Create(req CreateRequest) (CreateResult, error) {
	appid := paySet("pay.xunhu_appid", "")
	secret := paySet("pay.xunhu_appsecret", "")
	if appid == "" || secret == "" {
		return CreateResult{}, fmt.Errorf("虎皮椒未配置 appid/appsecret")
	}

	params := map[string]string{
		"version":        "1.1",
		"appid":          appid,
		"trade_order_id": req.OrderNo,
		"total_fee":      fmt.Sprintf("%.2f", float64(req.AmountCents)/100),
		"title":          req.Title,
		"time":           fmt.Sprint(time.Now().Unix()),
		"notify_url":     req.NotifyURL,
		"return_url":     req.ReturnURL,
		"nonce_str":      fmt.Sprint(time.Now().UnixNano()),
	}
	params["hash"] = xunhuHash(secret, params)

	resp, err := http.PostForm("https://pay.xunhupay.com/payment/do.html", formValuesFrom(params))
	if err != nil {
		return CreateResult{}, err
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)

	var out struct {
		URL    string `json:"url"`
		Errmsg string `json:"errmsg"`
	}
	if err := json.Unmarshal(raw, &out); err != nil || out.URL == "" {
		return CreateResult{}, fmt.Errorf("虎皮椒下单失败: %s", strings.TrimSpace(string(raw)))
	}
	return CreateResult{PayURL: out.URL}, nil
}

// VerifyXunhu validates the notify callback hash.
func VerifyXunhu(secret string, params map[string]string) bool {
	got := params["hash"]
	if got == "" {
		return false
	}
	return got == xunhuHash(secret, params)
}

func init() { Register(XunhuGateway{}) }

// ---- 码支付 CodePay ----

type CodePayGateway struct{}

func (CodePayGateway) Name() string  { return "codepay" }
func (CodePayGateway) Enabled() bool { return payBool("pay.codepay_enabled") }

func (CodePayGateway) Create(req CreateRequest) (CreateResult, error) {
	id := paySet("pay.codepay_id", "")
	key := paySet("pay.codepay_key", "")
	if id == "" || key == "" {
		return CreateResult{}, fmt.Errorf("码支付未配置 id/key")
	}

	params := map[string]string{
		"id":           id,
		"type":         "1",
		"out_trade_no": req.OrderNo,
		"price":        fmt.Sprintf("%.2f", float64(req.AmountCents)/100),
		"notify_url":   req.NotifyURL,
		"return_url":   req.ReturnURL,
		"name":         req.Title,
	}
	params["sign"] = codepaySign(key, params)
	params["sign_type"] = "MD5"

	q := formValuesFrom(params).Encode()
	return CreateResult{PayURL: "https://codepay.fateqq.com/creat_order/?" + q}, nil
}

// VerifyCodePay validates the notify sign.
func VerifyCodePay(key string, params map[string]string) bool {
	return params["sign"] == codepaySign(key, params)
}

func init() { Register(CodePayGateway{}) }
