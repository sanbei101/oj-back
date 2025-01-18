package lyoj

import "oj-back/app/model"

type Meta struct {
	Alias      string   `json:"alias"`
	Difficulty int      `json:"difficulty"`
	Id         int      `json:"id"`
	Tags       []string `json:"tags"`
	Title      string   `json:"title"`
}

type Config struct {
	Datas   []data  `json:"datas"`
	Input   string  `json:"input"`
	Output  string  `json:"output"`
	Spj     spj     `json:"spj"`
	Subtask subtask `json:"subtask"`
}

type data struct {
	Input   string `json:"input"`
	Output  string `json:"output"`
	Memory  int    `json:"memory"`
	Score   int    `json:"score"`
	Subtask int    `json:"subtask"`
	Time    int    `json:"time"`
}

type spj struct {
	CompileCmd string `json:"compile_cmd"`
	ExecName   string `json:"exec_name"`
	ExecParam  string `json:"exec_param"`
	ExecPath   string `json:"exec_path"`
	Source     string `json:"source"`
	Type       int    `json:"type"`
}

type subtask struct {
	Depends []int `json:"depends"`
	Id      int   `json:"id"`
	Title   int   `json:"title"`
	Type    int   `json:"type"`
}

type LyojProblem struct {
	Meta        Meta
	Description string
	Config      Config
	TestCases   []model.Case
}
