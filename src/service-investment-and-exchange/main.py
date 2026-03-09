from fastapi import FastAPI, Query
from datetime import datetime, timezone
from decimal import Decimal

app = FastAPI(title="Mock Exchange Service")

MOCK_RATES = {
    ("USD", "RUB"): Decimal("92.5030"),
    ("EUR", "RUB"): Decimal("100.2000"),
    ("RUB", "USD"): Decimal("0.0108"),
    ("RUB", "EUR"): Decimal("0.0099"),
    ("USD", "EUR"): Decimal("0.9200"),
    ("EUR", "USD"): Decimal("1.0900"),
}

@app.get("/api/v1/rates")
async def get_rate(
        base: str = Query(..., description="Базовая валюта (что продаем)"),
        target: str = Query(..., description="Целевая валюта (что покупаем)")
):
    base = base.upper()
    target = target.upper()

    if base == target:
        rate = Decimal("1.0000")
    else:
        rate = MOCK_RATES.get((base, target), Decimal("42.0000"))

    return {
        "base_currency": base,
        "target_currency": target,
        "rate": str(rate),
        "timestamp": datetime.now(timezone.utc).isoformat()
    }