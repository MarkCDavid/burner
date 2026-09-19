# burner

An event-driven blockchain consensus simulator. It models a network of
nodes that reach consensus with Proof of Work, SlimCoin-style Proof of
Burn, or a price-controlled Proof of Burn (PPoB), and records every
block, fork, and consensus adjustment into a SQLite database.

Built for the master's thesis
[*Evaluating the Interoperability of Proof of Burn and other Consensus
Algorithms Through Simulations*](https://gs.elaba.lt/object/elaba:239097808/)
(A. Šakalys, Vilnius Gediminas Technical University, 2025).

## Build and run

Requires go with CGO enabled (sqlite3). `flake.nix` provides a dev
shell with everything, including the python packages for the
processing scripts.

```
CGO_ENABLED=1 go build -o burner .
mkdir -p result
./burner simulation configuration/bitcoin.yaml --seed 1001
```

Each scenario is a YAML file in `configuration/`. The seed comes from
the file; `--seed` overrides it. Runs are deterministic per seed. One
simulated year finishes in under a minute on a desktop machine.

Results land in `result/<scenario> seed=<seed> (<timestamp>).sqlite`.
The plotting scripts in `processing/` read these databases
(`python3 processing/p_production_times.py <db>`, etc.).

## Dataset

The release
[`data-v1`](https://github.com/MarkCDavid/burner/releases/tag/data-v1)
holds the raw simulation output: 21 databases (6 scenarios × 3 seeds,
plus the PoB-without-PoW deadlock test), zstd-compressed, 520 MB total
(1.7 GB unpacked).

Unpack and inspect:

```
zstd -d bitcoin_seed1001.sqlite.zst
sqlite3 bitcoin_seed1001.sqlite ".tables"
```

Verify a download with the `SHA256SUMS` file from the release:

```
sha256sum -c SHA256SUMS
```

Main tables:

- `blocks` — one row per block, abandoned candidates included.
  `blockType`: 0 = Proof of Work, 1 = Proof of Burn. Times are in
  simulated seconds.
- `nodes` — per-node power draws (watts).
- `proof_of_work_consensus` — difficulty over time.
- `pricing_proof_of_burn_burn_consensus` — the PPoB burn price over time.
- `pricing_proof_of_burn_burn_transaction` — every burn transaction.
- `label` — the scenario name.

The `slimcoin_only_*` databases contain no blocks: that configuration
halts at t=0, because SlimCoin-style PoB can only mint after a
received PoW block.