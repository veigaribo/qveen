package templates

import (
	"fmt"
)

func TemplateAdd(xs ...int) (int, error) {
	acc := 0

	for _, x := range xs {
		acc += x
	}

	return acc, nil
}

func TemplateSub(xs ...int) (int, error) {
	if len(xs) < 2 {
		return 0, fmt.Errorf("Can't subtract less than 2 operands: %v.", xs)
	}

	acc := xs[0]

	for _, x := range xs[1:] {
		acc -= x
	}

	return acc, nil
}

func TemplateMul(xs ...int) (int, error) {
	acc := 1

	for _, x := range xs {
		acc *= x
	}

	return acc, nil
}

func TemplateDiv(xs ...int) (int, error) {
	if len(xs) < 2 {
		return 0, fmt.Errorf("Can't divide less than 2 operands: %v.", xs)
	}

	acc := xs[0]

	for _, x := range xs[1:] {
		if x == 0 {
			return 0, fmt.Errorf("Division by zero: %d / %d.", acc, x)
		}

		acc /= x
	}

	return acc, nil
}

func TemplateRem(dividend, divisor int) (int, error) {
	if divisor == 0 {
		return 0, fmt.Errorf("Division by zero: %d %% %d.", dividend, divisor)
	}

	return dividend % divisor, nil
}
