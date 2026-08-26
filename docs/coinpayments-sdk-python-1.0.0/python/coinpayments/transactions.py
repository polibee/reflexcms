"""Auto-generated. Do not edit by hand."""
from typing import Any, Optional
from ._http import HttpClient


class TransactionsApi:
    def __init__(self, http: HttpClient):
        self._http = http

    def get_merchant_wallets_consolidation_transactions_by_id_spend_id_v1(self, id: Any, spend_id: Any) -> Any:
        """Lists transactions used for consolidating funds to pay fee for the withdrawal - Supports Auth methods: [ OAuth, clientId/secret ]

        See: https://docs.coinpayments.net/api/transactions/routes/getMerchantWalletsConsolidationTransactionsByIdSpendIdV1

        Path:
            id: the id of the wallet
            spend_id: ID of the withdrawal (spend request)
        """
        path_params = {"id": id, "spendId": spend_id}
        return self._http.call("GET", "/api/v1/merchant/wallets/:id/consolidation-transactions/:spendId", path_params=path_params, authed=True)

    def get_merchant_wallets_consolidation_transactions_by_id_spend_id_v2(self, id: Any, spend_id: Any) -> Any:
        """Lists transactions used for consolidating funds to pay fee for the withdrawal - Supports Auth methods: [ OAuth, clientId/secret ]

        See: https://docs.coinpayments.net/api/transactions/routes/getMerchantWalletsConsolidationTransactionsByIdSpendIdV2

        Path:
            id: the id of the wallet
            spend_id: ID of the withdrawal (spend request)
        """
        path_params = {"id": id, "spendId": spend_id}
        return self._http.call("GET", "/api/v2/merchant/wallets/:id/consolidation-transactions/:spendId", path_params=path_params, authed=True)

    def get_merchant_wallets_transactions_by_id_v1(self, id: Any, *, skip: Optional[Any] = None, take: Optional[Any] = None) -> Any:
        """Lists transactions of the wallet - Supports Auth methods: [ OAuth, clientId/secret ]

        See: https://docs.coinpayments.net/api/transactions/routes/getMerchantWalletsTransactionsByIdV1

        Path:
            id: the id of the wallet

        Args:
            skip: how many transactions to skip (used for paging)
            take: how many transactions to take (used for paging)
        """
        path_params = {"id": id}
        query = {"skip": skip, "take": take}
        return self._http.call("GET", "/api/v1/merchant/wallets/:id/transactions", path_params=path_params, query=query, authed=True)

    def get_merchant_wallets_transactions_by_id_v2(self, id: Any, *, skip: Optional[Any] = None, take: Optional[Any] = None) -> Any:
        """Lists transactions of the wallet - Supports Auth methods: [ OAuth, clientId/secret ]

        See: https://docs.coinpayments.net/api/transactions/routes/getMerchantWalletsTransactionsByIdV2

        Path:
            id: the id of the wallet

        Args:
            skip: how many transactions to skip (used for paging)
            take: how many transactions to take (used for paging)
        """
        path_params = {"id": id}
        query = {"skip": skip, "take": take}
        return self._http.call("GET", "/api/v2/merchant/wallets/:id/transactions", path_params=path_params, query=query, authed=True)

    def get_merchant_wallets_transaction_by_id_v1(self, id: Any, *, transaction_id: Optional[Any] = None, spend_request_id: Optional[Any] = None) -> Any:
        """Get a specific transaction of the wallet, if transactionId is specified for search then the spendRequestId is ignored otherwise a first spending transaction with matching spendRequestId is returned - Supports Auth methods: [ OAuth, clientId/secret ]

        See: https://docs.coinpayments.net/api/transactions/routes/getMerchantWalletsTransactionByIdV1

        Path:
            id: the id of the wallet

        Args:
            transaction_id: Id of the transaction to look for
            spend_request_id: SpendRequestId of the transaction to look for
        """
        path_params = {"id": id}
        query = {"transactionId": transaction_id, "spendRequestId": spend_request_id}
        return self._http.call("GET", "/api/v1/merchant/wallets/:id/transaction", path_params=path_params, query=query, authed=True)

    def get_merchant_wallets_transaction_by_id_v2(self, id: Any, *, transaction_id: Optional[Any] = None, spend_request_id: Optional[Any] = None) -> Any:
        """Get a specific transaction of the wallet, if transactionId is specified for search then the spendRequestId is ignored otherwise a first spending transaction with matching spendRequestId is returned - Supports Auth methods: [ OAuth, clientId/secret ]

        See: https://docs.coinpayments.net/api/transactions/routes/getMerchantWalletsTransactionByIdV2

        Path:
            id: the id of the wallet

        Args:
            transaction_id: Id of the transaction to look for
            spend_request_id: SpendRequestId of the transaction to look for
        """
        path_params = {"id": id}
        query = {"transactionId": transaction_id, "spendRequestId": spend_request_id}
        return self._http.call("GET", "/api/v2/merchant/wallets/:id/transaction", path_params=path_params, query=query, authed=True)

    def post_merchant_wallets_consolidation_by_id_to_v1(self, id: Any, to: Any, *, address_ids: Optional[Any] = None) -> Any:
        """Executes merchant wallet consolidation (sending funds to the main wallet) - Supports Auth methods: [ OAuth, clientId/secret ]

        See: https://docs.coinpayments.net/api/transactions/routes/postMerchantWalletsConsolidationByIdToV1

        Path:
            id: the id of the wallet which will be sending funds
            to: the id of the wallet which will be receiving funds

        Args:
            address_ids: Comma-separated IDs of addresses for consolidation
        """
        path_params = {"id": id, "to": to}
        query = {"addressIds": address_ids}
        return self._http.call("POST", "/api/v1/merchant/wallets/:id/consolidation/:to", path_params=path_params, query=query, authed=True)

    def post_merchant_wallets_consolidation_by_id_to_v2(self, id: Any, to: Any, *, address_ids: Optional[Any] = None) -> Any:
        """Executes merchant wallet consolidation (sending funds to the main wallet) - Supports Auth methods: [ OAuth, clientId/secret ]

        See: https://docs.coinpayments.net/api/transactions/routes/postMerchantWalletsConsolidationByIdToV2

        Path:
            id: the id of the wallet which will be sending funds
            to: the id of the wallet which will be receiving funds

        Args:
            address_ids: Comma-separated IDs of addresses for consolidation
        """
        path_params = {"id": id, "to": to}
        query = {"addressIds": address_ids}
        return self._http.call("POST", "/api/v2/merchant/wallets/:id/consolidation/:to", path_params=path_params, query=query, authed=True)

    def get_merchant_wallets_consolidation_by_id_v1(self, id: Any, *, address_ids: Optional[Any] = None) -> Any:
        """Returns info about merchant wallet consolidation (sending funds to the main wallet) - Supports Auth methods: [ OAuth, clientId/secret ]

        See: https://docs.coinpayments.net/api/transactions/routes/getMerchantWalletsConsolidationByIdV1

        Path:
            id: the id of the wallet which will be sending funds to the main wallet

        Args:
            address_ids: Comma-separated IDs of addresses for consolidation
        """
        path_params = {"id": id}
        query = {"addressIds": address_ids}
        return self._http.call("GET", "/api/v1/merchant/wallets/:id/consolidation", path_params=path_params, query=query, authed=True)

    def get_merchant_wallets_consolidation_by_id_v2(self, id: Any, *, address_ids: Optional[Any] = None) -> Any:
        """Returns info about merchant wallet consolidation (sending funds to the main wallet) - Supports Auth methods: [ OAuth, clientId/secret ]

        See: https://docs.coinpayments.net/api/transactions/routes/getMerchantWalletsConsolidationByIdV2

        Path:
            id: the id of the wallet which will be sending funds to the main wallet

        Args:
            address_ids: Comma-separated IDs of addresses for consolidation
        """
        path_params = {"id": id}
        query = {"addressIds": address_ids}
        return self._http.call("GET", "/api/v2/merchant/wallets/:id/consolidation", path_params=path_params, query=query, authed=True)

    def post_merchant_wallets_consolidation_preview_v1(self, body: dict) -> Any:
        """Preview merchant wallets consolidation (sending funds to the main wallet) - Supports Auth methods: [ OAuth, clientId/secret ]

        See: https://docs.coinpayments.net/api/transactions/routes/postMerchantWalletsConsolidationPreviewV1

        Body: WalletsConsolidationRequestDto
        """
        return self._http.call("POST", "/api/v1/merchant/wallets/consolidation-preview", body=body, authed=True)

    def post_merchant_wallets_consolidation_preview_v2(self, body: dict) -> Any:
        """Preview merchant wallets consolidation (sending funds to the main wallet) - Supports Auth methods: [ OAuth, clientId/secret ]

        See: https://docs.coinpayments.net/api/transactions/routes/postMerchantWalletsConsolidationPreviewV2

        Body: WalletsConsolidationRequestDto
        """
        return self._http.call("POST", "/api/v2/merchant/wallets/consolidation-preview", body=body, authed=True)

    def post_merchant_wallets_consolidation_by_id_v1(self, id: Any, body: dict) -> Any:
        """Execute merchant wallets consolidation (sending funds to the main wallet) - Supports Auth methods: [ OAuth, clientId/secret ]

        See: https://docs.coinpayments.net/api/transactions/routes/postMerchantWalletsConsolidationByIdV1

        Path:
            id: the id of the wallet which will be receiving funds

        Body: WalletsConsolidationRequestDto
        """
        path_params = {"id": id}
        return self._http.call("POST", "/api/v1/merchant/wallets/consolidation/:id", path_params=path_params, body=body, authed=True)

    def post_merchant_wallets_consolidation_by_id_v2(self, id: Any, body: dict) -> Any:
        """Execute merchant wallets consolidation (sending funds to the main wallet) - Supports Auth methods: [ OAuth, clientId/secret ]

        See: https://docs.coinpayments.net/api/transactions/routes/postMerchantWalletsConsolidationByIdV2

        Path:
            id: the ID of the wallet which will be receiving funds (optional)

        Body: WalletsConsolidationRequestDto
        """
        path_params = {"id": id}
        return self._http.call("POST", "/api/v2/merchant/wallets/consolidation/:id", path_params=path_params, body=body, authed=True)

    def post_merchant_wallets_spend_request_by_id_v1(self, id: Any, body: dict) -> Any:
        """Sends a request to spend funds from the merchant client wallet. Also used to convert funds between supported cryptocurrencies - Supports Auth methods: [ OAuth, clientId/secret ]

        See: https://docs.coinpayments.net/api/transactions/routes/postMerchantWalletsSpendRequestByIdV1

        Path:
            id: the id of the wallet from which to spend funds from

        Body: SpendRequestDto
        """
        path_params = {"id": id}
        return self._http.call("POST", "/api/v1/merchant/wallets/:id/spend/request", path_params=path_params, body=body, authed=True)

    def post_merchant_wallets_spend_request_by_id_v2(self, id: Any, body: dict) -> Any:
        """Sends a request to spend funds from the merchant client wallet. Also used to convert funds between supported cryptocurrencies - Supports Auth methods: [ OAuth, clientId/secret ]

        See: https://docs.coinpayments.net/api/transactions/routes/postMerchantWalletsSpendRequestByIdV2

        Path:
            id: the id of the wallet from which to spend funds from

        Body: SpendRequestV2Dto
        """
        path_params = {"id": id}
        return self._http.call("POST", "/api/v2/merchant/wallets/:id/spend/request", path_params=path_params, body=body, authed=True)

    def post_merchant_wallets_spend_confirmation_by_id_v1(self, id: Any, body: dict) -> Any:
        """Sends a request to confirm spending funds from the merchant client wallet, Or to confirm converting funds. - Supports Auth methods: [ OAuth, clientId/secret ]

        See: https://docs.coinpayments.net/api/transactions/routes/postMerchantWalletsSpendConfirmationByIdV1

        Path:
            id: the id of the wallet which to spend funds from

        Body: WalletSpendConfirmationRequestDto
        """
        path_params = {"id": id}
        return self._http.call("POST", "/api/v1/merchant/wallets/:id/spend/confirmation", path_params=path_params, body=body, authed=True)

    def post_merchant_wallets_spend_confirmation_by_id_v2(self, id: Any, body: dict) -> Any:
        """Sends a request to confirm spending funds from the merchant client wallet, Or to confirm converting funds. - Supports Auth methods: [ OAuth, clientId/secret ]

        See: https://docs.coinpayments.net/api/transactions/routes/postMerchantWalletsSpendConfirmationByIdV2

        Path:
            id: the id of the wallet which to spend funds from

        Body: WalletSpendConfirmationRequestDto
        """
        path_params = {"id": id}
        return self._http.call("POST", "/api/v2/merchant/wallets/:id/spend/confirmation", path_params=path_params, body=body, authed=True)

    def get_merchant_wallets_transactions_count_by_id_v2(self, id: Any) -> Any:
        """Get transactions count of the wallet - Supports Auth methods: [ OAuth, clientId/secret ]

        See: https://docs.coinpayments.net/api/transactions/routes/getMerchantWalletsTransactionsCountByIdV2

        Path:
            id: the id of the wallet
        """
        path_params = {"id": id}
        return self._http.call("GET", "/api/v2/merchant/wallets/:id/transactions/count", path_params=path_params, authed=True)

    def get_merchant_wallets_transactions_count_by_label_currency_v3(self, label: Any, currency: Any) -> Any:
        """Get transactions count of the wallet - Supports Auth methods: [ OAuth, clientId/secret ]

        See: https://docs.coinpayments.net/api/transactions/routes/getMerchantWalletsTransactionsCountByLabelCurrencyV3

        Path:
            label: the unique id for the wallet provided by client
            currency: the currency of the wallet to count transactions for
        """
        path_params = {"label": label, "currency": currency}
        return self._http.call("GET", "/api/v3/merchant/wallets/:label/:currency/transactions/count", path_params=path_params, authed=True)

    def get_merchant_wallets_transactions_by_label_currency_v3(self, label: Any, currency: Any, *, skip: Optional[Any] = None, take: Optional[Any] = None) -> Any:
        """Lists transactions of the wallet - Supports Auth methods: [ OAuth, clientId/secret ]

        See: https://docs.coinpayments.net/api/transactions/routes/getMerchantWalletsTransactionsByLabelCurrencyV3

        Path:
            label: the unique id for the wallet provided by client
            currency: the currency of the wallet to list transactions for

        Args:
            skip: how many transactions to skip (used for paging)
            take: how many transactions to take (used for paging)
        """
        path_params = {"label": label, "currency": currency}
        query = {"skip": skip, "take": take}
        return self._http.call("GET", "/api/v3/merchant/wallets/:label/:currency/transactions", path_params=path_params, query=query, authed=True)

    def get_merchant_wallets_transaction_by_label_currency_v3(self, label: Any, currency: Any, *, transaction_id: Optional[Any] = None, spend_request_id: Optional[Any] = None) -> Any:
        """Get a specific transaction of the wallet, if transactionId is specified for search then the spendRequestId is ignored otherwise a first spending transaction with matching spendRequestId is returned - Supports Auth methods: [ OAuth, clientId/secret ]

        See: https://docs.coinpayments.net/api/transactions/routes/getMerchantWalletsTransactionByLabelCurrencyV3

        Path:
            label: the unique id for the wallet provided by client
            currency: the currency of the wallet to search transactions in

        Args:
            transaction_id: Id of the transaction to look for
            spend_request_id: SpendRequestId of the transaction to look for
        """
        path_params = {"label": label, "currency": currency}
        query = {"transactionId": transaction_id, "spendRequestId": spend_request_id}
        return self._http.call("GET", "/api/v3/merchant/wallets/:label/:currency/transaction", path_params=path_params, query=query, authed=True)

    def post_merchant_wallets_spend_request_by_label_currency_v3(self, label: Any, currency: Any, body: dict) -> Any:
        """Sends a request to spend funds from the merchant client wallet. Also used to convert funds between supported cryptocurrencies - Supports Auth methods: [ OAuth, clientId/secret ]

        See: https://docs.coinpayments.net/api/transactions/routes/postMerchantWalletsSpendRequestByLabelCurrencyV3

        Path:
            label: the unique label for the wallet provided by client
            currency: the currency of the wallet to spend funds from

        Body: SpendRequestV2Dto
        """
        path_params = {"label": label, "currency": currency}
        return self._http.call("POST", "/api/v3/merchant/wallets/:label/:currency/spend/request", path_params=path_params, body=body, authed=True)

    def post_merchant_wallets_spend_confirmation_by_label_currency_v3(self, label: Any, currency: Any, body: dict) -> Any:
        """Sends a request to confirm spending funds from the merchant client wallet, Or to confirm converting funds. - Supports Auth methods: [ OAuth, clientId/secret ]

        See: https://docs.coinpayments.net/api/transactions/routes/postMerchantWalletsSpendConfirmationByLabelCurrencyV3

        Path:
            label: the unique label for the wallet provided by client
            currency: the currency of the wallet to confirm spending from

        Body: WalletSpendConfirmationRequestDto
        """
        path_params = {"label": label, "currency": currency}
        return self._http.call("POST", "/api/v3/merchant/wallets/:label/:currency/spend/confirmation", path_params=path_params, body=body, authed=True)

    def get_merchant_wallets_consolidation_by_label_currency_v3(self, label: Any, currency: Any, *, unique_address_labels: Optional[Any] = None) -> Any:
        """Returns info about merchant wallet consolidation (sending funds to the main wallet) - Supports Auth methods: [ OAuth, clientId/secret ]

        See: https://docs.coinpayments.net/api/transactions/routes/getMerchantWalletsConsolidationByLabelCurrencyV3

        Path:
            label: the unique label for the wallet provided by client
            currency: the currency of the wallet to get consolidation info for

        Args:
            unique_address_labels: Comma-separated unique labels of addresses for consolidation
        """
        path_params = {"label": label, "currency": currency}
        query = {"uniqueAddressLabels": unique_address_labels}
        return self._http.call("GET", "/api/v3/merchant/wallets/:label/:currency/consolidation", path_params=path_params, query=query, authed=True)

    def post_merchant_wallets_consolidation_by_label_currency_to_v3(self, label: Any, currency: Any, to: Any, *, unique_address_labels: Optional[Any] = None) -> Any:
        """Executes merchant wallet consolidation (sending funds to the main wallet) - Supports Auth methods: [ OAuth, clientId/secret ]

        See: https://docs.coinpayments.net/api/transactions/routes/postMerchantWalletsConsolidationByLabelCurrencyToV3

        Path:
            label: the unique label for the wallet provided by client
            currency: the currency of the wallet to consolidate funds from
            to: the label of the wallet which will be receiving funds

        Args:
            unique_address_labels: Comma-separated unique address labels for consolidation
        """
        path_params = {"label": label, "currency": currency, "to": to}
        query = {"uniqueAddressLabels": unique_address_labels}
        return self._http.call("POST", "/api/v3/merchant/wallets/:label/:currency/consolidation/:to", path_params=path_params, query=query, authed=True)

    def post_merchant_wallets_consolidation_by_to_v3(self, to: Any, body: dict) -> Any:
        """Execute merchant wallets consolidation (sending funds to the main wallet) - Supports Auth methods: [ OAuth, clientId/secret ]

        See: https://docs.coinpayments.net/api/transactions/routes/postMerchantWalletsConsolidationByToV3

        Path:
            to: the ID of the wallet which will be receiving funds (optional)

        Body: WalletsConsolidationRequestV3Dto
        """
        path_params = {"to": to}
        return self._http.call("POST", "/api/v3/merchant/wallets/consolidation/:to", path_params=path_params, body=body, authed=True)

    def post_merchant_wallets_consolidation_preview_v3(self, body: dict) -> Any:
        """Preview merchant wallets consolidation (sending funds to the main wallet) - Supports Auth methods: [ OAuth, clientId/secret ]

        See: https://docs.coinpayments.net/api/transactions/routes/postMerchantWalletsConsolidationPreviewV3

        Body: WalletsConsolidationRequestV3Dto
        """
        return self._http.call("POST", "/api/v3/merchant/wallets/consolidation-preview", body=body, authed=True)

    def get_merchant_wallets_consolidation_transactions_by_label_currency_id_v3(self, label: Any, currency: Any, id: Any) -> Any:
        """Lists transactions used for consolidating funds to pay fee for the withdrawal - Supports Auth methods: [ OAuth, clientId/secret ]

        See: https://docs.coinpayments.net/api/transactions/routes/getMerchantWalletsConsolidationTransactionsByLabelCurrencyIdV3

        Path:
            label: the unique label for the wallet provided by client
            currency: the currency of the wallet to list consolidation transactions for
            id: ID of the withdrawal (spend request)
        """
        path_params = {"label": label, "currency": currency, "id": id}
        return self._http.call("GET", "/api/v3/merchant/wallets/:label/:currency/consolidation-transactions/:id", path_params=path_params, authed=True)
