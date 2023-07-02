package sokoban

import (
	"embed"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/zrcoder/rdor/pkg/game"
	"github.com/zrcoder/rdor/pkg/grid"
	"github.com/zrcoder/rdor/pkg/keys"
	"github.com/zrcoder/rdor/pkg/style"
	"github.com/zrcoder/rdor/pkg/style/color"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	Name             = "Sokoban"
	maxLevel         = 51
	inputPlaceholder = "1-51"

	wall      = '#'
	me        = '@'
	blank     = ' '
	slot      = 'X'
	box       = 'O'
	boxInSlot = '*'
	meInSlot  = '.'
)

//go:embed levels
var levelsFS embed.FS

func New() game.Game {
	return &sokoban{Base: game.New(Name)}
}

type sokoban struct {
	*game.Base
	blocks map[rune]string

	level    int
	upKey    key.Binding
	leftKey  key.Binding
	downKey  key.Binding
	rightKey key.Binding
	setKey   key.Binding
	input    textinput.Model

	helpGrid *grid.Grid
	grid     *grid.Grid
	myPos    grid.Position
	buf      *strings.Builder
}

func (s *sokoban) Init() tea.Cmd {
	s.ViewFunc = s.view
	s.HelpFunc = s.helpInfo
	s.KeyActionReset = s.reset
	s.blocks = map[rune]string{
		wall:      lipgloss.NewStyle().Background(color.Orange).Render(" = "),
		me:        " ⦿ ", // ♾ ⚉ ⚗︎ ⚘ ☻
		blank:     "   ",
		slot:      lipgloss.NewStyle().Background(color.Violet).Render("   "),
		box:       lipgloss.NewStyle().Background(color.Red).Render(" x "),
		boxInSlot: lipgloss.NewStyle().Background(color.Green).Render("   "),
		meInSlot:  lipgloss.NewStyle().Background(color.Violet).Render(" ⦿ "),
	}

	s.upKey = keys.Up
	s.leftKey = keys.Left
	s.downKey = keys.Down
	s.rightKey = keys.Right
	s.setKey = key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "set"),
	)
	s.Keys = []key.Binding{s.upKey, s.leftKey, s.downKey, s.rightKey, s.setKey}

	s.input = textinput.New()
	s.buf = &strings.Builder{}
	s.loadLever()
	s.input.Placeholder = inputPlaceholder
	return s.Base.Init()
}

func (s *sokoban) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	b, bcmd := s.Base.Update(msg)
	if b != s.Base {
		return b, bcmd
	}

	var icmd tea.Cmd
	s.input, icmd = s.input.Update(msg)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, s.upKey):
			s.move(grid.Up)
		case key.Matches(msg, s.leftKey):
			s.move(grid.Left)
		case key.Matches(msg, s.downKey):
			s.move(grid.Down)
		case key.Matches(msg, s.rightKey):
			s.move(grid.Right)
		case key.Matches(msg, s.setKey):
			return s, tea.Batch(bcmd, icmd, s.input.Focus())
		default:
			if msg.Type == tea.KeyEnter && s.input.Focused() {
				s.input.Blur()
				s.setted(s.input.Value())
				s.input.SetValue("")
			}
		}
	}
	return s, tea.Batch(bcmd, icmd)
}

func (s *sokoban) view() string {
	s.buf.Reset()
	s.grid.Range(func(pos grid.Position, char rune, isLineEnd bool) (end bool) {
		s.buf.WriteString(s.blocks[char])
		if isLineEnd {
			s.buf.WriteByte('\n')
		}
		return
	})
	state := style.Help.Render(fmt.Sprintf("- %d/%d - \n", s.level+1, maxLevel))
	res := lipgloss.JoinVertical(lipgloss.Center,
		strings.TrimRight(s.buf.String(), "\n"),
		state,
	)
	if s.input.Focused() {
		return lipgloss.JoinVertical(lipgloss.Left,
			res,
			"pick a level",
			s.input.View(),
		)
	}
	return res
}

func (s *sokoban) helpInfo() string {
	return "Our goal is to push all the boxes into the slots without been stuck somewhere."
}

func (s *sokoban) setted(level string) {
	n, err := strconv.Atoi(level)
	if err != nil {
		s.SetError(errors.New("invalid number"))
		return
	}
	if n < 1 || n > maxLevel+1 {
		s.SetError(errors.New("level out of range"))
		return
	}
	s.level = n - 1
	s.loadLever()
}

func (s *sokoban) loadLever() {
	data, err := levelsFS.ReadFile("levels/" + strconv.Itoa(s.level+1) + ".txt")
	if err != nil {
		panic(err)
	}
	s.grid = grid.New(string(data))
	s.helpGrid = grid.Copy(s.grid)
	s.grid.Range(func(pos grid.Position, char rune, isLineEnd bool) (end bool) {
		if char == me || char == meInSlot {
			s.myPos = pos
			return true
		}
		return
	})
}

func (s *sokoban) move(d grid.Direction) {
	pos := grid.TransForm(s.myPos, d)
	if s.grid.OutBound(pos) {
		return
	}
	switch s.grid.Get(pos) {
	case blank, slot:
		s.moveMe(pos)
	case box, boxInSlot:
		dest := grid.TransForm(pos, d)
		if s.grid.OutBound(dest) {
			return
		}
		char := s.grid.Get(dest)
		if char == blank || char == slot {
			s.moveBox(pos, dest)
			s.moveMe(pos)
		}
	}
	if s.success() {
		s.SetSuccess("")
	}
}

func (s *sokoban) moveMe(p grid.Position) {
	if s.grid.Get(p) == blank {
		s.grid.Set(p, me)
	} else {
		s.grid.Set(p, meInSlot)
	}
	if s.grid.Get(s.myPos) == me {
		s.grid.Set(s.myPos, blank)
	} else {
		s.grid.Set(s.myPos, slot)
	}
	s.myPos = p
}

func (s *sokoban) moveBox(src, dest grid.Position) {
	char := s.grid.Get(dest)
	if char == blank {
		s.grid.Set(dest, box)
	} else if char == slot {
		s.grid.Set(dest, boxInSlot)
	}
	if s.grid.Get(src) == box {
		s.grid.Set(src, blank)
	} else {
		s.grid.Set(src, slot)
	}
}

func (s *sokoban) success() bool {
	res := true
	s.grid.Range(func(pos grid.Position, char rune, isLineEnd bool) (end bool) {
		if char == box {
			res = false
			return true
		}
		return
	})
	return res
}

func (s *sokoban) reset() {
	s.grid.Copy(s.helpGrid)
	s.grid.Range(func(pos grid.Position, char rune, isLineEnd bool) (end bool) {
		if char == me || char == meInSlot {
			s.myPos = pos
			return true
		}
		return
	})
}
