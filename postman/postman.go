package postman

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func postman(ctx context.Context, wg *sync.WaitGroup, transferPoint chan<- string, n int, mail string) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Я почтальон номер:", n, "Мой рабочий день окончен!")
			return
		default:
			fmt.Println("Я почтальон номер:", n, "Взял письмо...")
			time.Sleep(1 * time.Second)
			fmt.Println("Я почтальон номер:", n, "Донёс письмо до почты:", mail)

			transferPoint <- mail

			fmt.Println("Я почтальон номер:", n, "Передал письмо:", mail)
		}
	}
}

func PostmanPool(ctx context.Context, postmanCount int) <-chan string {
	mailTransferPoint := make(chan string)

	wg := &sync.WaitGroup{}

	for i := 1; i <= postmanCount; i++ {
		wg.Add(1)
		go postman(ctx, wg, mailTransferPoint, i, postmanToMail(i))
	}

	go func() {
		wg.Wait()
		close(mailTransferPoint)
	}()

	return mailTransferPoint
}

func postmanToMail(postmanNumber int) string {
	ptm := map[int]string{
		1: "Семейный привет",
		2: "Приглашение от друга",
		3: "Информация из автосервиса",
	}

	mail, ok := ptm[postmanNumber]
	if !ok {
		return "Лотерея"
	}

	return mail
}
