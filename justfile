run CONFIGURATION:
  CGO_ENABLED=1 go run . simulation {{ CONFIGURATION }}

clean:
  mkdir -p ./result
  mkdir -p ./result/old
  mv ./result/*.sql* ./result/old || true

graph SQL:
  python3 ./processing/main.py "{{ SQL }}"

all CONFIGURATION:
  just run {{ CONFIGURATION }}
  just graph ./result/*.sqlite

