"""Auto-generated. Do not edit by hand."""
from typing import Any, Optional
from ._http import HttpClient


class RatesApi:
    def __init__(self, http: HttpClient):
        self._http = http

    def get_rates_v1(self, *, from_: Optional[Any] = None, to: Optional[Any] = None, point_in_time: Optional[Any] = None) -> Any:
        """lists the current conversion rates between currencies - Supports Auth methods: Anonymous

        See: https://docs.coinpayments.net/api/rates/routes/getRatesV1

        Args:
            from_: comma separated list of currencies to use as the source for rate calculations. The allowed input format is {currencyId}:{contractAddress}
            to: comma separated list of currencies for which to retrieve conversion rates for (from the `from` currencies). The allowed input format is {currencyId}:{contractAddress}
            point_in_time: Point in time for which to get rate, if null the latest rate is returned
        """
        query = {"from": from_, "to": to, "pointInTime": point_in_time}
        return self._http.call("GET", "/api/v1/rates", query=query, authed=False)

    def get_rates_v2(self, *, from_: Optional[Any] = None, to: Optional[Any] = None, point_in_time: Optional[Any] = None) -> Any:
        """lists the current conversion rates between currencies - Supports Auth methods: Anonymous

        See: https://docs.coinpayments.net/api/rates/routes/getRatesV2

        Args:
            from_: comma separated list of currencies to use as the source for rate calculations. The allowed input format is {currencyId}:{contractAddress}
            to: comma separated list of currencies for which to retrieve conversion rates for (from the `from` currencies). The allowed input format is {currencyId}:{contractAddress}
            point_in_time: Point in time for which to get rate, if null the latest rate is returned
        """
        query = {"from": from_, "to": to, "pointInTime": point_in_time}
        return self._http.call("GET", "/api/v2/rates", query=query, authed=False)
