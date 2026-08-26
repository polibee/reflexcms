"""Auto-generated. Do not edit by hand."""
from typing import Optional
from ._http import HttpClient
from .access_control import AccessControlApi
from .currencies import CurrenciesApi
from .fees import FeesApi
from .invoices import InvoicesApi
from .rates import RatesApi
from .transactions import TransactionsApi
from .wallets import WalletsApi
from .webhooks import WebhooksApi


class CoinPaymentsClient:
    """High-level client. Instantiate once, reuse everywhere."""

    def __init__(
        self,
        client_id: Optional[str] = None,
        client_secret: Optional[str] = None,
        base_url: str = "https://a-api.coinpayments.net",
    ):
        self._http = HttpClient(base_url, client_id, client_secret)
        self.access_control = AccessControlApi(self._http)
        self.currencies = CurrenciesApi(self._http)
        self.fees = FeesApi(self._http)
        self.invoices = InvoicesApi(self._http)
        self.rates = RatesApi(self._http)
        self.transactions = TransactionsApi(self._http)
        self.wallets = WalletsApi(self._http)
        self.webhooks = WebhooksApi(self._http)
