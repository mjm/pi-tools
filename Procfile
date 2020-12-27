homebase-web: cd homebase && yarn dev
homebase-bot: ibazel run //homebase/cmd/homebase-bot-srv -- -debug
deploy-srv: ibazel run //deploy/cmd/deploy-srv -- -debug -dry-run --github-token-path=$HOME/.github_auth.txt -poll-interval=15s
detect-presence: ibazel run //detect-presence/cmd/detect-presence-srv -- -debug -device-file $PWD/dev-devices.json -ping-interval 5s
go-links: ibazel run //go-links/cmd/go-links -- -debug
jaeger: /usr/local/bin/jaeger-all-in-one
tests: ibazel test //detect-presence/... //go-links/...
