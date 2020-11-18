homebase-web: cd homebase && yarn dev
homebase-bot: ibazel run //homebase/cmd/homebase-bot-srv -- -debug
detect-presence: ibazel run //detect-presence/cmd/detect-presence-srv -- -debug -device-file $PWD/dev-devices.json -ping-interval 5s
go-links: ibazel run //go-links/cmd/go-links -- -debug
