"""Webhook signature verification.

CoinPayments webhooks carry the same three signing headers
as outbound requests, so validation is the same HMAC math in reverse -
reconstruct the expected signature and constant-time compare.

https://docs.coinpayments.net/api/webhooks/authenticating-requests/
"""
import base64
import hashlib
import hmac
from datetime import datetime, timezone
from typing import Mapping, Optional, Union

_BOM = "\ufeff"


class WebhookResult:
    """Rich result so you can log *why* a webhook was rejected."""
    __slots__ = ("ok", "reason")

    def __init__(self, ok: bool, reason: str = ""):
        self.ok = ok
        self.reason = reason

    def __bool__(self) -> bool:
        return self.ok

    def __repr__(self) -> str:
        return f"WebhookResult(ok={self.ok}, reason={self.reason!r})"


def verify_webhook(
    method: str,
    url: str,
    raw_body: Union[str, bytes, bytearray],
    headers: Mapping[str, str],
    client_secret: str,
    *,
    expected_client_id: Optional[str] = None,
    max_age_seconds: int = 300,
) -> WebhookResult:
    """Verify a CoinPayments webhook request.

    Args:
        method: HTTP method the webhook arrived with (typically "POST").
        url: The full URL the webhook was delivered to (scheme + host + path +
             query). Use the URL you registered with CoinPayments - not what
             your framework reconstructs, which may differ behind proxies.
        raw_body: Request body BEFORE any parsing. Accepts either `str` or
                  `bytes` - frameworks differ (Django/FastAPI return bytes,
                  Flask returns str). JSON parsers may re-serialise with
                  different whitespace, which invalidates the signature -
                  capture the body before decoding.
        headers: Request headers. Case-insensitive lookup is performed for
                 the three X-CoinPayments-* fields.
        client_secret: Your integration's client secret.
        expected_client_id: If provided, reject webhooks whose Client header
                            does not match. Useful when one endpoint serves
                            multiple integrations.
        max_age_seconds: Reject webhooks whose timestamp drifts from `now`
                         by more than this many seconds. Set to 0 to skip
                         the freshness check.
    """
    # Normalise bytes → str so f-string concatenation below is well-defined.
    # Django's `request.body` and FastAPI's `await request.body()` both return
    # bytes; without this, f-string interpolation would inject the repr of the
    # bytes object (`b'...'`) and every signature would mismatch.
    if isinstance(raw_body, (bytes, bytearray)):
        raw_body = raw_body.decode("utf-8")

    # Case-insensitive header lookup.
    h = {k.lower(): v for k, v in headers.items()}
    client_id = h.get("x-coinpayments-client")
    ts        = h.get("x-coinpayments-timestamp")
    sig       = h.get("x-coinpayments-signature")

    if not (client_id and ts and sig):
        return WebhookResult(False, "missing signature headers")

    if expected_client_id is not None and client_id != expected_client_id:
        return WebhookResult(False, "client id mismatch")

    if max_age_seconds > 0:
        try:
            sent = datetime.strptime(ts, "%Y-%m-%dT%H:%M:%S").replace(tzinfo=timezone.utc)
        except ValueError:
            return WebhookResult(False, "malformed timestamp")
        age = abs((datetime.now(timezone.utc) - sent).total_seconds())
        if age > max_age_seconds:
            return WebhookResult(False, f"timestamp too old ({int(age)}s)")

    message = f"{_BOM}{method}{url}{client_id}{ts}{raw_body}".encode("utf-8")
    expected = base64.b64encode(
        hmac.new(client_secret.encode("utf-8"), message, hashlib.sha256).digest()
    ).decode("ascii")

    # Constant-time comparison - never use plain `==` for signature checks.
    if not hmac.compare_digest(expected, sig):
        return WebhookResult(False, "signature mismatch")

    return WebhookResult(True)
