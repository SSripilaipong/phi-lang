package program

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/SSripilaipong/muto/core/base"
	program2 "github.com/SSripilaipong/muto/program"
)

func TestBuildFromString(t *testing.T) {
	t.Parallel()

	t.Run("should resolve to string", func(t *testing.T) {
		program := buildOrPanic(`(main) = "hello world"`)
		assert.Equal(t, base.NewString("hello world"), execute(program))
	})

	t.Run("should resolve to number", func(t *testing.T) {
		program := BuildProgramFromString(`(main) = 123.45`).Value()
		assert.Equal(t, base.NewNumberFromString("123.45"), execute(program))
	})

	t.Run("should resolve to object", func(t *testing.T) {
		program := BuildProgramFromString(`main = hello "world" 123`).Value()
		assert.Equal(t, base.NewNamedOneLayerObject("hello", base.NewString("world"), base.NewNumberFromString("123")), execute(program))
	})

	t.Run("should resolve variable", func(t *testing.T) {
		program := BuildProgramFromString(`
(hello X) = X
main = hello "world"
`).Value()
		assert.Equal(t, base.NewString("world"), execute(program))
	})

	t.Run("should match rule by string value", func(t *testing.T) {
		program := BuildProgramFromString(`
(hello "a") = 1
(hello "b") = 2
main = hello "b"
`).Value()
		assert.Equal(t, base.NewNumberFromString("2"), execute(program))
	})

	t.Run("should match rule by number value", func(t *testing.T) {
		program := BuildProgramFromString(`
(hello 1) = "hello"
(hello 2) = "world"
(hello 3) = "muto"
main = hello 2
`).Value()
		assert.Equal(t, base.NewString("world"), execute(program))
	})

	t.Run("should resolve nested object", func(t *testing.T) {
		program := BuildProgramFromString(`
(hello (f X Y)) = X
main = hello (f "abc" 123)
`).Value()
		assert.Equal(t, base.NewString("abc"), execute(program))
	})

	t.Run("should resolve variadic variable", func(t *testing.T) {
		program := BuildProgramFromString(`f X Xs... = g Xs...
main = f 1 2 3
`).Value()
		assert.Equal(t, base.NewNamedOneLayerObject("g", base.NewNumberFromString("2"), base.NewNumberFromString("3")), execute(program))
	})

	t.Run("should match nested variadic variable with size 0", func(t *testing.T) {
		program := BuildProgramFromString(`g (f Xs...) = h Xs...
main = g (f)
`).Value()
		assert.Equal(t, base.NewNamedOneLayerObject("h"), execute(program))
	})

	t.Run("should match children strictly for nested pattern", func(t *testing.T) {
		program := BuildProgramFromString(`g (f 1) = 555
main = g (f 1 2)
`).Value()
		assert.True(t, base.NewNamedOneLayerObject("g", base.NewNamedOneLayerObject("f", base.NewNumberFromString("1"), base.NewNumberFromString("2"))).
			Equals(base.UnsafeNodeToObject(execute(program))))
	})

	t.Run("should resolve to object with data head when there are children left", func(t *testing.T) {
		program := BuildProgramFromString(`f X = 999
main = f 1 2
`).Value()
		assert.Equal(t, base.NewOneLayerObject(base.NewNumberFromString("999"), base.NewNumberFromString("2")), execute(program))
	})

	t.Run("should extract nested object with variable object name pattern", func(t *testing.T) {
		program := BuildProgramFromString(`f (G X) = h (G X)
main = f (hello "world")
`).Value()
		assert.Equal(t, base.NewNamedOneLayerObject("h", base.NewNamedOneLayerObject("hello", base.NewString("world"))), execute(program))
	})

	t.Run("should build nested variable object with variadic params", func(t *testing.T) {
		program := BuildProgramFromString(`f (H X...) = g (H X...)
main = f (h "1" "2")
`).Value()
		assert.Equal(t, base.NewNamedOneLayerObject("g", base.NewNamedOneLayerObject("h", base.NewString("1"), base.NewString("2"))), execute(program))
	})

	t.Run("should not fail when variadic param part tries to match with no children", func(t *testing.T) {
		program := BuildProgramFromString(`f (G S... 0) = 0
main = f $
`).Value()
		assert.True(t, base.NewNamedOneLayerObject("f", base.NewUnlinkedRuleBasedClass("$")).
			Equals(base.UnsafeNodeToObject(execute(program))))
	})

	t.Run("should not fail when variadic right param part tries to match with no children", func(t *testing.T) {
		program := BuildProgramFromString(`f 0 S... = 0
(main) = f
`).Value()
		assert.True(t, base.NewUnlinkedRuleBasedClass("f").Equals(base.UnsafeNodeToRuleBasedClass(execute(program))))
	})

	t.Run("should apply active mutation before normal mutation", func(t *testing.T) {
		program := BuildProgramFromString(`@ (f (+ 1 X)) = X
main = f (+ 1 999)
`).Value()
		assert.Equal(t, base.NewNumberFromString("999"), execute(program))
	})

	t.Run("should be able to actively mutate children while mutating parent", func(t *testing.T) {
		program := BuildProgramFromString(`@ (f (g X)) = X
@ h = g 123
main = f (h)
`).Value()
		assert.Equal(t, base.NewNumberFromString("123"), execute(program))
	})

	t.Run("should match variable rule pattern with anonymous object", func(t *testing.T) {
		program := BuildProgramFromString(`(f (G X)) = X
main = f ((g 456) 123)
`).Value()
		assert.Equal(t, base.NewNumberFromString("123"), execute(program))
	})

	t.Run("should match variable rule pattern with anonymous object (when using active mutation)", func(t *testing.T) {
		program := BuildProgramFromString(`@ (f (G X)) = X
main = f ((g 456) 123)
`).Value()
		assert.Equal(t, base.NewNumberFromString("123"), execute(program))
	})

	t.Run("should resolve result with multiple variadic variables in param part", func(t *testing.T) {
		program := BuildProgramFromString(`f Xs... = $ Xs... Xs...
main = f 1 2 3
`).Value()
		assert.Equal(t, base.NewConventionalList(
			base.NewNumberFromString("1"), base.NewNumberFromString("2"), base.NewNumberFromString("3"),
			base.NewNumberFromString("1"), base.NewNumberFromString("2"), base.NewNumberFromString("3"),
		), execute(program))
	})

	t.Run("should match nested anonymous object", func(t *testing.T) {
		program := BuildProgramFromString(`f (G X) = h (G X)
@ (h ((g) X)) = 999
(h (g X)) = X
main = f (g 123)
`).Value()
		assert.Equal(t, base.NewNumberFromString("123"), execute(program))
	})

	t.Run("should match nested anonymous object when forcing named class to be anonymous object", func(t *testing.T) {
		program := BuildProgramFromString(`f (G X) = h ((G) X)
@ (h ((g) X)) = 999
(h (g X)) = X
main = f (g 123)
`).Value()
		assert.Equal(t, base.NewNumberFromString("999"), execute(program))
	})

	t.Run("should be able to check equality of nodes when extracting mutation in pattern", func(t *testing.T) {
		program := BuildProgramFromString(`(f X X) = 1
main = f f f
`).Value()
		assert.Equal(t, base.NewNumberFromString("1"), execute(program))
	})

	t.Run("should mutate after variable rule", func(t *testing.T) {
		program := BuildProgramFromString(`(g X) = X
f G = G 123
main = f g
`).Value()
		assert.Equal(t, base.NewNumberFromString("123"), execute(program))
	})

	t.Run("should auto bubble up when building object from variable", func(t *testing.T) {
		program := BuildProgramFromString(`f X = X
main = f (g 1)
`).Value()
		assert.True(t, base.NewNamedOneLayerObject("f", base.NewNamedOneLayerObject("g", base.NewNumberFromString("1"))).Equals(
			base.UnsafeNodeToObject(mutateOnce(program))))
	})

	t.Run("should mutate children after bubbling up to data object", func(t *testing.T) {
		program := BuildProgramFromString(`main = (f 1) (g)
(f X) = X
(g) = 2
`).Value()
		assert.Equal(t, base.NewOneLayerObject(base.NewNumberFromString("1"), base.NewNumberFromString("2")), execute(program))
	})

	t.Run("should match cases", func(t *testing.T) {
		program := BuildProgramFromString(`main = (match
 \1 [.ok]
 \X [+ 40 X]
) 2
`).Value()
		assert.Equal(t, base.NewNumberFromString("42"), execute(program))
	})

	t.Run("should extract determinant conjunction binding", func(t *testing.T) {
		// Pattern (f S)^P Y captures P as the intermediate object (f 1)
		program := BuildProgramFromString(`(f S)^P Y = $ S P Y
main = (f 1) 2
`).Value()
		expected := base.NewConventionalList(
			base.NewNumberFromString("1"),
			base.NewNamedOneLayerObject("f", base.NewNumberFromString("1")),
			base.NewNumberFromString("2"),
		)
		assert.True(t, expected.Equals(base.UnsafeNodeToObject(execute(program))))
	})

	t.Run("should extract nested determinant conjunctions at different levels", func(t *testing.T) {
		// Pattern ((f S)^P Y)^Q Z captures:
		// - P as the inner object (f 1)
		// - Q as the middle object ((f 1) 2)
		program := BuildProgramFromString(`((f S)^P Y)^Q Z = $ S P Y Q Z
main = ((f 1) 2) 3
`).Value()
		expected := base.NewConventionalList(
			base.NewNumberFromString("1"),
			base.NewNamedOneLayerObject("f", base.NewNumberFromString("1")),
			base.NewNumberFromString("2"),
			base.NewCompoundObject(
				base.NewUnlinkedRuleBasedClass("f"),
				base.NewParamChain([][]base.Node{
					{base.NewNumberFromString("1")},
					{base.NewNumberFromString("2")},
				}),
			),
			base.NewNumberFromString("3"),
		)
		assert.True(t, expected.Equals(base.UnsafeNodeToObject(execute(program))))
	})

	t.Run("should mutate constructor with no params", func(t *testing.T) {
		program := BuildProgramFromString(`main = \[hello "world"]
`).Value()
		assert.Equal(t, base.NewNamedOneLayerObject("hello", base.NewString("world")), execute(program))
	})

	t.Run("should mutate constructor and append extra params at same level", func(t *testing.T) {
		program := BuildProgramFromString(`main = (\[hello "world"] 123)
`).Value()
		expected := base.NewNamedOneLayerObject("hello", base.NewString("world"), base.NewNumberFromString("123"))
		assert.True(t, expected.Equals(base.UnsafeNodeToObject(execute(program))))
	})

	t.Run("should work with constructor in match", func(t *testing.T) {
		// Constructor produces .default and appends the matched value
		program := BuildProgramFromString(`main = (match
 \[.default]
) 123
`).Value()
		expected := base.NewOneLayerObject(base.NewTag("default"), base.NewNumberFromString("123"))
		assert.True(t, expected.Equals(base.UnsafeNodeToObject(execute(program))))
	})
}

func mutateOnce(program program2.Program) base.Node {
	return program.MutateOnce(program.InitialObject())
}

func execute(program program2.Program) base.Node {
	return program.MutateUntilTerminated(program.InitialObject())
}

func buildOrPanic(src string) program2.Program {
	program := BuildProgramFromString(src)
	if program.IsErr() {
		panic(program.Error())
	}
	return program.Value()
}
