"""Auto-generated. Do not edit by hand."""
from typing import Any, Optional
from ._http import HttpClient


class FeesApi:
    def __init__(self, http: HttpClient):
        self._http = http

    def get_fees_blockchain_by_id_v1(self, id: Any) -> Any:
        """Returns the current blockchain fee for sending funds from the CoinPayments system to an external address - Supports Auth methods: Anonymous

        See: https://docs.coinpayments.net/api/fees/routes/getFeesBlockchainByIdV1

        Path:
            id: Currency in format `1` or `4:0xdac17f958d2ee523a2206206994597c13d831ec7` for smart contracts
        """
        path_params = {"id": id}
        return self._http.call("GET", "/api/v1/fees/blockchain/:id", path_params=path_params, authed=False)

    def get_fees_blockchain_by_id_v2(self, id: Any) -> Any:
        """Returns the current blockchain fee for sending funds from the CoinPayments system to an external address - Supports Auth methods: Anonymous

        See: https://docs.coinpayments.net/api/fees/routes/getFeesBlockchainByIdV2

        Path:
            id: Currency in format `1` or `4:0xdac17f958d2ee523a2206206994597c13d831ec7` for smart contracts
        """
        path_params = {"id": id}
        return self._http.call("GET", "/api/v2/fees/blockchain/:id", path_params=path_params, authed=False)
