test:
  go test --count=1 ./tests

run:
  go run ./cmd/burner

push MESSAGE:
  git add -i
  git commit -m "{{ MESSAGE }}"
  git push origin main

chain:
  just run
  dot -Tpng chain.dot -o chain.png
  feh chain.png
