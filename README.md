 
# getho

---
**getho** is a low-level Ethereum **base-layer observability and debugging tool** for inspecting execution-layer behavior.
---


## What getho does

###  Transaction Inspection

* Decode raw Ethereum transactions (RLP)
* Inspect calldata and function selectors
* Display sender, recipient, value, nonce, and type

###  Gas & Fee Breakdown

* Base fee
* Priority fee (tip)
* Blob fee (EIP-4844)
* Gas used vs gas limit

###  Execution Tracing

* Opcode-level execution traces
* Track:

  * `CALL`
  * `DELEGATECALL`
  * `STATICCALL`
  * `SSTORE` / `SLOAD`
* Identify execution paths and state changes

###  Calldata & ABI Insight

* Decode function selectors
* Parse arguments
* Highlight unknown or malformed calldata

###  Execution-Layer Focus

* L1 / base-layer first
* Geth-compatible execution semantics
* No abstraction leaks

---
 
## Example usage

```bash
# Inspect a transaction
getho tx 0xTX_HASH

# Decode calldata
getho calldata 0xTX_HASH

# Gas & fee analysis
getho gas 0xTX_HASH

# Full execution trace
getho trace 0xTX_HASH

# Decode raw RLP
getho rlp decode 0xF86B...
```



## Use cases

* Debug failed or reverted transactions
* Understand unexpected gas spikes
* Inspect blob transactions (EIP-4844)
* Learn how the EVM executes contracts step-by-step
* Validate execution-client behavior

---

getho is **not** a dApp tool — it is an **execution-layer microscope**.

---
 

## Contributing

Contributions are welcome from:

* Protocol researchers
* Execution-client devs
* Infra & tooling builders

Open issues, propose features, or submit PRs.

---
