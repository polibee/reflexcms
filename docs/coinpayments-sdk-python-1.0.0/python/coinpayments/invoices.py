"""Auto-generated. Do not edit by hand."""
from typing import Any, Optional
from ._http import HttpClient


class InvoicesApi:
    def __init__(self, http: HttpClient):
        self._http = http

    def post_merchant_invoices_v1(self, body: dict) -> Any:
        """Creates a new invoice - Supports Auth methods: [ OAuth, clientId/secret ]

        See: https://docs.coinpayments.net/api/invoices/routes/postMerchantInvoicesV1

        Body: CreateInvoicePaymentWithCurrencyRequestDtoCreateMerchantInvoiceRequestDto
        """
        return self._http.call("POST", "/api/v1/merchant/invoices", body=body, authed=True)

    def post_merchant_invoices_v2(self, body: dict) -> Any:
        """Creates a new invoice - Supports Auth methods: [ OAuth, clientId/secret ]

        See: https://docs.coinpayments.net/api/invoices/routes/postMerchantInvoicesV2

        Body: CreateMerchantInvoiceRequestV2Dto
        """
        return self._http.call("POST", "/api/v2/merchant/invoices", body=body, authed=True)

    def post_merchant_invoices_cancel_by_id_v1(self, id: Any) -> Any:
        """Cancels an invoice. - Supports Auth methods: [ OAuth, clientId/secret ]

        See: https://docs.coinpayments.net/api/invoices/routes/postMerchantInvoicesCancelByIdV1

        Path:
            id: (no description)
        """
        path_params = {"id": id}
        return self._http.call("POST", "/api/v1/merchant/invoices/:id/cancel", path_params=path_params, authed=True)

    def get_merchant_invoices_v1(self, *, status: Optional[Any] = None, from_: Optional[Any] = None, to: Optional[Any] = None, q: Optional[Any] = None, integration: Optional[Any] = None, payout_wallet_id: Optional[Any] = None, after: Optional[Any] = None, limit: Optional[Any] = None) -> Any:
        """Get a list of the current merchant's invoices - Supports Auth methods: [ OAuth, clientId/secret ]

        See: https://docs.coinpayments.net/api/invoices/routes/getMerchantInvoicesV1

        Args:
            status: optional query to fetch invoices that were created with the specific client
            from_: optional query to fetch from and including the time specified up to the current time
            to: optional query to fetch all invoices up to and including the specified time
            q: optional search string to find invoices with these words
            integration: optional integration by which the invoice was created
            payout_wallet_id: optional query to filter the invoices by the wallet they were paid out to (for 'paid' and 'completed' invoices)
            after: cursor that points to the end of the page of data that has been returned
            limit: the maximum number of objects that may be returned, the query may return fewer results than the requested maximum.  If `after` is specified then `limit` specifies the number of items to return starting from `after`.  If not specified then `limit` specifies the number of items to return from the beginning.
        """
        query = {"status": status, "from": from_, "to": to, "q": q, "integration": integration, "payoutWalletId": payout_wallet_id, "after": after, "limit": limit}
        return self._http.call("GET", "/api/v1/merchant/invoices", query=query, authed=True)

    def get_merchant_invoices_v2(self, *, status: Optional[Any] = None, from_: Optional[Any] = None, to: Optional[Any] = None, q: Optional[Any] = None, integration: Optional[Any] = None, payout_wallet_id: Optional[Any] = None, after: Optional[Any] = None, limit: Optional[Any] = None) -> Any:
        """Get a list of the current merchant's invoices - Supports Auth methods: [ OAuth, clientId/secret ]

        See: https://docs.coinpayments.net/api/invoices/routes/getMerchantInvoicesV2

        Args:
            status: optional query to fetch invoices that were created with the specific client
            from_: optional query to fetch from and including the time specified up to the current time
            to: optional query to fetch all invoices up to and including the specified time
            q: optional search string to find invoices with these words
            integration: optional integration by which the invoice was created
            payout_wallet_id: optional query to filter the invoices by the wallet they were paid out to (for 'paid' and 'completed' invoices)
            after: cursor that points to the end of the page of data that has been returned
            limit: the maximum number of objects that may be returned, the query may return fewer results than the requested maximum.  If `after` is specified then `limit` specifies the number of items to return starting from `after`.  If not specified then `limit` specifies the number of items to return from the beginning.
        """
        query = {"status": status, "from": from_, "to": to, "q": q, "integration": integration, "payoutWalletId": payout_wallet_id, "after": after, "limit": limit}
        return self._http.call("GET", "/api/v2/merchant/invoices", query=query, authed=True)

    def post_merchant_invoices_buy_now_button_v1(self, body: dict) -> Any:
        """Creates payment button code for the specified invoice - Supports Auth methods: [ OAuth, clientId/secret ]

        See: https://docs.coinpayments.net/api/invoices/routes/postMerchantInvoicesBuyNowButtonV1

        Body: CreateMerchantInvoiceBuyNowButtonHtmlRequestDto
        """
        return self._http.call("POST", "/api/v1/merchant/invoices/buy-now-button", body=body, authed=True)

    def post_merchant_invoices_buy_now_button_v2(self, body: dict) -> Any:
        """Creates payment button code for the specified invoice - Supports Auth methods: [ OAuth, clientId/secret ]

        See: https://docs.coinpayments.net/api/invoices/routes/postMerchantInvoicesBuyNowButtonV2

        Body: CreateMerchantInvoiceBuyNowButtonHtmlRequestV2Dto
        """
        return self._http.call("POST", "/api/v2/merchant/invoices/buy-now-button", body=body, authed=True)

    def get_invoices_payment_currencies_by_id_currency_v1(self, id: Any, currency: Any) -> Any:
        """Provides an object with details for sending payments to the invoice - Supports Auth methods: Anonymous

        See: https://docs.coinpayments.net/api/invoices/routes/getInvoicesPaymentCurrenciesByIdCurrencyV1

        Path:
            id: If of active invoice
            currency: Currency in format 'id' or 'id:contractAddress' for smart contracts
        """
        path_params = {"id": id, "currency": currency}
        return self._http.call("GET", "/api/v1/invoices/:id/payment-currencies/:currency", path_params=path_params, authed=False)

    def get_invoices_payment_currencies_status_by_id_currency_v1(self, id: Any, currency: Any) -> Any:
        """Returns the current state of the payment object - Supports Auth methods: Anonymous

        See: https://docs.coinpayments.net/api/invoices/routes/getInvoicesPaymentCurrenciesStatusByIdCurrencyV1

        Path:
            id: If of active invoice
            currency: Currency in format 'id' or 'id:contractAddress' for smart contracts
        """
        path_params = {"id": id, "currency": currency}
        return self._http.call("GET", "/api/v1/invoices/:id/payment-currencies/:currency/status", path_params=path_params, authed=False)

    def get_merchant_invoices_by_id_v1(self, id: Any, *, include_full_details: Optional[Any] = None) -> Any:
        """Find invoice belonging to merchant by the invoice ID - Supports Auth methods: [ OAuth, clientId/secret ]

        See: https://docs.coinpayments.net/api/invoices/routes/getMerchantInvoicesByIdV1

        Path:
            id: Invoice id

        Args:
            include_full_details: Indicates whether to return information about Merchant, Metadata, Notes to recipient and email delivery settings
        """
        path_params = {"id": id}
        query = {"include_full_details": include_full_details}
        return self._http.call("GET", "/api/v1/merchant/invoices/:id", path_params=path_params, query=query, authed=True)

    def get_merchant_invoices_by_id_v2(self, id: Any, *, include_full_details: Optional[Any] = None) -> Any:
        """Find invoice belonging to merchant by the invoice ID - Supports Auth methods: [ OAuth, clientId/secret ]

        See: https://docs.coinpayments.net/api/invoices/routes/getMerchantInvoicesByIdV2

        Path:
            id: Invoice Id

        Args:
            include_full_details: Indicates whether to return information about Merchant, Metadata, Notes to recipient and email delivery settings
        """
        path_params = {"id": id}
        query = {"include_full_details": include_full_details}
        return self._http.call("GET", "/api/v2/merchant/invoices/:id", path_params=path_params, query=query, authed=True)

    def get_merchant_invoices_payouts_by_id_v1(self, id: Any) -> Any:
        """Get payout details for an invoice, including if invoice has been fully paid out, the exact amount they will receive
        and in what currency, which address payout will be deposited, and who (Buyer) performed the payment. - Supports Auth methods: [ OAuth, clientId/secret ]

        See: https://docs.coinpayments.net/api/invoices/routes/getMerchantInvoicesPayoutsByIdV1

        Path:
            id: Invoice Id
        """
        path_params = {"id": id}
        return self._http.call("GET", "/api/v1/merchant/invoices/:id/payouts", path_params=path_params, authed=True)

    def get_merchant_invoices_payouts_by_id_v2(self, id: Any) -> Any:
        """Get payout details for an invoice, including if invoice has been fully paid out, the exact amount they will receive
        and in what currency, which address payout will be deposited, and who (Buyer) performed the payment. - Supports Auth methods: [ OAuth, clientId/secret ]

        See: https://docs.coinpayments.net/api/invoices/routes/getMerchantInvoicesPayoutsByIdV2

        Path:
            id: Invoice id
        """
        path_params = {"id": id}
        return self._http.call("GET", "/api/v2/merchant/invoices/:id/payouts", path_params=path_params, authed=True)

    def get_merchant_invoices_history_by_id_v1(self, id: Any) -> Any:
        """Lists the history events of an invoice by the invoice ID - Supports Auth methods: [ OAuth, clientId/secret ]

        See: https://docs.coinpayments.net/api/invoices/routes/getMerchantInvoicesHistoryByIdV1

        Path:
            id: Invoice id
        """
        path_params = {"id": id}
        return self._http.call("GET", "/api/v1/merchant/invoices/:id/history", path_params=path_params, authed=True)

    def get_merchant_invoices_history_by_id_v2(self, id: Any) -> Any:
        """Lists the history events of an invoice by the invoice ID - Supports Auth methods: [ OAuth, clientId/secret ]

        See: https://docs.coinpayments.net/api/invoices/routes/getMerchantInvoicesHistoryByIdV2

        Path:
            id: Invoice Id
        """
        path_params = {"id": id}
        return self._http.call("GET", "/api/v2/merchant/invoices/:id/history", path_params=path_params, authed=True)
