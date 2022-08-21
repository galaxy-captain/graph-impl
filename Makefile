.SILENT: try

clean:
	rm -rf ./.build/*

compile: clean
	go build -o ./.build ./cmd/...
	cp -r ./configs ./.build/configs

run: compile
	@printf "\n============================Run============================\n\n"
	cd ./.build && ./server

try: compile
	@printf "\n============================Run============================\n"
	cd ./.build && ENV=live SRV_SCENARIO=s1 SRV_TYPE=prd PORT=8888 ./server
