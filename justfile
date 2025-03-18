test:
  go test --count=1 ./tests

run:
  go run ./cmd/burner

run-seeded SEED:
  go run ./cmd/burner --seed={{ SEED }}

push MESSAGE:
  git add .
  git commit --allow-empty -m "{{ MESSAGE }}"
  git push origin main

chain:
  just run
  dot -Tsvg chain.dot -o chain.svg
  open chain.svg

chain-seeded SEED:
  just run-seeded {{ SEED }}
  dot -Tsvg chain.dot -o chain.svg
  open chain.svg
