"""Auto-generated. Do not edit by hand."""
from typing import Any, Optional
from ._http import HttpClient


class CurrenciesApi:
    def __init__(self, http: HttpClient):
        self._http = http

    def get_currencies_v1(self, *, q: Optional[Any] = None, types: Optional[Any] = None, capabilities: Optional[Any] = None) -> Any:
        """lists platform supported currencies and their capabilities. - Supports Auth methods: Anonymous

        See: https://docs.coinpayments.net/api/currencies/routes/getCurrenciesV1

        Args:
            q: optional search query to find currencies with names and/or codes similar to the specified search string
            types: comma separated list of the types of currencies to return (e.g. 'coin', 'token', 'fiat', etc.).  By default currencies of all types are returned
            capabilities: comma separated list of capabilities, currencies without the specified capabilities won't be returned
        """
        query = {"q": q, "types": types, "capabilities": capabilities}
        return self._http.call("GET", "/api/v1/currencies", query=query, authed=False)

    def get_currencies_v2(self, *, q: Optional[Any] = None, types: Optional[Any] = None, capabilities: Optional[Any] = None) -> Any:
        """lists platform supported currencies and their capabilities. - Supports Auth methods: Anonymous

        See: https://docs.coinpayments.net/api/currencies/routes/getCurrenciesV2

        Args:
            q: optional search query to find currencies with names and/or codes similar to the specified search string
            types: comma separated list of the types of currencies to return (e.g. 'coin', 'token', 'fiat', etc.).  By default currencies of all types are returned
            capabilities: comma separated list of capabilities, currencies without the specified capabilities won't be returned
        """
        query = {"q": q, "types": types, "capabilities": capabilities}
        return self._http.call("GET", "/api/v2/currencies", query=query, authed=False)

    def get_currencies_by_id_v1(self, id: Any) -> Any:
        """finds a currency by its id - Supports Auth methods: Anonymous

        See: https://docs.coinpayments.net/api/currencies/routes/getCurrenciesByIdV1

        Path:
            id: the id of the currency to retrieve
        """
        path_params = {"id": id}
        return self._http.call("GET", "/api/v1/currencies/:id", path_params=path_params, authed=False)

    def get_currencies_by_id_v2(self, id: Any) -> Any:
        """finds a currency by its id - Supports Auth methods: Anonymous

        See: https://docs.coinpayments.net/api/currencies/routes/getCurrenciesByIdV2

        Path:
            id: the id of the currency to retrieve
        """
        path_params = {"id": id}
        return self._http.call("GET", "/api/v2/currencies/:id", path_params=path_params, authed=False)

    def get_merchant_currencies_v1(self) -> Any:
        """Gets the merchant's currently accepted currencies.
        Currencies that are ranked (ordered) will be returned at the top of the list. - Supports Auth methods: [ OAuth, clientId/secret ]

        See: https://docs.coinpayments.net/api/currencies/routes/getMerchantCurrenciesV1
        """
        return self._http.call("GET", "/api/v1/merchant/currencies", authed=True)

    def get_currencies_blockchain_nodes_latest_block_number_by_id_v1(self, id: Any) -> Any:
        """Gets the latest blockchain block number by currency - Supports Auth methods: Anonymous

        See: https://docs.coinpayments.net/api/currencies/routes/getCurrenciesBlockchainNodesLatestBlockNumberByIdV1

        Path:
            id: ID of the currency.
        """
        path_params = {"id": id}
        return self._http.call("GET", "/api/v1/currencies/blockchain-nodes/:id/latest-block-number", path_params=path_params, authed=False)

    def get_currencies_blockchain_nodes_latest_block_number_by_id_v2(self, id: Any) -> Any:
        """Gets the latest blockchain block number by currency - Supports Auth methods: Anonymous

        See: https://docs.coinpayments.net/api/currencies/routes/getCurrenciesBlockchainNodesLatestBlockNumberByIdV2

        Path:
            id: ID of the currency.
        """
        path_params = {"id": id}
        return self._http.call("GET", "/api/v2/currencies/blockchain-nodes/:id/latest-block-number", path_params=path_params, authed=False)

    def get_currencies_required_confirmations_v1(self) -> Any:
        """Gets the required confirmations for each currency - Supports Auth methods: Anonymous

        See: https://docs.coinpayments.net/api/currencies/routes/getCurrenciesRequiredConfirmationsV1
        """
        return self._http.call("GET", "/api/v1/currencies/required-confirmations", authed=False)

    def get_currencies_required_confirmations_v2(self) -> Any:
        """Gets the required confirmations for each currency - Supports Auth methods: Anonymous

        See: https://docs.coinpayments.net/api/currencies/routes/getCurrenciesRequiredConfirmationsV2
        """
        return self._http.call("GET", "/api/v2/currencies/required-confirmations", authed=False)

    def get_currencies_conversions_v1(self) -> Any:
        """gets a list of all possible currency conversions - Supports Auth methods: Anonymous

        See: https://docs.coinpayments.net/api/currencies/routes/getCurrenciesConversionsV1
        """
        return self._http.call("GET", "/api/v1/currencies/conversions", authed=False)

    def get_currencies_conversions_v2(self) -> Any:
        """gets a list of all possible currency conversions - Supports Auth methods: Anonymous

        See: https://docs.coinpayments.net/api/currencies/routes/getCurrenciesConversionsV2
        """
        return self._http.call("GET", "/api/v2/currencies/conversions", authed=False)

    def get_currencies_limits_by_from_to_v1(self, from_: Any, to: Any) -> Any:
        """Returns conversion limits by currency pair - Supports Auth methods: Anonymous

        See: https://docs.coinpayments.net/api/currencies/routes/getCurrenciesLimitsByFromToV1

        Path:
            from_: From currency in format `1` or `4:0xdac17f958d2ee523a2206206994597c13d831ec7` for smart contracts
            to: To currency in format `1` or `4:0xdac17f958d2ee523a2206206994597c13d831ec7` for smart contracts
        """
        path_params = {"from": from_, "to": to}
        return self._http.call("GET", "/api/v1/currencies/limits/:from/:to", path_params=path_params, authed=False)

    def get_currencies_limits_by_from_to_v2(self, from_: Any, to: Any) -> Any:
        """Returns conversion limits by currency pair - Supports Auth methods: Anonymous

        See: https://docs.coinpayments.net/api/currencies/routes/getCurrenciesLimitsByFromToV2

        Path:
            from_: From currency in format `1` or `4:0xdac17f958d2ee523a2206206994597c13d831ec7` for smart contracts
            to: To currency in format `1` or `4:0xdac17f958d2ee523a2206206994597c13d831ec7` for smart contracts
        """
        path_params = {"from": from_, "to": to}
        return self._http.call("GET", "/api/v2/currencies/limits/:from/:to", path_params=path_params, authed=False)
