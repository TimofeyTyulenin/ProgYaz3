package main

import (
	"fmt"
	"math/rand"
	"time"
)

var ring []*Token

type Token struct {
	id          int
	left_token  <-chan Message
	right_token chan Message
}

type Message struct {
	data      string
	recipient int
	ttl       int
}

func NewMessage(data string, recipient, ttl int) Message {
	return Message{
		data:      data,
		recipient: recipient - 1,
		ttl:       ttl,
	}
}

func initialize(n int) []*Token {
	ring := make([]*Token, 0, n)

	ring = append(ring, &Token{id: 0, right_token: make(chan Message)})

	for i := 1; i < n; i++ {
		ring = append(ring, &Token{id: i, left_token: ring[i-1].right_token, right_token: make(chan Message)})
	}

	ring[0].left_token = ring[n-1].right_token

	return ring

}

func Send(n, num int, m string) {
	rand.Seed(time.Now().UnixNano())
	msg := NewMessage(m, num, rand.Intn(n))
	for i := 0; i < n; i++ {
		t := ring[i]

		go Start(*t)
		if i == 0 {
			ring[len(ring)-1].right_token <- msg
		}
		time.Sleep(time.Millisecond * 10)
	}

}

func Start(t Token) {

	msg_fake := <-t.left_token

	switch {
	case msg_fake.recipient == t.id:

		fmt.Println("ID", t.id+1, "Доставлено ( Получатель:", msg_fake.recipient+1, "Сообщение:", msg_fake.data, " Ttl:", msg_fake.ttl, ")")

		return
	case msg_fake.ttl > 0:

		msg_fake.ttl -= 1
		t.right_token <- msg_fake

	default:

		fmt.Println("на ID", t.id+1, "Ttl истек: Не доставлено ( Получатель:", msg_fake.recipient+1, "Сообщение:", msg_fake.data, " Ttl:", msg_fake.ttl, ")")
		return
	}

	time.Sleep(time.Millisecond * 500)

}

var n int
var num int
var m string

func main() {
	fmt.Println("Введите количество N")
	fmt.Scanf("%d\n", &n)
	ring = initialize(n)
	fmt.Println("Введите номер получателя и сообщение")
	fmt.Scanf("%d\n%s", &num, &m)
	if n < num {
		fmt.Println("Такого получателя не существует. Сообщение не доставлено")
		return
	}
	Send(n, num, m)

}
