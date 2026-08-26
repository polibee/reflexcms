"""CoinPayments HTTP transport - signs, calls, returns JSON.

"""
import json
import re
from typing import Any, Dict, Optional, Tuple, Union
from urllib.parse import urlencode, quote

import requests

from ._signer import sign

_PLACEHOLDER_RE = re.compile(r":[A-Za-z_]\w*")


class CoinPaymentsError(Exception):
    def __init__(self, status: int, message: str, payload: Any = None):
        super().__init__(f"[{status}] {message}")
        self.status = status
        self.payload = payload


# (connect_timeout, read_timeout) in seconds.
DEFAULT_TIMEOUT: Tuple[float, float] = (10.0, 30.0)


def _query_str(v: Any) -> str:
    if isinstance(v, bool):
        return "true" if v else "false"
    return str(v)


class HttpClient:
    def __init__(
        self,
        base_url: str,
        client_id: Optional[str] = None,
        client_secret: Optional[str] = None,
        timeout: Union[float, Tuple[float, float]] = DEFAULT_TIMEOUT,
    ):
        self.base_url = base_url.rstrip("/")
        self.client_id = client_id
        self.client_secret = client_secret
        self.timeout = timeout
        self._session = requests.Session()

    def call(
        self,
        method: str,
        path: str,
        *,
        path_params: Optional[Dict[str, Any]] = None,
        query: Optional[Dict[str, Any]] = None,
        body: Any = None,
        authed: bool = True,
    ) -> Any:
        # 1. Substitute path params. Routes use Express-style `:name` tokens -
        #    `\b` keeps `:id` from clobbering a longer name like `:identity`.
        if path_params:
            for key, value in path_params.items():
                pattern = re.compile(r":" + re.escape(key) + r"\b")
                new_path, n = pattern.subn(quote(str(value), safe=""), path)
                if n == 0:
                    raise ValueError(f"path template {path!r} has no placeholder :{key}")
                path = new_path
        leftover = _PLACEHOLDER_RE.search(path)
        if leftover:
            raise ValueError(f"unresolved path placeholder {leftover.group(0)} in {path}")

        # 2. Build URL (including query string - it's part of the signed message).
        qs = ""
        if query:
            flat: Dict[str, str] = {}
            for k, v in query.items():
                if v is None:
                    continue
                if isinstance(v, (list, tuple)):
                    flat[k] = ",".join(_query_str(x) for x in v)
                else:
                    flat[k] = _query_str(v)
            if flat:
                qs = "?" + urlencode(flat)
        url = self.base_url + path + qs

        # 3. Serialize body exactly once. The same string is signed and sent.
        payload = "" if body is None else json.dumps(body, separators=(",", ":"))

        # 4. Headers.
        headers = {"Accept": "application/json"}
        if payload:
            headers["Content-Type"] = "application/json"

        if authed:
            if not (self.client_id and self.client_secret):
                raise RuntimeError(f"{method} {path} requires client credentials")
            ts, sig = sign(method, url, self.client_id, self.client_secret, payload)
            headers["X-CoinPayments-Client"] = self.client_id
            headers["X-CoinPayments-Timestamp"] = ts
            headers["X-CoinPayments-Signature"] = sig

        # 5. Send. `data=payload` ensures the wire bytes match what we signed.
        resp = self._session.request(
            method, url, data=payload or None, headers=headers, timeout=self.timeout,
        )

        # 6. Decode.
        if resp.status_code >= 400:
            detail: Any = None
            try:
                detail = resp.json()
            except ValueError:
                detail = resp.text
            raise CoinPaymentsError(resp.status_code, resp.reason or "", detail)

        if not resp.content:
            return None
        return resp.json()
