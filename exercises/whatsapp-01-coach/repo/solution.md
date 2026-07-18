# My Design: WhatsApp

## Step 1 — Use cases, constraints, estimates

<!-- Which use cases are in scope (1:1, group, delivery states)? What's
     out (calls, stories, payments)? Estimate: concurrent connections,
     messages/sec, and connections one chat server can hold -- that
     last number sizes the fleet. -->

## Step 2 — High-level design

<!-- Persistent connections (WebSocket/TCP) from clients to chat
     servers, plus a session registry mapping user -> server. Narrate
     one message end to end: sender -> server A -> registry lookup ->
     server B -> recipient. -->

## Step 3 — Core components

<!-- Go deep where this question hinges:
     - Delivery semantics: how do sent/delivered/read receipts flow
       back? What happens to a message sent to someone offline --
       where does it wait, and how is it drained on reconnect?
     - Group messaging: how does one message become N deliveries? What
       breaks first in a very large group, and how do you mitigate it?
     - Storage & encryption: are messages stored server-side, and for
       how long? What does end-to-end encryption make impossible for
       the server to do? -->

## Step 4 — Scale it

<!-- Where does this break at 10x? The session registry is now the
     critical shared state -- how does IT scale and stay available?
     How is presence/last-seen handled without hammering the registry
     on every status change? What about a user with multiple devices? -->
