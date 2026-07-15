# Design a Call Center

Design a call center with three employee levels — `"respondent"`,
`"manager"`, and `"director"` — created with a fixed count of each.

An incoming call goes to a free **respondent**; if none is free it
escalates to a free **manager**, then a free **director**. If everyone
is busy the call waits in a FIFO queue. When any call ends, the freed
employee immediately takes the longest-waiting queued call.

- `dispatch(callID)` — route a new call; return the level that took it
  (`"respondent"`, `"manager"`, `"director"`) or `"queued"`
- `end_call(callID)` — finish an active call (freeing its employee) or
  abandon a queued one; return whether the call existed
- `handler_of(callID)` — the level currently handling the call,
  `"queued"` if waiting, or `""` if unknown/ended

## Examples

```
c = new CallCenter(1 respondent, 1 manager, 0 directors)
dispatch(1) -> "respondent"
dispatch(2) -> "manager"
dispatch(3) -> "queued"
end_call(1)
handler_of(3) -> "respondent"
```

## Constraints

- Call IDs are positive and unique per dispatch.
- Queued calls are assigned strictly in arrival order.
