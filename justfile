GIT := if path_exists(home_directory() + "/.ssh/id_mcd_ed25519") == "true" {
  "GIT_SSH_COMMAND='ssh -i " + home_directory() + "/.ssh/id_mcd_ed25519 -o IdentitiesOnly=yes' git" 
} else { 
  "git" 
}


test:
  go test --count=1 ./tests

run:
  go run ./cmd/burner

run-seeded SEED:
  go run ./cmd/burner --seed={{ SEED }}

push MESSAGE:
  {{ GIT }} add .
  {{ GIT }} commit --allow-empty -m "{{ MESSAGE }}"
  {{ GIT }} push origin main

chain:
  just run
  dot -Tsvg chain.dot -o chain.svg
  open chain.svg

chain-seeded SEED:
  just run-seeded {{ SEED }}
  dot -Tsvg chain.dot -o chain.svg
  open chain.svg
