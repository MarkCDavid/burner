test:
  go test ./tests

run:
  go run ./cmd/burner

push MESSAGE:
  git add -i
  git commit -m "{{ MESSAGE }}"
  git push origin main
