package factory

import "fmt"

type pen interface {
	Write()
}

type pencil struct {
}

func (p *pencil) Write() {
	fmt.Println("pencil")
}

type brushPen struct {
}

func (p *brushPen) Write() {
	fmt.Println("brushPen")
}

// 简单工厂模式
type PenFactory struct {
}

func (PenFactory) ProducePencil() pen {
	return new(pencil)
}

func (PenFactory) ProduceBrushPen() pen {
	return new(brushPen)
}

func (p PenFactory) Produce(typ string) pen {
	switch typ {
	case "pencil":
		return p.ProducePencil()
	case "brush":
		return p.ProduceBrushPen()
	default:
		return nil
	}
}

// 抽象工厂模式
type AbstractFacory interface {
	Produce() pen
}

type PencilFactory struct {
}

func (PencilFactory) Produce() pen {
	return new(pencil)
}

type BrushPen struct {
}

func (BrushPen) Produce() pen {
	return new(brushPen)
}
