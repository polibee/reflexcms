"""Auto-generated. Do not edit by hand."""
from typing import Any, Optional
from ._http import HttpClient


class AccessControlApi:
    def __init__(self, http: HttpClient):
        self._http = http

    def post_merchant_clients_access_control_v1(self, body: dict) -> Any:
        """Allows merchant clients to downgrade their access scope by reducing permissions - Supports Auth methods: [ OAuth, clientId/secret ]

        See: https://docs.coinpayments.net/api/accessControl/routes/postMerchantClientsAccessControlV1

        Body: UpdateMerchantClientAccessControlDto
        """
        return self._http.call("POST", "/api/v1/merchant/clients/access-control", body=body, authed=True)

    def get_merchant_clients_access_control_v1(self) -> Any:
        """Retrieves the current allowed access control scope for the authenticated merchant client - Supports Auth methods: [ OAuth, clientId/secret ]

        See: https://docs.coinpayments.net/api/accessControl/routes/getMerchantClientsAccessControlV1
        """
        return self._http.call("GET", "/api/v1/merchant/clients/access-control", authed=True)
