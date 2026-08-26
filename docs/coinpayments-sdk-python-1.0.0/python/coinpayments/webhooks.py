"""Auto-generated. Do not edit by hand."""
from typing import Any, Optional
from ._http import HttpClient


class WebhooksApi:
    def __init__(self, http: HttpClient):
        self._http = http

    def put_merchant_wallets_webhook_by_id_v1(self, id: Any, body: dict) -> Any:
        """Updates the webhook configuration for a specific wallet - Supports Auth methods: [ OAuth, clientId/secret ]

        See: https://docs.coinpayments.net/api/webhooks/routes/putMerchantWalletsWebhookByIdV1

        Path:
            id: the id of the wallet to update the webhook for

        Body: UpdateWalletWebhookRequestDto
        """
        path_params = {"id": id}
        return self._http.call("PUT", "/api/v1/merchant/wallets/:id/webhook", path_params=path_params, body=body, authed=True)

    def put_merchant_wallets_webhook_by_id_v2(self, id: Any, body: dict) -> Any:
        """Updates the webhook configuration for a specific wallet - Supports Auth methods: [ OAuth, clientId/secret ]

        See: https://docs.coinpayments.net/api/webhooks/routes/putMerchantWalletsWebhookByIdV2

        Path:
            id: the id of the wallet to update the webhook for

        Body: UpdateWalletWebhookRequestDto
        """
        path_params = {"id": id}
        return self._http.call("PUT", "/api/v2/merchant/wallets/:id/webhook", path_params=path_params, body=body, authed=True)

    def put_merchant_wallets_addresses_webhook_by_id_a_id_v1(self, id: Any, a_id: Any, body: dict) -> Any:
        """Updates the webhook configuration for a specific wallet address - Supports Auth methods: [ OAuth, clientId/secret ]

        See: https://docs.coinpayments.net/api/webhooks/routes/putMerchantWalletsAddressesWebhookByIdAIdV1

        Path:
            id: the id of the wallet containing the address
            a_id: the id of the address to update the webhook for

        Body: UpdateWalletWebhookRequestDto
        """
        path_params = {"id": id, "aId": a_id}
        return self._http.call("PUT", "/api/v1/merchant/wallets/:id/addresses/:aId/webhook", path_params=path_params, body=body, authed=True)

    def put_merchant_wallets_addresses_webhook_by_id_a_id_v2(self, id: Any, a_id: Any, body: dict) -> Any:
        """Updates the webhook configuration for a specific wallet address - Supports Auth methods: [ OAuth, clientId/secret ]

        See: https://docs.coinpayments.net/api/webhooks/routes/putMerchantWalletsAddressesWebhookByIdAIdV2

        Path:
            id: the id of the wallet containing the address
            a_id: the id of the address to update the webhook for

        Body: UpdateWalletWebhookRequestDto
        """
        path_params = {"id": id, "aId": a_id}
        return self._http.call("PUT", "/api/v2/merchant/wallets/:id/addresses/:aId/webhook", path_params=path_params, body=body, authed=True)

    def post_merchant_clients_webhooks_by_id_v1(self, id: Any, body: dict) -> Any:
        """Create a new webhook for a client - Supports Auth methods: [ OAuth, clientId/secret ]

        See: https://docs.coinpayments.net/api/webhooks/routes/postMerchantClientsWebhooksByIdV1

        Path:
            id: The public ID of a client

        Body: CreateMerchantClientWebhookRequestDto
        """
        path_params = {"id": id}
        return self._http.call("POST", "/api/v1/merchant/clients/:id/webhooks", path_params=path_params, body=body, authed=True)

    def get_merchant_clients_webhooks_by_id_v1(self, id: Any, *, type: Optional[Any] = None) -> Any:
        """List all of the webhooks for a particular client - Supports Auth methods: [ OAuth, clientId/secret ]

        See: https://docs.coinpayments.net/api/webhooks/routes/getMerchantClientsWebhooksByIdV1

        Path:
            id: Target `ClientId`

        Args:
            type: Notification type. Null for all
        """
        path_params = {"id": id}
        query = {"type": type}
        return self._http.call("GET", "/api/v1/merchant/clients/:id/webhooks", path_params=path_params, query=query, authed=True)

    def put_merchant_clients_webhooks_by_id_w_id_v1(self, id: Any, w_id: Any, body: dict) -> Any:
        """Update an existing webhook for a client - Supports Auth methods: [ OAuth, clientId/secret ]

        See: https://docs.coinpayments.net/api/webhooks/routes/putMerchantClientsWebhooksByIdWIdV1

        Path:
            id: The public ID of a client
            w_id: the ID of the webhook to update

        Body: UpdateMerchantClientWebhookDto
        """
        path_params = {"id": id, "wId": w_id}
        return self._http.call("PUT", "/api/v1/merchant/clients/:id/webhooks/:wId", path_params=path_params, body=body, authed=True)

    def delete_merchant_clients_webhooks_by_id_w_id_v1(self, id: Any, w_id: Any) -> Any:
        """Delete a webhook for a client - Supports Auth methods: [ OAuth, clientId/secret ]

        See: https://docs.coinpayments.net/api/webhooks/routes/deleteMerchantClientsWebhooksByIdWIdV1

        Path:
            id: the public ID of a client
            w_id: the ID of the webhook to delete
        """
        path_params = {"id": id, "wId": w_id}
        return self._http.call("DELETE", "/api/v1/merchant/clients/:id/webhooks/:wId", path_params=path_params, authed=True)

    def put_merchant_wallets_webhook_by_label_currency_v3(self, label: Any, currency: Any, body: dict) -> Any:
        """Updates the webhook configuration for a specific wallet - Supports Auth methods: [ OAuth, clientId/secret ]

        See: https://docs.coinpayments.net/api/webhooks/routes/putMerchantWalletsWebhookByLabelCurrencyV3

        Path:
            label: the unique id for the wallet provided by client
            currency: the currency of the wallet to update webhook for

        Body: UpdateWalletWebhookRequestDto
        """
        path_params = {"label": label, "currency": currency}
        return self._http.call("PUT", "/api/v3/merchant/wallets/:label/:currency/webhook", path_params=path_params, body=body, authed=True)

    def put_merchant_wallets_addresses_webhook_by_w_label_currency_a_label_v3(self, w_label: Any, currency: Any, a_label: Any, body: dict) -> Any:
        """Updates the webhook configuration for a specific wallet address - Supports Auth methods: [ OAuth, clientId/secret ]

        See: https://docs.coinpayments.net/api/webhooks/routes/putMerchantWalletsAddressesWebhookByWLabelCurrencyALabelV3

        Path:
            w_label: the unique label for the wallet provided by client
            currency: the currency of the wallet containing the address
            a_label: the unique label for the address provided by client

        Body: UpdateWalletWebhookRequestDto
        """
        path_params = {"wLabel": w_label, "currency": currency, "aLabel": a_label}
        return self._http.call("PUT", "/api/v3/merchant/wallets/:wLabel/:currency/addresses/:aLabel/webhook", path_params=path_params, body=body, authed=True)
