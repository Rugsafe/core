
# RugSafe Protocol - Core

![RugSafe Logo](https://rugsafe.io/_next/static/media/logo5.7217ba98.png)

RugSafe is a multichain protocol designed to transform the DeFi landscape by addressing rug-pull risks and introducing innovative financial instruments. Built on the Cosmos SDK, the RugSafe blockchain enables secure, decentralized recovery mechanisms, dynamic intents, and advanced asset management capabilities.

| Status Type          | Status                                                                 |
|----------------------|-------------------------------------------------------------------------|
| **Development Build**| [![Development Build](https://github.com/rugsafe/core/actions/workflows/pipeline.yml/badge.svg)](https://github.com/rugsafe/core/actions/workflows/pipeline.yml) |
| **Issues**           | [![Issues](https://img.shields.io/github/issues/rugsafe/core.svg)](https://github.com/rugsafe/core/issues) |
| **Last Commit**      | [![Last commit](https://img.shields.io/github/last-commit/rugsafe/core.svg)](https://github.com/rugsafe/core/commits/main) |
| **License**          | [![License](https://img.shields.io/badge/license-APACHE-blue.svg)](https://github.com/rugsafe/core/blob/main/LICENSE) |

## Protocol Overview

The RugSafe blockchain integrates recovery mechanisms and financial instruments to empower DeFi users:

1. **Vault Mechanism**: Secure deposits of rugged tokens and issuance of anti-coins, inversely pegged to token value.
2. **Will Module**: Automates asset actions such as transfers, claims, and contract calls based on defined will intents.
3. **IBC Integration**: Supports cross-chain communication and asset interoperability.

## Key Features

### Will Module
- **Dynamic Automation**: Define will-based actions such as scheduled transfers, claims, or contract executions.
- **Secure Execution**: Enforces permissions and ensures compliance with user-defined access controls.
- **Customizable Conditions**: Automate actions based on custom conditions

### Rug Detection Mechanisms
- **Liquidity Monitoring**: Detects sudden liquidity shifts indicative of rug-pull risks.
- **Proactive Interventions**: Executes protective actions like freezing or front-running suspicious transactions.
- **Mitigation Systems**: Enables automated user intents for responding to risky scenarios.

### IBC and Interoperability
- **Cross-Chain Communication**: Supports secure asset transfers and interactions with other Cosmos SDK chains.
- **Interchain Intents**: Automate asset management across multiple chains using IBC.

## This Repository

This repository contains the source code for the RugSafe Cosmos SDK blockchain. The primary module implemented is the **Will Module**, which handles:

1. **Will Creation**: Allows users to define intents such as transfers, claims, and contract executions.
2. **Will Execution**: Executes defined actions based on user-defined conditions like block height or time.
3. **Claims Processing**: Verifies and processes claims through mechanisms like Schnorr signatures and Pedersen commitments.

### Folder Structure
```bash
.
├── Makefile
├── app/
├── cmd/
├── docs/
├── modules/
│   └── will/
├── scripts/
├── testnet/
└── README.md
```

## Privacy Features

The RugSafe protocol integrates advanced cryptographic techniques to ensure the confidentiality and integrity of claims processing. The **Will Module** leverages three primary methods for securing claims:

### Schnorr Signatures
- **Overview**: Schnorr signatures provide a robust mechanism for validating claims by ensuring that only authorized entities can sign and submit claims.
- **Benefits**:
  - Lightweight and efficient for validation.
  - Offers provable security under the discrete logarithm assumption.
- **Implementation**:
  - Uses the `edwards25519` curve for cryptographic operations.
  - Verifies claims by reconstructing the signature using the public key, signature components, and the hashed message.

### Pedersen Commitments
- **Overview**: Pedersen commitments allow users to securely commit to a value while keeping it hidden, enabling zero-knowledge verification.
- **Benefits**:
  - Perfectly hiding and computationally binding properties.
  - Protects the confidentiality of sensitive data during claims.
- **Implementation**:
  - Claims are validated by checking that the sum of commitments matches a target commitment stored in the will.
  - Utilizes the `ristretto` curve for efficient point operations and serialization.

### zkSNARKs (Zero-Knowledge Succinct Non-Interactive Arguments of Knowledge)
- **Overview**: zkSNARKs provide a way to prove the validity of a claim without revealing any underlying information.
- **Benefits**:
  - Enables full privacy-preserving claims.
  - Compact proofs with fast verification times.
- **Implementation**:
  - Allows claimants to generate and submit proofs based on predefined zkSNARK circuits.
  - Public inputs (if any) are securely provided to the chain for verification without compromising privacy.

### Unified Approach
- Claims processing is designed to support multiple cryptographic schemes for flexibility and enhanced security.
- The protocol automatically selects the appropriate mechanism based on the type of claim and its associated conditions.
- All cryptographic operations are performed on-chain to maintain trustlessness and decentralization.

**Note**: These privacy-preserving methods are key components of RugSafe's commitment to secure and confidential DeFi recovery mechanisms.


## Quick Start

### Prerequisites

Ensure you have the following tools installed:
- **Go**: [Install Go](https://golang.org/doc/install)
- **Cosmos SDK**: [Cosmos SDK Documentation](https://docs.cosmos.network/)
- **Node.js**: For frontend integration (if applicable).

### Build and Run

Build the blockchain:
```bash
$ make build
```

Run using the `run.sh` script:
```bash
$ bash run.sh
```

### Testing

Run unit tests:
```bash
$ make test
```




## Contributing

We welcome contributions to RugSafe! Join our community and help shape the future of DeFi:
- **Discord**: [Join our community](https://discord.gg/ecMQ2D6nsu)
- **Telegram**: [Join the chat](https://t.me/rugsafe)

## License

RugSafe is released under the [APACHE License](LICENSE).

---

**Note**: This repository is under active development and may undergo significant changes. For a detailed understanding of RugSafe, refer to our [white paper](https://rugsafe.io/assets/paper/rugsafe.pdf).
