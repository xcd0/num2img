
all: a
	go run main.go < a

rand:
	./random.sh 10000000 > a
