
clean:
	rm -rf ./.build/*

compile: clean
	go build -o ./.build ./cmd/...
	cp -r ./configs ./.build/configs

run: compile
	@printf "\n============================Run============================\n\n"
	cd ./.build && ./service
