

.FORCE:

dist: .FORCE clean
	go build
	mkdir -p build
	cp ./gl4-game ./build/
	cp -R ./resources ./build/

	mkdir -p dist
	cd build && tar czf ../dist/gl4-game.tgz .

clean:
	rm -rf ./dist ./build ./gl4-game.tgz