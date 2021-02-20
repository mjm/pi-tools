homebase-web: cd homebase && yarn dev
homebase-api: ibazel run //homebase/cmd/homebase-api-srv -- -debug -schema-path=$PWD/homebase/schema.graphql
homebase-bot: ibazel run //homebase/cmd/homebase-bot-srv -- -debug
deploy-srv: ibazel run //deploy/cmd/deploy-srv -- -debug -dry-run -github-token-path=$HOME/.github_auth.txt -poll-interval=15s -terraform=/usr/local/bin/terraform
detect-presence: ibazel run //detect-presence/cmd/detect-presence-srv -- -debug -device-file $PWD/dev-devices.json -ping-interval 5s -github-token-path=$HOME/.github_auth.txt
go-links: ibazel run //go-links/cmd/go-links -- -debug
vault-proxy: ibazel run //vault-proxy/cmd/vault-proxy -- -debug -auth-path webauthn-debug -cookie-domain "" -static-dir=$PWD/vault-proxy/static
backup: ibazel run //backup/cmd/backup-srv -- -debug -tarsnap-keyfile $HOME/Downloads/tarsnap-raspberrypi.key
jaeger: /usr/local/bin/jaeger-all-in-one
tests: ibazel test --keep_going //detect-presence/... //go-links/... //homebase/bot/...
relay-web: cd homebase && npx relay-compiler --watch
