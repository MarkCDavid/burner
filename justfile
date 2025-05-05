run COMMAND="":
  CGO_ENABLED=1 go run . {{ COMMAND }}

chain:
  just run simulation
  dot -Tsvg chain.dot -o chain.svg
  firefox chain.svg
