.PHONY: test
test:
	cd test
	mkdir -p dist
	tinygo build -o dist/testmain.wasm -target wasm ./...
	cat dist/testmain.wasm | deno run https://denopkg.com/syumai/binpack/mod.ts > dist/testmainwasm.ts && rm dist/testmain.wasm

