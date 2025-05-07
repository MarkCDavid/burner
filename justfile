run CONFIGURATION:
  mkdir -p ./result
  mkdir -p ./result/old
  mv ./result/*.sql* ./result/old || true
  CGO_ENABLED=1 go run . simulation {{ CONFIGURATION }}

graph +SQL:
  for sql in {{SQL}}; do python3 ./processing/main.py "$sql"; done

all CONFIGURATION:
  just run {{ CONFIGURATION }}
  just graph ./result/*.sqlite

ppob:
  just run simulation
  just graph ./result/PricingOnly_*.sqlite
