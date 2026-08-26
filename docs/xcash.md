API 对接

网关地址与链码
所有 Django/DRF API 路由均不带尾部 /，示例中的路径请按原样请求。

自部署
API 网关地址即为 .env 中配置的 SITE_DOMAIN，例如 https://{你的域名}/v1/invoice；管理后台同域名访问。

Xcash 官方服务
用途	URL
API 网关	https://pay.xca.sh（如 /v1/invoice）
EPay 网关	https://pay.xca.sh/epay/submit.php
SaaS 控制台	https://dash.xca.sh
公链 - chain
名称	API 参数	类型	Gas 代币	Chain ID
Ethereum	ethereum	EVM	ETH	1
BNB Smart Chain	bsc	EVM	BNB	56
Polygon PoS	polygon	EVM	POL	137
Arbitrum One	arbitrum-one	EVM	ETH	42161
Optimism	optimism	EVM	ETH	10
Base	base	EVM	ETH	8453
Tron	tron	Tron	TRX	-
Sepolia	sepolia	EVM 测试网	ETH	11155111
Nile	nile	Tron 测试网	TRX	-
Anvil Local	anvil	EVM 本地链	ETH	31337
代币 - crypto
crypto 使用币种 symbol，一律大写。常见 symbol 如下：

名称	API 参数	类型	说明
Tether USD	USDT	稳定币	ERC-20 / TRC-20，最常用收款币种
USD Coin	USDC	稳定币	ERC-20，多条 EVM 链可用
Dai	DAI	稳定币	ERC-20
Ether	ETH	原生币	Ethereum 及 Arbitrum / Optimism / Base 的原生 Gas 币
BNB	BNB	原生币	BNB Smart Chain 的原生 Gas 币
Polygon	POL	原生币	Polygon PoS 的原生 Gas 币
TRON	TRX	原生币	Tron 的原生 Gas 币
上表仅为常见示例：EVM 链支持任意 ERC-20 代币，在后台添加代币合约地址并启用即可收款，可用 symbol 不止于此；Tron 当前仅放行 USDT 与原生 TRX。实际可用的链币组合取决于后台启用的链、币种与链上部署关系。测试项目只能使用测试网或本地链，非测试项目只能使用主网链。

认证与签名
除明确标注为公开的端点外，/v1/* 接口都需要 HMAC-SHA256 签名。在管理后台创建项目后，系统生成 appid（如 XC-A3BK7NMG）与 hmac_key。

请求头
http

复制
XC-Appid: {appid}
XC-Timestamp: {unix_timestamp}
XC-Nonce: {unique_nonce}
XC-Signature: {hmac_signature}
Content-Type: application/json
Header	说明
XC-Appid	项目 AppID
XC-Timestamp	当前 Unix 时间戳，生产环境允许与服务器相差 300 秒
XC-Nonce	同一 AppID 下 300 秒内不可重复
XC-Signature	HMAC-SHA256 签名，小写十六进制
签名计算
text

复制
message   = XC-Nonce + XC-Timestamp + request_body
signature = HMAC-SHA256(message, hmac_key).hexdigest()
request_body 必须是实际发送的原始请求体字符串。GET 请求没有 body 时使用空字符串 ""。

Python 示例
python

复制
import hashlib
import hmac
import json
import time
import uuid

appid = "XC-A3BK7NMG"
hmac_key = "your_hmac_key"

payload = {
    "out_no": "order-001",
    "title": "Premium Plan",
    "currency": "USD",
    "amount": "29.99",
}
body = json.dumps(payload, separators=(",", ":"), ensure_ascii=False)
timestamp = str(int(time.time()))
nonce = str(uuid.uuid4())
signature = hmac.new(
    hmac_key.encode(),
    f"{nonce}{timestamp}{body}".encode(),
    hashlib.sha256,
).hexdigest()

headers = {
    "XC-Appid": appid,
    "XC-Timestamp": timestamp,
    "XC-Nonce": nonce,
    "XC-Signature": signature,
    "Content-Type": "application/json",
}
Node.js 示例
javascript

复制
const crypto = require("crypto");

const appid = "XC-A3BK7NMG";
const hmacKey = "your_hmac_key";
const body = JSON.stringify({
  out_no: "order-001",
  title: "Premium Plan",
  currency: "USD",
  amount: "29.99",
});
const timestamp = Math.floor(Date.now() / 1000).toString();
const nonce = crypto.randomUUID();
const signature = crypto
  .createHmac("sha256", hmacKey)
  .update(`${nonce}${timestamp}${body}`)
  .digest("hex");
响应与错误码
成功时直接返回业务 JSON。创建类接口通常返回 HTTP 201，查询类接口返回 HTTP 200。业务错误响应：

json

复制
{
  "code": "1001",
  "message": "AppID无效",
  "detail": ""
}
框架级错误（如资源不存在 404、方法错误 405、限流 429）可能返回 DRF 默认格式 { "detail": "Not found." }。

接口列表
方法	路径	说明	签名
POST	/v1/invoice	创建账单收款	需要
GET	/v1/invoice/{sys_no}	查询账单收款公开状态	不需要
GET	/v1/deposit/address	获取充值收款地址	需要
GET/POST	/epay/submit.php	易支付 V1 创建订单	EPay MD5
错误码
错误码	说明	HTTP
1000	参数错误	400
1001	AppID 无效	400
1002	IP 禁止	403
1003	签名错误	403
1004	项目未配置	400
1007	单号 out_no 重复	400
1008	Timestamp 未设置或过期	400
1009	请求重复	400
2000 / 2001 / 2002	无效链 / 无效加密货币 / 本链不支持此加密货币	400
4000 / 4001 / 4002	无效 UID / 项目未配置该链归集地址 / 充值用户数达上限	400 / 400 / 403
5000 / 5008 / 5009	账单收款类型错误 / 无可用账单收款方式 / 待支付记录过多	400
6000 / 6002 / 6004	内部令牌无效 / 项目不存在 / 账户已冻结	401 / 404 / 403
创建账单收款
POST /v1/invoice（需要签名）。创建一笔账单收款，成功后返回 pay_url，买家打开账单收款页完成选币、选链和付款。

请求参数
字段	类型	必填	说明
out_no	string	是	商户订单号，最长 32 位，同一项目内唯一
title	string	是	账单标题，最长 32 位
currency	string	是	计价法币代码（如 USD、CNY），收款加密货币由 methods 指定
amount	string	是	计价金额，范围 0.00000001 – 1000000
duration	integer	否	有效期分钟数，5 – 30，默认 10
methods	object	否	限定收款方式，格式 {"币种": ["链码"]}
notify_url	string	否	账单级 Webhook 地址，优先于项目默认通知地址
return_url	string	否	账单完成后的同步跳转地址
methods 规则
不传 methods：系统按项目配置生成当前可用的链币组合。
传入 methods：必须是系统生成组合的子集，否则返回无可用收款方式。
currency 只决定 amount 的计价单位（法币），与买家实际支付的加密货币解耦——后者由 methods 限定。
买家最终应支付的链、币、地址与数量以账单页/查询接口返回的 chain、crypto、pay_address、pay_amount 为准，不要自行按 amount 推导链上付款数量。
请求示例
json

复制
{
  "out_no": "order-20260602-001",
  "title": "Premium Plan",
  "currency": "USD",
  "amount": "29.99",
  "duration": 15,
  "methods": {
    "USDT": ["ethereum"],
    "USDC": ["base"]
  },
  "notify_url": "https://merchant.example.com/xcash/webhook",
  "return_url": "https://merchant.example.com/payment/success"
}
限定只允许买家使用某种稳定币（USD 计价 + methods 限定 USDT）：

json

复制
{
  "out_no": "order-20260602-002",
  "title": "Contract Invoice",
  "currency": "USD",
  "amount": "100",
  "duration": 15,
  "methods": {
    "USDT": ["ethereum", "base"]
  },
  "notify_url": "https://merchant.example.com/xcash/webhook"
}
响应示例
json

复制
{
  "appid": "XC-A3BK7NMG",
  "sys_no": "INV2606028X7K2P9Q",
  "out_no": "order-20260602-001",
  "title": "Premium Plan",
  "currency": "USD",
  "amount": "29.99",
  "methods": { "USDT": ["ethereum"], "USDC": ["base"] },
  "chain": null,
  "crypto": null,
  "crypto_address": null,
  "pay_address": null,
  "pay_amount": null,
  "pay_url": "https://pay.xca.sh/pay/INV2606028X7K2P9Q",
  "started_at": "2026-06-02T12:00:00Z",
  "created_at": "2026-06-02T12:00:00Z",
  "expires_at": "2026-06-02T12:15:00Z",
  "notify_url": "https://merchant.example.com/xcash/webhook",
  "return_url": "https://merchant.example.com/payment/success",
  "payment": null,
  "status": "waiting",
  "risk_level": null,
  "risk_score": null
}
如果最终只剩一个收款组合，系统会在创建时自动选择，此时 chain、crypto、pay_address、pay_amount 可能已返回具体值。默认匿名限流 256/分钟。

查询账单收款
GET /v1/invoice/{sys_no}（公开接口，无需签名）。用于账单收款页或买家侧轮询状态，不返回 appid、out_no、notify_url。

关键响应字段
字段	说明
sys_no	系统单号；前缀 + 6 位日期(YYMMDD) + 8 位大写字母数字。账单前缀 INV，充值前缀 DXC
chain	已选链，未选时为 null
crypto	已选币种，未选时为 null
pay_address	收款地址
pay_amount	买家应付加密货币数量
payment_uri	EVM 链可用的 EIP-681 支付 URI；非 EVM 或无法精确编码金额时为空
status	waiting / completed / expired
payment	匹配到的链上转账对象，含 hash、block、from/to、amount、confirm_progress 等
risk_level	风险等级
risk_score	风险分数
限流 60/分钟，按 sys_no + IP 维度。

获取充值收款地址
GET /v1/deposit/address（需要签名）。为项目下的终端客户获取充值收款地址，同一项目、同一 uid、同一链稳定返回同一地址。

虽然从实现上看，同一项目、单个 uid 的充值收款地址理论上可以在所有链之间共享；但为了语义清晰，建议按每个 chain + crypto 组合重新调用本接口获取地址，不要只把某个链或币种下取到的地址直接复用于其他组合。

字段	类型	必填	说明
uid	string	是	终端客户标识，1 – 128 位，仅字母、数字、下划线、中划线
chain	string	是	链 code，如 ethereum、base、tron
crypto	string	是	币种 symbol，如 USDT
请求示例
http

复制
GET /v1/deposit/address?uid=user-10001&chain=base&crypto=USDC
GET 请求签名时 request_body 为空字符串。响应：

json

复制
{
  "deposit_address": "0xAbCd1234..."
}
请求的链、币种与链上币种关系必须均已启用，且项目的测试/主网属性必须与链匹配。限流 60/分钟，按 appid + IP 维度。

Webhook 回调
Xcash 在账单收款或充值收款进入关键状态时向商户投递 Webhook。生产环境默认只允许投递到 HTTPS 公网地址，拒绝 http、localhost 和私有网段。

签名头
Xcash 原生协议事件使用 POST application/json，带 HMAC 头，签名算法与 API 请求一致：

http

复制
XC-Appid: {appid}
XC-Nonce: {event_nonce}
XC-Timestamp: {unix_timestamp}
XC-Signature: {hmac_signature}
Content-Type: application/json
响应与重试
成功响应：HTTP 200，响应体去除首尾空白后等于 ok（EPay V1 通知为 success）。
单次请求超时 5 秒；只有网络错误或 5xx 会按指数退避重试，2xx 非 200、3xx、4xx 不重试。
商户应验证签名，并对同一 XC-Nonce 做幂等处理。项目通知开关必须开启，否则即使传入独立 notify_url 也不投递。
账单收款 Webhook
账单进入 completed 后发送一次通知（confirmed=true）：

json

复制
{
  "type": "invoice",
  "data": {
    "sys_no": "INV2606028X7K2P9Q",
    "out_no": "order-20260602-001",
    "crypto": "USDT",
    "chain": "ethereum",
    "pay_address": "0xAbCd1234...",
    "pay_amount": "29.870001",
    "hash": "0xabc123...",
    "block": 12345678,
    "confirmed": true,
    "risk_level": null,
    "risk_score": null
  }
}
充值收款 Webhook
充值对应链上转账达到确认要求后发送一次通知（confirmed=true）：

json

复制
{
  "type": "deposit",
  "data": {
    "sys_no": "DXC2606026K9P2QWX",
    "uid": "user-10001",
    "chain": "base",
    "block": 12345678,
    "hash": "0xabc123...",
    "crypto": "USDC",
    "amount": "500",
    "confirmed": true,
    "risk_level": null,
    "risk_score": null
  }
}