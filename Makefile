
all:
	protoc --java_out=protoc wirebulletin.proto
	protoc --go_out=protoc wirebulletin.proto
	protoc --python_out=protoc wirebulletin.proto

clean:
	rm -f protoc/*
