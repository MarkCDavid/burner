run COMMAND="":
  mkdir -p ./result
  mkdir -p ./result/old
  mv ./result/*.sql* ./result/old || true
  CGO_ENABLED=1 go run . {{ COMMAND }}

graph +SQL:
  for sql in {{SQL}}; do python3 ./processing/main.py "$sql"; done

ppob:
  just run simulation
  just graph ./result/PricingOnly_*.sqlite
