gen_grpc:
	cd protos && \
	protoc \
	--go_out=./workshopParser \
	--go_opt=paths=source_relative \
	--go-grpc_out=./workshopParser \
	--go-grpc_opt=paths=source_relative \
	workshopParser.proto