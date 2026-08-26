"""CoinPayments HMAC-SHA256 request signer.

Spec: https://docs.coinpayments.net/api/auth/generate-api-signature
"""
import base64
import hashlib
import hmac
from datetime import datetime, timezone
from typing import Tuple

_BOM = "﻿"


def sign(
    method: str,
    url: str,
    client_id: str,
    client_secret: str,
    payload: str,
) -> Tuple[str, str]:
    """Return (timestamp, signature) for a CoinPayments API request.

    The timestamp and signature must be sent as X-CoinPayments-Timestamp
    and X-CoinPayments-Signature headers, together with X-CoinPayments-Client.

    IMPORTANT: `payload` MUST be the exact JSON string that will be sent as
    the request body - any whitespace drift invalidates the signature. Pass
    "" for a request with no body. The HTTP client is responsible for
    serializing once and signing the same bytes that go on the wire.
    """
    ts = datetime.now(timezone.utc).strftime("%Y-%m-%dT%H:%M:%S")
    message = f"{_BOM}{method}{url}{client_id}{ts}{payload}".encode("utf-8")
    digest = hmac.new(
        client_secret.encode("utf-8"),
        message,
        hashlib.sha256,
    ).digest()
    return ts, base64.b64encode(digest).decode("ascii")
