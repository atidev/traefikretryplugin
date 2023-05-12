# Traefik Retry Plugin

Plugin that allows to retry HTTP requests based on policy specified in the header.

## The Header

Plugin uses structured header named `Retry-Policy`.
That header implements [RFC-8941 dictionary](https://www.rfc-editor.org/rfc/rfc8941.html#name-dictionaries) structure.

Header fields description:

| Key        | Value                                                                                                                                                                                                |
|------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `codes`    | Interval of response codes in mathematical notation of intervals, with spaces used as a separator, eg: <br/>`[502 504]` — 502 <= codes >= 504<br/>`[502 504) 429` — 502 <= codes > 504, codes == 429 |
| `attempts` | A number with a self-explanatory name, eg: `3`                                                                                                                                                       |
