all: sqr.zip
	@echo "done"

sqrcpp: main.cpp
	g++ -o sqrcpp main.cpp -lgmp -lgmpxx

sqrcpp2: main2.cpp
	g++ -o sqrcpp2 main2.cpp ./bigint/BigUnsigned.cc ./bigint/BigIntegerAlgorithms.cc ./bigint/BigInteger.cc ./bigint/BigIntegerUtils.cc ./bigint/BigUnsignedInABase.cc

sqrc: main.c
	gcc -o sqrc main.c -lgmp

sqrgo: ./gomathbig/main.go
	go build -o sqrgo ./gomathbig/

sqrgo2: ./gogmp/main.go
	go build -o sqrgo2 ./gogmp/

Sqr.class: Sqr.java
	javac Sqr.java

sqr.zip: sqrc sqrcpp sqrcpp2 sqrgo sqrgo2 Sqr.class
	zip sqr.zip sqrc sqrcpp sqrcpp2 sqrgo sqrgo2 Sqr.class
