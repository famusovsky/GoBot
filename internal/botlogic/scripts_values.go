package botlogic

const (
	defaultScript = "default"
)

var (
	Scripts = NewScriptList(
		map[string]*Script{

			defaultScript: NewScript(
				[]map[string]AnswerConfig{
					0: {
						defaultScript: {
							Text:       "Заглушка",
							NextStep:   0,
							NextScript: defaultScript,
						},
					},
				},
			),
		}, // end
	)
)
