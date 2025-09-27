🌌 Nebryth Protocol

Nebryth is a modular, mythic-layer blockchain protocol built on the Cosmos SDK (v0.50+), designed for scalable token ecosystems, dual-realm governance, and future-proof chain architecture. It powers the native coin NBYT and anchors the TRNX/RIGR token lifecycle.

🚀 Overview

Nebryth blends cosmic identity with technical precision. It offers:

• ⚛️ ABCI++ lifecycle support for modern app chains
• 🧬 Modular governance and staking architecture
• 🔁 Dual-ecosystem token strategy (TRNX/RIGR)
• 🛠️ CLI-ready devnet setup with reproducible builds
• 🧩 Genesis funding, validator wiring, and module orchestration


🧱 Architecture

Built on Cosmos SDK v0.50+, Nebryth includes:

• `auth`, `bank`, `staking`, `gov`, `mint`, and custom modules
• ABCI++ pre-blockers and post-blockers wired via `app.go`
• Clean `go.mod` and `go.sum` for reproducible builds
• Genesis file with full NBYT allocation and validator config


🧪 Devnet Setup

# Clone and build
git clone https://github.com/nebryth/nebryth
cd nebryth
make install

# Initialize chain
nebrythd init nebryth-devnet --chain-id nebryth-1

# Add genesis accounts and validators
nebrythd add-genesis-account <address> 1000000000nbyt
nebrythd gentx <keyname> 1000000000nbyt --chain-id nebryth-1
nebrythd collect-gentxs

# Start node
nebrythd start


🪐 Token Ecosystem

Nebryth anchors the TRNX/RIGR lifecycle:

Token	Role	Realm	Notes	
NBYT	Native coin	Protocol	Used for staking, fees, governance	
TRNX	Primary token	Public	Launchpad, liquidity, access	
RIGR	Companion token	Devnet	Validator incentives, testnet utility	


📖 Lore & Identity

Nebryth is more than a chain—it’s a mythic protocol. Inspired by cosmic cycles, scientific precision, and modular clarity, it invites builders to shape the future of decentralized ecosystems.

🛡️ Governance & Staking

• Proposal lifecycle powered by `gov` module
• Validator setup via `staking` and `slashing` modules
• Customizable voting periods and thresholds


🧰 Troubleshooting

Windows quirks? Go build errors? Genesis resets?
Check the docs or open an issue. Nebryth is built for reproducibility and rapid iteration.
