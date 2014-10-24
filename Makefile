
all:
	protoc --java_out=protoc wirebulletin.proto
	protoc --go_out=ahimsa/wirebulletin wirebulletin.proto
	protoc --python_out=protoc wirebulletin.proto

clean:
	rm -rf protoc/*
