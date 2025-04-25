# Integration types

- Required features:

  - Adjustable difficulty:

    - PoW: Discrete Steps every 2016 blocks
    - PoB:
      - Discrete Steps every N blocks
      - Price Adjustment every M blocks

  - Node Power
  - Node Efficiency

- Simulations:

  - Baseline (done):

    - Proof of Work (periodic difficulty adjustment)

  - Simple Integration (done):

    - Proof of Work (periodic difficulty adjustment)
    - Proof of Burn (per block difficulty adjustment) (must be after PoW)
    - Expectation:
      - If PoB cannot produce a block, it continues to mine PoW
      - Transactions per unit time increase
      - Power per transaction decrease
      - Power per unit time unchanged

  - Uncontrolled Burn:

    - Proof of Work (periodic difficulty adjustment)
    - Proof of Burn (very low difficulty)
    - Expectation:
      - Very large amount of blocks being mined
      - Very little transactions per block
      - Energy consumption similar to PoW

  - Pure Burn:

    - Proof of Burn (per block difficulty adjustment)
    - Expectation:
      - Getting stuck!

  - Adjustable Integration:

    - Proof of Work (discrete difficulty adjustment)
    - Proof of Burn (continuous difficulty adjustment)
      - NOTE: The continuous difficulty adjustment is a difficult task. Here we
        will assume it is possible to implement and that it depends on some time
        T that nodes across the network can agree on.
      - NOTE: There could potentially be a need for us to adjust baseline
        difficulty from which we start decaying, this could be done similarly to
        PoW and adjusted every 2 weeks or something.
      - NOTE: The true difficulty adjustment in this scenario would be done with
        price adjustment - if we do not manage (even with continuous reduction)
        to mine a block within a specified time (e.g. 5 minutes, it should be
        lower than PoW interval), we decrease the price for burning or increase
        the power per burned currency.

  - Adjustable Burn:
    - Proof of Burn (continuous difficulty adjustment)
    - Expectation:
      - Same as Adjustable Integration, but because we do not have PoW in the
        same network, we can adjust out mine on decay time as we see fit.
