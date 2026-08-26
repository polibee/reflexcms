"""Auto-generated. Do not edit by hand."""
from typing import Any, Optional
from ._http import HttpClient


class WalletsApi:
    def __init__(self, http: HttpClient):
        self._http = http

    def post_merchant_wallets_v1(self, body: dict) -> Any:
        """Creates a new merchant wallet - Supports Auth methods: [ OAuth, clientId/secret ]

        See: https://docs.coinpayments.net/api/wallets/routes/postMerchantWalletsV1

        Body: NewWalletRequestDto
        """
        return self._http.call("POST", "/api/v1/merchant/wallets", body=body, authed=True)

    def post_merchant_wallets_v2(self, body: dict) -> Any:
        """Creates a new merchant wallet - Supports Auth methods: [ OAuth, clientId/secret ]

        See: https://docs.coinpayments.net/api/wallets/routes/postMerchantWalletsV2

        Body: NewWalletRequestV2Dto
        """
        return self._http.call("POST", "/api/v2/merchant/wallets", body=body, authed=True)

    def get_merchant_wallets_v1(self, *, skip: Optional[Any] = None, take: Optional[Any] = None) -> Any:
        """Lists merchant client wallets (up to 100 wallets per response) - Supports Auth methods: [ OAuth, clientId/secret ]

        See: https://docs.coinpayments.net/api/wallets/routes/getMerchantWalletsV1

        Args:
            skip: How may transaction to skip (used for paging)
            take: How may transaction to take (used for paging)
        """
        query = {"skip": skip, "take": take}
        return self._http.call("GET", "/api/v1/merchant/wallets", query=query, authed=True)

    def get_merchant_wallets_v2(self, *, skip: Optional[Any] = None, take: Optional[Any] = None) -> Any:
        """Lists merchant client wallets (up to 100 wallets per response) - Supports Auth methods: [ OAuth, clientId/secret ]

        See: https://docs.coinpayments.net/api/wallets/routes/getMerchantWalletsV2

        Args:
            skip: how many wallets to skip (used for paging)
            take: how many wallets to take (used for paging)
        """
        query = {"skip": skip, "take": take}
        return self._http.call("GET", "/api/v2/merchant/wallets", query=query, authed=True)

    def get_merchant_wallets_by_id_v1(self, id: Any) -> Any:
        """Finds a merchant client wallet by id - Supports Auth methods: [ OAuth, clientId/secret ]

        See: https://docs.coinpayments.net/api/wallets/routes/getMerchantWalletsByIdV1

        Path:
            id: the id of the wallet
        """
        path_params = {"id": id}
        return self._http.call("GET", "/api/v1/merchant/wallets/:id", path_params=path_params, authed=True)

    def get_merchant_wallets_by_id_v2(self, id: Any) -> Any:
        """Finds a merchant client wallet by id - Supports Auth methods: [ OAuth, clientId/secret ]

        See: https://docs.coinpayments.net/api/wallets/routes/getMerchantWalletsByIdV2

        Path:
            id: the id of the wallet
        """
        path_params = {"id": id}
        return self._http.call("GET", "/api/v2/merchant/wallets/:id", path_params=path_params, authed=True)

    def post_merchant_wallets_addresses_by_id_v1(self, id: Any, body: dict) -> Any:
        """Creates an address under a wallet by the wallet ID - Supports Auth methods: [ OAuth, clientId/secret ]

        See: https://docs.coinpayments.net/api/wallets/routes/postMerchantWalletsAddressesByIdV1

        Path:
            id: the ID of the wallet to create the address under

        Body: CreateMerchantWalletAddressRequestDto
        """
        path_params = {"id": id}
        return self._http.call("POST", "/api/v1/merchant/wallets/:id/addresses", path_params=path_params, body=body, authed=True)

    def post_merchant_wallets_addresses_by_id_v2(self, id: Any, body: dict) -> Any:
        """Creates an address under a wallet by the wallet ID - Supports Auth methods: [ OAuth, clientId/secret ]

        See: https://docs.coinpayments.net/api/wallets/routes/postMerchantWalletsAddressesByIdV2

        Path:
            id: the id of the wallet to create the address under

        Body: CreateMerchantWalletAddressRequestDto
        """
        path_params = {"id": id}
        return self._http.call("POST", "/api/v2/merchant/wallets/:id/addresses", path_params=path_params, body=body, authed=True)

    def get_merchant_wallets_addresses_by_id_v1(self, id: Any, *, skip: Optional[Any] = None, take: Optional[Any] = None) -> Any:
        """Lists all merchant addresses of a specific wallet by the wallet Id - Supports Auth methods: [ OAuth, clientId/secret ]

        See: https://docs.coinpayments.net/api/wallets/routes/getMerchantWalletsAddressesByIdV1

        Path:
            id: the id of the wallet

        Args:
            skip: how many addresses to skip (used for paging)
            take: how many addresses to take (used for paging)
        """
        path_params = {"id": id}
        query = {"skip": skip, "take": take}
        return self._http.call("GET", "/api/v1/merchant/wallets/:id/addresses", path_params=path_params, query=query, authed=True)

    def get_merchant_wallets_addresses_by_id_v2(self, id: Any, *, skip: Optional[Any] = None, take: Optional[Any] = None) -> Any:
        """Lists all merchant addresses of a specific wallet by the wallet Id - Supports Auth methods: [ OAuth, clientId/secret ]

        See: https://docs.coinpayments.net/api/wallets/routes/getMerchantWalletsAddressesByIdV2

        Path:
            id: the id of the wallet

        Args:
            skip: how many addresses to skip (used for paging)
            take: how many addresses to take (used for paging)
        """
        path_params = {"id": id}
        query = {"skip": skip, "take": take}
        return self._http.call("GET", "/api/v2/merchant/wallets/:id/addresses", path_params=path_params, query=query, authed=True)

    def get_merchant_wallets_addresses_by_id_a_id_v1(self, id: Any, a_id: Any) -> Any:
        """Get a specific address by its Id and the Id of the wallet it belongs to - Supports Auth methods: [ OAuth, clientId/secret ]

        See: https://docs.coinpayments.net/api/wallets/routes/getMerchantWalletsAddressesByIdAIdV1

        Path:
            id: the id of the wallet containing the address
            a_id: the id of the address to retrieve
        """
        path_params = {"id": id, "aId": a_id}
        return self._http.call("GET", "/api/v1/merchant/wallets/:id/addresses/:aId", path_params=path_params, authed=True)

    def get_merchant_wallets_addresses_by_id_a_id_v2(self, id: Any, a_id: Any) -> Any:
        """Get a specific address by its Id and the Id of the wallet it belongs to - Supports Auth methods: [ OAuth, clientId/secret ]

        See: https://docs.coinpayments.net/api/wallets/routes/getMerchantWalletsAddressesByIdAIdV2

        Path:
            id: the id of the wallet containing the address
            a_id: the id of the address to retrieve
        """
        path_params = {"id": id, "aId": a_id}
        return self._http.call("GET", "/api/v2/merchant/wallets/:id/addresses/:aId", path_params=path_params, authed=True)

    def get_merchant_wallets_count_v2(self) -> Any:
        """Get merchant client wallets count - Supports Auth methods: [ OAuth, clientId/secret ]

        See: https://docs.coinpayments.net/api/wallets/routes/getMerchantWalletsCountV2
        """
        return self._http.call("GET", "/api/v2/merchant/wallets/count", authed=True)

    def get_merchant_wallets_addresses_count_by_id_v2(self, id: Any) -> Any:
        """Get merchant addresses count of a specific wallet by the wallet Id - Supports Auth methods: [ OAuth, clientId/secret ]

        See: https://docs.coinpayments.net/api/wallets/routes/getMerchantWalletsAddressesCountByIdV2

        Path:
            id: the id of the wallet
        """
        path_params = {"id": id}
        return self._http.call("GET", "/api/v2/merchant/wallets/:id/addresses/count", path_params=path_params, authed=True)

    def get_merchant_wallets_balance_by_id_date_v1(self, id: Any, date: Any) -> Any:
        """Returns the wallet balance in a specific date time offset - Supports Auth methods: [ OAuth, clientId/secret ]

        See: https://docs.coinpayments.net/api/wallets/routes/getMerchantWalletsBalanceByIdDateV1

        Path:
            id: the wallet id
            date: the specific date time offset to retrieve what the wallet balance was at that time
        """
        path_params = {"id": id, "date": date}
        return self._http.call("GET", "/api/v1/merchant/wallets/:id/balance/:date", path_params=path_params, authed=True)

    def get_merchant_wallets_balance_by_id_date_v2(self, id: Any, date: Any) -> Any:
        """Returns the wallet balance in a specific date time offset - Supports Auth methods: [ OAuth, clientId/secret ]

        See: https://docs.coinpayments.net/api/wallets/routes/getMerchantWalletsBalanceByIdDateV2

        Path:
            id: the wallet id
            date: the specific date time offset to retrieve what the wallet balance was at that time
        """
        path_params = {"id": id, "date": date}
        return self._http.call("GET", "/api/v2/merchant/wallets/:id/balance/:date", path_params=path_params, authed=True)

    def get_merchant_wallets_v3(self, *, skip: Optional[Any] = None, take: Optional[Any] = None) -> Any:
        """Lists merchant client wallets (up to 100 wallets per response) - Supports Auth methods: [ OAuth, clientId/secret ]

        See: https://docs.coinpayments.net/api/wallets/routes/getMerchantWalletsV3

        Args:
            skip: how many wallets to skip (used for paging)
            take: how many wallets to take (used for paging)
        """
        query = {"skip": skip, "take": take}
        return self._http.call("GET", "/api/v3/merchant/wallets", query=query, authed=True)

    def put_merchant_wallets_v3(self, body: dict) -> Any:
        """Creates or retrieves a wallet and address by external Ids - Supports Auth methods: [ OAuth, clientId/secret ]

        See: https://docs.coinpayments.net/api/wallets/routes/putMerchantWalletsV3

        Body: NewWalletAndAddressRequestV3Dto
        """
        return self._http.call("PUT", "/api/v3/merchant/wallets", body=body, authed=True)

    def get_merchant_wallets_count_v3(self) -> Any:
        """Get merchant client wallets count - Supports Auth methods: [ OAuth, clientId/secret ]

        See: https://docs.coinpayments.net/api/wallets/routes/getMerchantWalletsCountV3
        """
        return self._http.call("GET", "/api/v3/merchant/wallets/count", authed=True)

    def get_merchant_wallets_addresses_count_by_label_currency_v3(self, label: Any, currency: Any) -> Any:
        """Get merchant addresses count of a specific wallet by the wallet label - Supports Auth methods: [ OAuth, clientId/secret ]

        See: https://docs.coinpayments.net/api/wallets/routes/getMerchantWalletsAddressesCountByLabelCurrencyV3

        Path:
            label: the unique label for the wallet provided by client
            currency: the currency of the wallet
        """
        path_params = {"label": label, "currency": currency}
        return self._http.call("GET", "/api/v3/merchant/wallets/:label/:currency/addresses/count", path_params=path_params, authed=True)

    def get_merchant_wallets_addresses_by_label_currency_v3(self, label: Any, currency: Any, *, skip: Optional[Any] = None, take: Optional[Any] = None) -> Any:
        """Lists all merchant addresses of a specific wallet by unique wallet label and currency Id - Supports Auth methods: [ OAuth, clientId/secret ]

        See: https://docs.coinpayments.net/api/wallets/routes/getMerchantWalletsAddressesByLabelCurrencyV3

        Path:
            label: the unique id for the wallet provided by client
            currency: the currency of the wallet to list addresses for

        Args:
            skip: How many addresses to skip (used for paging)
            take: How many addresses to take (used for paging)
        """
        path_params = {"label": label, "currency": currency}
        query = {"skip": skip, "take": take}
        return self._http.call("GET", "/api/v3/merchant/wallets/:label/:currency/addresses", path_params=path_params, query=query, authed=True)

    def get_merchant_wallets_addresses_by_w_label_currency_a_label_v3(self, w_label: Any, currency: Any, a_label: Any) -> Any:
        """Get a specific address by its label and the label of the wallet it belongs to - Supports Auth methods: [ OAuth, clientId/secret ]

        See: https://docs.coinpayments.net/api/wallets/routes/getMerchantWalletsAddressesByWLabelCurrencyALabelV3

        Path:
            w_label: the unique label for the wallet provided by client
            currency: the currency of the wallet containing the address
            a_label: the unique label for the address provided by client
        """
        path_params = {"wLabel": w_label, "currency": currency, "aLabel": a_label}
        return self._http.call("GET", "/api/v3/merchant/wallets/:wLabel/:currency/addresses/:aLabel", path_params=path_params, authed=True)
