package main

import (
	"context"
	"fmt"
	"strings"
	"unicode"
)

// информация о количестве цифр в каждом слове
type counter map[string]int

// слово и количество цифр в нем
type pair struct {
	word  string
	count int
}

// начало решения

// считает количество цифр в словах
func countDigitsInWords(ctx context.Context, words []string) counter {
	childCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	select {
	case <-ctx.Done():
		cancel()
		return counter{}
	default:
		pending := submitWords(childCtx, words)
		counted := countWords(childCtx, pending)
		return fillStats(childCtx, counted)
	}
	
}

// отправляет слова на подсчет
func submitWords(ctx context.Context, words []string) <-chan string {
	out := make(chan string)
	go func() {
		for _, word := range words {
			select {
			case <- ctx.Done():
				close(out)
				return
			case out <- word:
			}
		}
		close(out)
	}()
	return out
}

// считает цифры в словах
func countWords(ctx context.Context, in <-chan string) <-chan pair {
	out := make(chan pair)
	go func() {
		for word := range in {
			count := countDigits(word)
			select {
			case <- ctx.Done():
				close(out)
				return
			case out <- pair{word, count}:
			}
		}
		close(out)
	}()
	return out
}

// готовит итоговую статистику
func fillStats(ctx context.Context, in <-chan pair) counter {
	stats := counter{}
	for p := range in {
		select {
		case <-ctx.Done():
			return counter{}
		default:
			stats[p.word] = p.count
		}
	}
	return stats
}

// конец решения

// считает количество цифр в слове
func countDigits(str string) int {
	count := 0
	for _, char := range str {
		if unicode.IsDigit(char) {
			count++
		}
	}
	return count
}

func main() {
	phrase := "0ne 1wo thr33 4068"
	words := strings.Fields(phrase)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stats := countDigitsInWords(ctx, words)
	fmt.Println(stats)
}
